package challenger

import (
	"context"

	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/contracts-caller/caller"
	common2 "github.com/the-web3/contracts-caller/common"
	"github.com/the-web3/contracts-caller/ethereumcli"
)

func Main(gitVersion string) func(ctx *cli.Context) error {
	return func(cliCtx *cli.Context) error {
		cfg, err := NewConfig(cliCtx)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		callerPrivateKey, _, err := common2.ParseWalletPrivKeyAndContractAddr(
			"ContractCaller", cfg.Mnemonic, cfg.SequencerHDPath,
			cfg.PrivateKey, cfg.TreasureManagerContractAddress, cfg.Passphrase,
		)
		if err != nil {
			return err
		}
		chainClient, err := ethereumcli.EthClientWithTimeout(ctx, cfg.ChainRpcUrl)
		if err != nil {
			return err
		}
		log.Info("Contract Caller Client init success")

		chainID, err := chainClient.ChainID(ctx)
		if err != nil {
			return err
		}
		callerConfig := &caller.ContractCallerConfig{
			ChainClient:               chainClient,
			ChainID:                   chainID,
			TreasureManagerAddr:       common.Address(common.HexToAddress(cfg.TreasureManagerContractAddress)),
			WithdrawManageAddr:        cfg.WithdrawManagerAddress,
			PrivateKey:                callerPrivateKey,
			LoopInterval:              cfg.LoopInterval,
			NumConfirmations:          cfg.NumConfirmations,
			SafeAbortNonceTooLowCount: cfg.SafeAbortNonceTooLowCount,
			EnableHsm:                 cfg.EnableHsm,
			HsmCreden:                 cfg.HsmCreden,
			HsmAPIName:                cfg.HsmAPIName,
			HsmAddress:                cfg.HsmAddress,
		}
		log.Info("Contract caller hsm", "EnableHsm", cfg.EnableHsm, "HsmAPIName", cfg.HsmAPIName, "HsmAddress", cfg.HsmAddress)
		cCaller, err := caller.NewContractCaller(ctx, callerConfig)
		if err != nil {
			return err
		}
		if err := cCaller.Start(); err != nil {
			return err
		}
		log.Info("Contract caller service start")
		defer cCaller.Stop()
		return nil
	}
}
