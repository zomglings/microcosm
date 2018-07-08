package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

// Tests that CreateKeys creates a JSON keyfile of the correct format when passed valid arguments
func TestCreateKeys1(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "microcosm")
	if err != nil {
		t.Fatal("Unable to create temporary directory")
	}

	defer os.RemoveAll(tempDir)

	password := "lol"
	addresses, err := CreateKeys(tempDir, password, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(addresses) != 1 {
		t.Fatalf("Incorrect number of addresses returned -- expected: %d, actual: %d", 1, len(addresses))
	}
	address := addresses[0]

	tempDirContents, err := ioutil.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Unable to list files in temporary directory: %s", tempDir)
	}

	if len(tempDirContents) != 1 {
		t.Fatalf("Incorrect number of files in temporary directory -- expected: %d, actual: %d", 1, len(tempDirContents))
	}

	keyFileInfo := tempDirContents[0]
	keyFile := path.Join(tempDir, keyFileInfo.Name())
	keyFileContents, err := ioutil.ReadFile(keyFile)
	if err != nil {
		t.Fatal(err)
	}

	var keyBody interface{}
	err = json.Unmarshal(keyFileContents, &keyBody)

	if err != nil {
		t.Fatal(err)
	}

	storedAddress := keyBody.(map[string]interface{})["address"]
	prefixedStoredAddressLower := strings.ToLower(fmt.Sprintf("0x%s", storedAddress))

	expectedStoredAddressLower := strings.ToLower(address.String())

	if expectedStoredAddressLower != prefixedStoredAddressLower {
		t.Fatalf("Incorrect address in keyfile -- expected: %s, actual: %s", expectedStoredAddressLower, prefixedStoredAddressLower)
	}
}

// Tests that CreateKeys returns an error when passed a non-existent directory as an argument
func TestCreateKeys2(t *testing.T) {
	keydir := "./if-you-create-a-subdirectory-with-this-name-you-are-a-jerk"
	password := "lol"
	_, err := CreateKeys(keydir, password, 1)
	if err == nil {
		t.Fatalf("CreateKeys call with non-existent directory did not raise an error -- expected: nil, actual: %v", err)
	}
}

// Tests that CreateKeys creates multiple keys when required to
func TestCreateKeys3(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "microcosm")
	if err != nil {
		t.Fatal("Unable to create temporary directory")
	}

	defer os.RemoveAll(tempDir)

	password := "lol"
	addresses, err := CreateKeys(tempDir, password, 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(addresses) != 2 {
		t.Fatalf("Incorrect number of addresses returned -- expected: %d, actual: %d", 2, len(addresses))
	}

	tempDirContents, err := ioutil.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Unable to list files in temporary directory: %s", tempDir)
	}

	if len(tempDirContents) != 2 {
		t.Fatalf("Incorrect number of files in temporary directory -- expected: %d, actual: %d", 1, len(tempDirContents))
	}

	keyFileAddresses := make(map[string]bool)
	for _, keyFileInfo := range tempDirContents {
		keyFile := path.Join(tempDir, keyFileInfo.Name())
		keyFileContents, err := ioutil.ReadFile(keyFile)
		if err != nil {
			t.Fatal(err)
		}
		var keyBody interface{}
		err = json.Unmarshal(keyFileContents, &keyBody)

		if err != nil {
			t.Fatal(err)
		}

		storedAddress := keyBody.(map[string]interface{})["address"]
		prefixedStoredAddressLower := strings.ToLower(fmt.Sprintf("0x%s", storedAddress))
		keyFileAddresses[prefixedStoredAddressLower] = true
	}

	for _, address := range addresses {
		expectedStoredAddressLower := strings.ToLower(address.String())
		if !keyFileAddresses[expectedStoredAddressLower] {
			t.Fatalf("Could not find address in generated keyfiles: %s", address)
		}
	}
}
