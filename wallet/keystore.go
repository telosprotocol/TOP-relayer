package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wonderivan/logger"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

func createAccount(store *keystore.KeyStore, pass string) (accounts.Account, error) {
	acc, err := store.NewAccount(pass)
	if err != nil {
		return accounts.Account{}, err
	}
	logger.Info("Created new account:", acc.Address.Hex())
	return acc, nil
}

func loadAccount(store *keystore.KeyStore, path, pass string) (accounts.Account, error) {
	var keyfiles []string
	{
		_, err := os.Stat(path)
		if err != nil {
			_, err := createAccount(store, pass)
			if err != nil {
				return accounts.Account{}, err
			}
		}
	}
	kfs, err := getKeyfiles(path, keyfiles)
	if err != nil {
		logger.Error("getKeyfiles error:", err, "path:", path)
		return accounts.Account{}, err
	}
	keyfiles = kfs
	if len(keyfiles) == 0 {
		acc, err := createAccount(store, pass)
		if err != nil {
			return accounts.Account{}, err
		}
		return acc, nil
	}

	file, err := os.OpenFile(keyfiles[0], os.O_RDONLY, 0600)
	if err != nil {
		return accounts.Account{}, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return accounts.Account{}, err
	}

	filesize := fi.Size()
	data := make([]byte, filesize)
	_, err = file.Read(data)
	if err != nil {
		return accounts.Account{}, err
	}
	if len(data) == 0 {
		return accounts.Account{}, fmt.Errorf("empty or wrong keystore file")
	}

	var jsonkey encryptedKeyJSONV3
	err = json.Unmarshal(data, &jsonkey)
	if err != nil {
		logger.Error("json.Unmarshal keystore error:", err)
		return accounts.Account{}, err
	}
	var acc accounts.Account
	acc.Address = common.HexToAddress("0x" + jsonkey.Address)
	acc.URL.Scheme = keystore.KeyStoreScheme
	acc.URL.Path = keyfiles[0]

	return acc, nil
}

func getKeyfiles(keyfilepath string, files []string) ([]string, error) {
	dir, err := ioutil.ReadDir(keyfilepath)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if !fi.IsDir() {
			dir, err := filepath.Abs(keyfilepath)
			if err != nil {
				return nil, err
			}
			fullName := dir + "/" + fi.Name()
			files = append(files, fullName)
		}
	}
	return files, nil

}
