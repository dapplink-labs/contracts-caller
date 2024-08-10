GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

TM_ABI_ARTIFACT := ./abis/TreasureManager.sol/TreasureManager.json


contracts-caller:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/contracts-caller

clean:
	rm contracts-caller

test:
	go test -v ./...

lint:
	golangci-lint run ./...

bindings:
	$(eval temp := $(shell mktemp))

	cat $(TM_ABI_ARTIFACT) \
		| jq -r .bytecode > $(temp)

	cat $(TM_ABI_ARTIFACT) \
		| jq .abi \
		| abigen --pkg bindings \
		--abi - \
		--out bindings/treasure_manager.go \
		--type TreasureManager \
		--bin $(temp)

		rm $(temp)

.PHONY: \
	contracts-caller \
	bindings \
	clean \
	test \
	lint