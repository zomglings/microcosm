package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// CreateKeys creates a new account in the specified key directory, unlockable with the given
// password. If the key creation is successful, the function returns the address of the new account
// with a nil error. Otherwise, it returns the 0 address with a non-nil error.
func CreateKeys(keydir, password string, numKeys uint) ([]common.Address, error) {
	var address common.Address
	addresses := make([]common.Address, numKeys)

	keydirInfo, err := os.Lstat(keydir)
	if err != nil {
		return addresses, err
	}

	if !keydirInfo.IsDir() {
		err = fmt.Errorf("keydir: %s -- not a directory", keydir)
		return addresses, err
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	var i uint
	for i = 0; i < numKeys; i++ {
		address, err = keystore.StoreKey(keydir, password, scryptN, scryptP)
		if err != nil {
			return addresses, err
		}
		addresses[i] = address
	}

	return addresses, err
}

// GetAddress examines an Ethereum keyfile, and returns the address it represents. It returns an
// error if it failed to read the file or to identify an address in it.
func GetAddress(keyfile string) (common.Address, error) {
	var address common.Address
	var decodedContent interface{}

	content, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return address, err
	}

	err = json.Unmarshal(content, &decodedContent)
	if err != nil {
		return address, err
	}

	decodedAddress := decodedContent.(map[string]interface{})["address"].(string)
	address = common.BytesToAddress([]byte(common.FromHex(decodedAddress)))
	return address, nil
}
