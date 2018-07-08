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

func TestCreateAccount1(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "microcosm")
	if err != nil {
		t.Fatal("Unable to create temporary directory")
	}

	defer os.RemoveAll(tempDir)

	password := "lol"
	address, err := CreateKey(tempDir, password)
	if err != nil {
		t.Fatal(err)
	}

	tempDirContents, err := ioutil.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Unable to list files in temporary directory: %s", tempDir)
	}

	if len(tempDirContents) != 1 {
		t.Fatalf("Incorrect number of files in temporary directory -- expected: %d, actual %d", 1, len(tempDirContents))
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

func TestCreateAccount2(t *testing.T) {
	keydir := "./if-you-create-a-subdirectory-with-this-name-you-are-a-jerk"
	password := "lol"
	_, err := CreateKey(keydir, password)
	if err == nil {
		t.Fatalf("CreateKey call with non-existent directory did not raise an error -- expected: nil, actual: %v", err)
	}
}
