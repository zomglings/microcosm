package genesis

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Internal structure for testing different invocations of CreateGenesisBlock
// All members except resultsInError are part of the CreateGenesisBlock signature and explained in
// more detail in its doc string.
// resultsInError specifies whether or not the given invocation should return a non-nil error -- if
// it is true, the expectation is that the returned error object will be non-nil and if false, the
// error object should be nil.
type testCaseCreateGenesisBlock struct {
	chainID        int64
	difficulty     int64
	defaultBalance int64
	gasLimit       uint64
	accounts       []common.Address
	resultsInError bool
}

func TestCreateGenesisBlock1(t *testing.T) {
	var chainID, difficulty, defaultBalance int64
	var gasLimit uint64
	var accounts = make([]common.Address, 1)

	// Make a random address to allocate ether to
	addressBytes := make([]byte, common.AddressLength)
	_, err := rand.Read(addressBytes)
	if err != nil {
		t.Fatal(err)
	}
	address := common.BytesToAddress(addressBytes)
	accounts[0] = address

	chainID = 42
	difficulty = 999999999999
	defaultBalance = 10000000000
	gasLimit = 1000000

	genesisBlock, err := CreateGenesisBlock(chainID, difficulty, defaultBalance, gasLimit, accounts)
	if err != nil {
		t.Fatal(err)
	}

	if blockChainID := genesisBlock.Config.ChainID.Int64(); blockChainID != chainID {
		t.Fatalf("Incorrect chain ID assigned to genesis block -- expected: %d, actual: %d", chainID, blockChainID)
	}

	if blockDifficulty := genesisBlock.Difficulty.Int64(); blockDifficulty != difficulty {
		t.Fatalf("Incorrect difficulty assigned to genesis block -- expected: %d, actual: %d", difficulty, blockDifficulty)
	}

	if genesisBlock.GasLimit != gasLimit {
		t.Fatalf("Incorrect gas limit assigned to genesis block -- expected: %d, actual: %d", gasLimit, genesisBlock.GasLimit)
	}

	allocatedAccount := genesisBlock.Alloc[address]
	if accountBalance := allocatedAccount.Balance.Int64(); accountBalance != defaultBalance {
		t.Fatalf("Incorrect balance allocated to address %s -- expected %d, actual %d", address.String(), defaultBalance, accountBalance)
	}
}

func TestCreateGenesisBlock2(t *testing.T) {
	// Make random addresses to allocate ether to
	numAccounts := 5
	var accounts = make([]common.Address, numAccounts)

	for i := 0; i < numAccounts; i++ {
		addressBytes := make([]byte, common.AddressLength)
		_, err := rand.Read(addressBytes)
		if err != nil {
			t.Fatal(err)
		}
		address := common.BytesToAddress(addressBytes)
		accounts[i] = address
	}

	var testCases = []testCaseCreateGenesisBlock{
		{0, 999999999, 10000000000, 1000000, accounts, false},
		{123989813, 999999999, 10000000000, 1000000, accounts, false},
		{-1, 999999999, 10000000000, 1000000, accounts, true},
		{0, 0, 10000000000, 1000000, accounts, true},
		{12938921832, 0, 10000000000, 1000000, accounts, true},
		{12938921832, 10, -5, 1000000, accounts, true},
		{12938921832, 10, 10000000000, 0, accounts, true},
	}

	for _, tt := range testCases {
		_, err := CreateGenesisBlock(tt.chainID, tt.difficulty, tt.defaultBalance, tt.gasLimit, tt.accounts)
		if tt.resultsInError && err == nil {
			t.Errorf("CreateGenesisBlock call did not result in error -- %v\n", tt)
		}
		if !tt.resultsInError && err != nil {
			t.Error(err)
		}
	}
}
