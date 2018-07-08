package accounts

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// CreateKey creates a new account in the specified key directory, unlockable with the given
// password. If the key creation is successful, the function returns the address of the new account
// with a nil error. Otherwise, it returns the 0 address with a non-nil error.
func CreateKey(keydir, password string) (common.Address, error) {
	var address common.Address

	keydirInfo, err := os.Lstat(keydir)
	if err != nil {
		return address, err
	}

	if !keydirInfo.IsDir() {
		err = fmt.Errorf("keydir: %s -- not a directory", keydir)
		return address, err
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	address, err = keystore.StoreKey(keydir, password, scryptN, scryptP)

	return address, err
}
