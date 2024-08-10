package caller

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethc "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/the-web3/contracts-caller/bindings"
	common2 "github.com/the-web3/contracts-caller/common"
	"github.com/the-web3/contracts-caller/txmgr"
)

type SignerFn func(context.Context, ethc.Address, *types.Transaction) (*types.Transaction, error)

var (
	errMaxPriorityFeePerGasNotFound = errors.New(
		"Method eth_maxPriorityFeePerGas not found",
	)
	FallbackGasTipCap = big.NewInt(1500000000)
)

type ContractCallerConfig struct {
	ChainClient               *ethclient.Client
	ChainID                   *big.Int
	TreasureManagerAddr       ethc.Address
	WithdrawManageAddr        string
	PrivateKey                *ecdsa.PrivateKey
	LoopInterval              time.Duration
	SignerFn                  SignerFn
	NumConfirmations          uint64
	SafeAbortNonceTooLowCount uint64
	EnableHsm                 bool
	HsmAPIName                string
	HsmCreden                 string
	HsmAddress                string
}

type ContractCaller struct {
	Ctx                        context.Context
	Cfg                        *ContractCallerConfig
	TreasureManagerContract    *bindings.TreasureManager
	RawTreasureManagerContract *bind.BoundContract
	WalletAddr                 ethc.Address
	TreasureManagerABI         *abi.ABI
	txMgr                      txmgr.TxManager
	cancel                     func()
	wg                         sync.WaitGroup
	once                       sync.Once
}

func NewContractCaller(ctx context.Context, cfg *ContractCallerConfig) (*ContractCaller, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	treasureManagerContract, err := bindings.NewTreasureManager(
		ethc.Address(cfg.TreasureManagerAddr), cfg.ChainClient,
	)
	if err != nil {
		return nil, err
	}
	parsed, err := abi.JSON(strings.NewReader(
		bindings.TreasureManagerMetaData.ABI,
	))
	if err != nil {
		log.Error("MtChallenger parse eigen layer contract abi fail", "err", err)
		return nil, err
	}
	treasureManagerABI, err := bindings.TreasureManagerMetaData.GetAbi()
	if err != nil {
		log.Error("MtChallenger get eigen layer contract abi fail", "err", err)
		return nil, err
	}
	rawTreasureManagerContract := bind.NewBoundContract(
		cfg.TreasureManagerAddr, parsed, cfg.ChainClient, cfg.ChainClient,
		cfg.ChainClient,
	)

	txManagerConfig := txmgr.Config{
		ResubmissionTimeout:       time.Second * 5,
		ReceiptQueryInterval:      time.Second,
		NumConfirmations:          cfg.NumConfirmations,
		SafeAbortNonceTooLowCount: cfg.SafeAbortNonceTooLowCount,
	}
	txMgr := txmgr.NewSimpleTxManager(txManagerConfig, cfg.ChainClient)
	var walletAddr ethc.Address
	if cfg.EnableHsm {
		walletAddr = ethc.HexToAddress(cfg.HsmAddress)
	} else {
		walletAddr = crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)
	}
	return &ContractCaller{
		Cfg:                        cfg,
		Ctx:                        ctx,
		TreasureManagerContract:    treasureManagerContract,
		RawTreasureManagerContract: rawTreasureManagerContract,
		WalletAddr:                 walletAddr,
		TreasureManagerABI:         treasureManagerABI,
		txMgr:                      txMgr,
		cancel:                     cancel,
	}, nil
}

func (c *ContractCaller) UpdateGasPrice(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	var opts *bind.TransactOpts
	var err error
	if !c.Cfg.EnableHsm {
		opts, err = bind.NewKeyedTransactorWithChainID(
			c.Cfg.PrivateKey, c.Cfg.ChainID,
		)
	} else {
		opts, err = common2.NewHSMTransactOpts(ctx, c.Cfg.HsmAPIName,
			c.Cfg.HsmAddress, c.Cfg.ChainID, c.Cfg.HsmCreden)
	}
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.NoSend = true

	finalTx, err := c.RawTreasureManagerContract.RawTransact(opts, tx.Data())
	switch {
	case err == nil:
		return finalTx, nil

	case c.IsMaxPriorityFeePerGasNotFoundError(err):
		log.Info("MtChallenger eth_maxPriorityFeePerGas is unsupported by current backend, using fallback gasTipCap", "txData", tx.Data())
		opts.GasTipCap = FallbackGasTipCap
		return c.RawTreasureManagerContract.RawTransact(opts, tx.Data())

	default:
		return nil, err
	}
}

func (c *ContractCaller) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.Cfg.ChainClient.SendTransaction(ctx, tx)
}

func (c *ContractCaller) IsMaxPriorityFeePerGasNotFoundError(err error) bool {
	return strings.Contains(
		err.Error(), errMaxPriorityFeePerGasNotFound.Error(),
	)
}

func (c *ContractCaller) setWithdrawManagerTx(ctx context.Context, address ethc.Address) (*types.Transaction, error) {
	balance, err := c.Cfg.ChainClient.BalanceAt(
		c.Ctx, ethc.Address(c.WalletAddr), nil,
	)
	if err != nil {
		log.Error("Contract caller unable to get current balance", "err", err)
		return nil, err
	}
	log.Info("Contract wallet address balance", "balance", balance)

	nonce64, err := c.Cfg.ChainClient.NonceAt(
		c.Ctx, ethc.Address(c.WalletAddr), nil,
	)
	if err != nil {
		log.Error("Contract wallet unable to get current nonce", "err", err)
		return nil, err
	}
	nonce := new(big.Int).SetUint64(nonce64)
	var opts *bind.TransactOpts
	if !c.Cfg.EnableHsm {
		opts, err = bind.NewKeyedTransactorWithChainID(
			c.Cfg.PrivateKey, c.Cfg.ChainID,
		)
	} else {
		opts, err = common2.NewHSMTransactOpts(ctx, c.Cfg.HsmAPIName,
			c.Cfg.HsmAddress, c.Cfg.ChainID, c.Cfg.HsmCreden)
	}
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = nonce
	opts.NoSend = true

	tx, err := c.TreasureManagerContract.SetWithdrawManager(opts, address)
	switch {
	case err == nil:
		return tx, nil

	case c.IsMaxPriorityFeePerGasNotFoundError(err):
		log.Warn("contract callet eth_maxPriorityFeePerGas is unsupported by current backend, using fallback gasTipCap")
		opts.GasTipCap = FallbackGasTipCap
		return c.TreasureManagerContract.SetWithdrawManager(opts, address)
	default:
		return nil, err
	}
}

func (c *ContractCaller) setWithdrawManager(address string) (*types.Transaction, error) {
	tx, err := c.setWithdrawManagerTx(c.Ctx, ethc.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	updateGasPrice := func(ctx context.Context) (*types.Transaction, error) {
		log.Info("Contract caller setWithdrawManager update gas price")
		return c.UpdateGasPrice(ctx, tx)
	}
	receipt, err := c.txMgr.Send(
		c.Ctx, updateGasPrice, c.SendTransaction,
	)
	if err != nil {
		return nil, err
	}
	log.Info("Contract caller set withdraw manager success", "TxHash", receipt.TxHash)
	return tx, nil
}

func (c *ContractCaller) Start() error {
	c.wg.Add(1)
	go c.eventLoop()
	c.once.Do(func() {
		log.Info("Contract caller start exec set withdraw manager")
		tx, err := c.setWithdrawManager(c.Cfg.WithdrawManageAddr)
		if err != nil {
			log.Error("Contract caller set withdraw manager fail", "WithdrawManageAddr", c.Cfg.WithdrawManageAddr, "err", err)
		}
		log.Info("Contract caller set withdraw manager success", "WithdrawManageAddr", c.Cfg.WithdrawManageAddr, "txHash", tx.Hash().String())
	})
	return nil
}

func (c *ContractCaller) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *ContractCaller) eventLoop() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.Cfg.LoopInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info("Contract caller get loop")
			addressList, err := c.TreasureManagerContract.GetTokenWhiteList(&bind.CallOpts{})
			if err != nil {
				log.Error("get token white list fail", "err", err)
				continue
			}
			for _, address := range addressList {
				log.Info("token white list address", "address", address.String())
			}

			withdrawManagerAddr, _ := c.TreasureManagerContract.WithdrawManager(&bind.CallOpts{})
			log.Info("withdraw manager address", "withdrawManagerAddr", withdrawManagerAddr.String())

			treasureManageAddress, _ := c.TreasureManagerContract.TreasureManager(&bind.CallOpts{})
			log.Info("treasure manage address", "treasureManageAddress", treasureManageAddress.String())

		case err := <-c.Ctx.Done():
			log.Error("Contract caller service shutting down", "err", err)
			return
		}
	}
}
