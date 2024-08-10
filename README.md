<!--
parent:
  order: false
-->

<div align="center">
  <h1> Contract Caller Repo </h1>
</div>

<div align="center">
  <a href="https://github.com/the-web3/contracts-caller/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/the-web3/contracts-caller.svg" />
  </a>
  <a href="https://github.com/the-web3/contracts-caller/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/the-web3/contracts-caller.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/the-web3/contracts-caller">
    <img alt="GoDoc" src="https://godoc.org/github.com/the-web3/contracts-caller?status.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/the-web3/contracts-caller">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/the-web3/contracts-caller"/>
  </a>
</div>

Contract Caller is contracts call template project, everyone can develop theirs business requirement base this project

**Note**: Requires [Go 1.18+](https://golang.org/dl/)

## Installation

For prerequisites and detailed build instructions please read the [Installation](https://github.com/the-web3/contracts-caller/) instructions. Once the dependencies are installed, run:

```bash
go mod tidy
```

Or check out the latest [release](https://github.com/the-web3/contracts-caller).

## Quick Start

Generate bindings 
```
make bindings
```
Build project
```
make 
make contracts-caller
```

Setup and run project

- change env config according to your requirement, please env example [.evn](https://github.com/the-web3/contracts-caller/.env)
```
source .env
```

Run
```
./contracts-caller
```

If you run succcess, you can see following logs
```
INFO [08-10|20:51:03.084] ContractCaller wallet params parsed successfully wallet_address=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 contract_address=0x0B306BF915C4d645ff596e518fAf3F9669b97016
INFO [08-10|20:51:03.084] Contract Caller Client init success
INFO [08-10|20:51:03.085] Contract caller hsm                      EnableHsm=false HsmAPIName= HsmAddress=
INFO [08-10|20:51:03.086] Contract caller start exec set withdraw manager
INFO [08-10|20:51:03.087] Contract wallet address balance          balance=9,999,993,061,903,973,142,891
INFO [08-10|20:51:03.088] Contract caller setWithdrawManager update gas price
INFO [08-10|20:51:04.092] Contract caller set withdraw manager success TxHash=e7d3d3..1f22cb
INFO [08-10|20:51:04.093] Contract caller set withdraw manager success WithdrawManageAddr=0x0B306BF915C4d645ff596e518fAf3F9669b97016 txHash=0xe7d3d3010e6d358df4f1b682688067c06e72bd342aa45b0921a08d76e31f22cb
INFO [08-10|20:51:04.093] Contract caller service start
INFO [08-10|20:51:08.087] Contract caller get loop
INFO [08-10|20:51:08.089] token white list address                 address=0xdAC17F958D2ee523a2206206994597C13D831ec7
INFO [08-10|20:51:08.089] token white list address                 address=0x8D983cb9388EaC77af0474fA441C4815500Cb7BB
INFO [08-10|20:51:08.089] token white list address                 address=0x3c3a81e81dc49A522A592e7622A7E711c06bf354
INFO [08-10|20:51:08.090] withdraw manager address                 withdrawManagerAddr=0x0B306BF915C4d645ff596e518fAf3F9669b97016
INFO [08-10|20:51:08.091] treasure manage address                  treasureManageAddress=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

## Contributing

Looking for a good place to start contributing? Check out some [`good first issues`](https://github.com/the-web3/contracts-caller/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).

For additional instructions, standards and style guides, please refer to the [Contributing](./CONTRIBUTING.md) document.
