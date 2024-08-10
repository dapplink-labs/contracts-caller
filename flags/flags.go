package flags

import (
	"github.com/urfave/cli"
	"time"
)

const envVarPrefix = "CONTRACTS_CALLER"

func prefixEnvVar(suffix string) string {
	return envVarPrefix + "_" + suffix
}

var (
	ChainRpcUrlFlag = cli.StringFlag{
		Name:   "chain-rpc-url",
		Usage:  "HTTP provider URL for chain",
		EnvVar: prefixEnvVar("CHAIN_RPC_URL"),
		Value:  "http://127.0.0.1:8545",
	}
	ChainIdFlag = cli.Uint64Flag{
		Name:   "chain-id",
		Usage:  "Chain id for evm chain",
		EnvVar: prefixEnvVar("CHAIN_ID"),
		Value:  31337,
	}
	PrivateKeyFlag = cli.StringFlag{
		Name:   "private-key",
		Usage:  "Ethereum private key for node operator",
		EnvVar: prefixEnvVar("PRIVATE_KEY"),
		Value:  "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	}
	TreasureManagerContractAddressFlag = cli.StringFlag{
		Name:   "treasure-manage-address",
		Usage:  "Address of the treasure manager contract",
		EnvVar: prefixEnvVar("TREASURE_MANAGER_ADDRESS"),
		Value:  "0x0B306BF915C4d645ff596e518fAf3F9669b97016",
	}
	WithdrawManagerAddressFlag = cli.StringFlag{
		Name:   "treasure-manage-address",
		Usage:  "Address of the treasure manager contract",
		EnvVar: prefixEnvVar("WITHDRAW_MANAGER_ADDRESS"),
		Value:  "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720",
	}
	LoopIntervalFlag = cli.DurationFlag{
		Name:   "loop-interval",
		Usage:  "main worker lopp interval",
		EnvVar: prefixEnvVar("LOOP_INTERVAL"),
		Value:  time.Second * 5,
	}
	NumConfirmationsFlag = cli.Uint64Flag{
		Name: "num-confirmations",
		Usage: "Number of confirmations which we will wait after " +
			"appending a new batch",
		EnvVar: prefixEnvVar("NUM_CONFIRMATIONS"),
		Value:  1,
	}
	SafeAbortNonceTooLowCountFlag = cli.Uint64Flag{
		Name: "safe-abort-nonce-too-low-count",
		Usage: "Number of ErrNonceTooLow observations required to " +
			"give up on a tx at a particular nonce without receiving " +
			"confirmation",
		EnvVar: prefixEnvVar("SAFE_ABORT_NONCE_TOO_LOW_COUNT"),
		Value:  3,
	}
	MnemonicFlag = cli.StringFlag{
		Name: "mnemonic",
		Usage: "The mnemonic used to derive the wallets for either the " +
			"sequencer or the proposer",
		EnvVar: prefixEnvVar("MNEMONIC"),
	}
	CallerHDPathFlag = cli.StringFlag{
		Name: "sequencer-hd-path",
		Usage: "The HD path used to derive the sequencer wallet from the " +
			"mnemonic. The mnemonic flag must also be set.",
		EnvVar: prefixEnvVar("CALLER_HD_PATH"),
	}
	PassphraseFlag = cli.StringFlag{
		Name:   "passphrase",
		Usage:  "passphrase for the seed generation process to increase the seed's security",
		EnvVar: prefixEnvVar("PASSPHRASE"),
	}
	EnableHsmFlag = cli.BoolFlag{
		Name:   "enable-hsm",
		Usage:  "Enalbe the hsm",
		EnvVar: prefixEnvVar("ENABLE_HSM"),
	}
	HsmAPINameFlag = cli.StringFlag{
		Name:   "hsm-api-name",
		Usage:  "the api name of hsm",
		EnvVar: prefixEnvVar("HSM_API_NAME"),
	}
	HsmAddressFlag = cli.StringFlag{
		Name:   "hsm-address",
		Usage:  "the address of hsm key",
		EnvVar: prefixEnvVar("HSM_ADDRESS"),
	}
	HsmCredenFlag = cli.StringFlag{
		Name:   "hsm-creden",
		Usage:  "the creden of hsm key",
		EnvVar: prefixEnvVar("HSM_CREDEN"),
	}
)

var requiredFlags = []cli.Flag{
	ChainRpcUrlFlag,
	ChainIdFlag,
	PrivateKeyFlag,
	NumConfirmationsFlag,
	SafeAbortNonceTooLowCountFlag,
	TreasureManagerContractAddressFlag,
	LoopIntervalFlag,
}

var optionalFlags = []cli.Flag{
	MnemonicFlag,
	CallerHDPathFlag,
	PassphraseFlag,
	EnableHsmFlag,
	HsmAddressFlag,
	HsmAPINameFlag,
	HsmCredenFlag,
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

var Flags []cli.Flag
