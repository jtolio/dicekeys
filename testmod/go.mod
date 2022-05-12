module github.com/jtolio/dicekeys/testmod

go 1.18

require (
	github.com/dsnet/try v0.0.3
	github.com/ethereum/go-ethereum v1.10.17
	github.com/jtolio/dicekeys v0.0.0
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/jtolio/dicekeys => ./..
