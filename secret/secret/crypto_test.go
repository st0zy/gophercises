package secret_test

import (
	"testing"

	"github.com/st0zy/gophercises/secret/secret"
)

func TestEncryptDecrypt(t *testing.T) {
	stringToTest := "Testing"
	want := stringToTest

	encryptedString, err := secret.Encrypt("test_key", stringToTest)
	if err != nil {
		t.Error("Failed during encryption")
	}
	got, err := secret.Decrypt("test_key", encryptedString)

	if err != nil {
		t.Error("Failed during decryption")
	}
	if got != want {
		t.Errorf("got %s : want %s", got, want)
	}

}
