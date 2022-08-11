package wallet

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func TestCreateAccount(t *testing.T) {
	store := keystore.NewKeyStore(defaultPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := createAccount(store, defaultPass)
	if err != nil {
		t.Fatal("CreateAccount error:", err)
	}
	t.Log("create account ok!", acc.Address)
}

func TestLoadAccount(t *testing.T) {
	store := keystore.NewKeyStore(defaultPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := loadAccount(store, defaultPath, defaultPass)
	if err != nil {
		t.Fatal("load accounts failed,error:", err)
	}

	t.Logf("account:%v", acc.Address)
}

func TestUnlockAccount(t *testing.T) {
	store := keystore.NewKeyStore(defaultPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := loadAccount(store, defaultPath, defaultPass)
	if err != nil {
		t.Fatal("load accounts failed,error:", err)
	}

	err = store.Unlock(acc, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("unlock acc ok! ", acc.Address)
}
