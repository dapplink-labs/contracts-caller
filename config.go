package challenger

import (
	"time"

	"github.com/urfave/cli"

	"github.com/the-web3/contracts-caller/flags"
)

type Config struct {
	ChainRpcUrl                    string
	ChainId                        uint64
	PrivateKey                     string
	Mnemonic                       string
	SequencerHDPath                string
	LoopInterval                   time.Duration
	Passphrase                     string
	TreasureManagerContractAddress string
	WithdrawManagerAddress         string
	ResubmissionTimeout            time.Duration
	NumConfirmations               uint64
	SafeAbortNonceTooLowCount      uint64

	EnableHsm  bool
	HsmAPIName string
	HsmCreden  string
	HsmAddress string
}

func NewConfig(ctx *cli.Context) (Config, error) {
	cfg := Config{
		ChainRpcUrl:                    ctx.GlobalString(flags.ChainRpcUrlFlag.Name),
		ChainId:                        ctx.GlobalUint64(flags.ChainIdFlag.Name),
		PrivateKey:                     ctx.GlobalString(flags.PrivateKeyFlag.Name),
		Mnemonic:                       ctx.GlobalString(flags.MnemonicFlag.Name),
		Passphrase:                     ctx.GlobalString(flags.PassphraseFlag.Name),
		TreasureManagerContractAddress: ctx.GlobalString(flags.TreasureManagerContractAddressFlag.Name),
		WithdrawManagerAddress:         ctx.GlobalString(flags.WithdrawManagerAddressFlag.Name),
		NumConfirmations:               ctx.GlobalUint64(flags.NumConfirmationsFlag.Name),
		SafeAbortNonceTooLowCount:      ctx.GlobalUint64(flags.SafeAbortNonceTooLowCountFlag.Name),
		LoopInterval:                   ctx.GlobalDuration(flags.LoopIntervalFlag.Name),
		EnableHsm:                      ctx.GlobalBool(flags.EnableHsmFlag.Name),
		HsmAddress:                     ctx.GlobalString(flags.HsmAddressFlag.Name),
		HsmAPIName:                     ctx.GlobalString(flags.HsmAPINameFlag.Name),
		HsmCreden:                      ctx.GlobalString(flags.HsmCredenFlag.Name),
	}
	return cfg, nil
}
