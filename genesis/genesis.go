package genesis

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

// CreateGenesisBlock generates a genesis block with the given configuration
// For more information, look at the Genesis struct: https://github.com/ethereum/go-ethereum/blob/v1.8.12/core/genesis.go
func CreateGenesisBlock(chainID, difficulty, defaultBalance int64, gasLimit uint64, accounts []common.Address) (core.Genesis, error) {
	var genesisBlock core.Genesis
	var err error

	if chainID < 0 {
		err = fmt.Errorf("Received negative chainID: %d", chainID)
		return genesisBlock, err
	}

	if difficulty <= 0 {
		err = fmt.Errorf("Received non-positive difficulty: %d", difficulty)
		return genesisBlock, err
	}

	if defaultBalance < 0 {
		err = fmt.Errorf("Received negative default balance: %d", defaultBalance)
		return genesisBlock, err
	}

	if gasLimit <= 0 {
		err = fmt.Errorf("Received non-positive gas limit: %d", gasLimit)
		return genesisBlock, err
	}

	genesisBlock.Difficulty = big.NewInt(difficulty)
	genesisBlock.GasLimit = gasLimit

	// Generate a random nonce for the genesis block
	nonceBuffer := make([]byte, 8)
	_, err = rand.Read(nonceBuffer)
	if err != nil {
		return genesisBlock, err
	}
	genesisBlock.Nonce = binary.BigEndian.Uint64(nonceBuffer)

	bigChainID := big.NewInt(chainID)
	chainConfig := params.ChainConfig{ChainID: bigChainID}
	genesisBlock.Config = &chainConfig

	bigDefaultBalance := big.NewInt(defaultBalance)
	genesisAllocation := make(core.GenesisAlloc)
	for _, address := range accounts {
		genesisAllocation[address] = core.GenesisAccount{Balance: bigDefaultBalance}
	}
	genesisBlock.Alloc = genesisAllocation

	return genesisBlock, nil
}
