package secret_test

import (
	"testing"

	"github.com/st0zy/gophercises/secret/secret"
)

func TestInMemoryVault(t *testing.T) {
	key := "test-key"
	value := "test-value"
	want := value

	vault := secret.NewInMemoryVault()

	err := vault.Set(key, value)
	if err != nil {
		t.Error("Failed during encryption")
	}
	got, err := vault.Get(key)
	if err != nil {
		t.Error("Failed during decryption")
	}
	if got != want {
		t.Errorf("got %s : want %s", got, want)
	}

}

func TestFileVault(t *testing.T) {
	key := "test-key"
	value := "test-value"
	want := value

	vault := secret.NewFileVault(secret.WithVaultPath("test.db"))

	err := vault.Set(key, value)
	if err != nil {
		t.Error(err.Error())
	}

	err = vault.Set("somethingelse", "somethingelse")
	if err != nil {
		t.Error(err.Error())
	}
	got, err := vault.Get(key)
	if err != nil {
		t.Error(err.Error())
	}
	if got != want {
		t.Errorf("got %s : want %s", got, want)
	}

}
