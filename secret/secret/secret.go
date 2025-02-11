package secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

type InMemoryVault struct {
	encryptionKey string
	secretMap     map[string]string
}

type FileVault struct {
	encryptionKey string
	filePath      string
	mutex         sync.Mutex
	secretMap     map[string]string
}

func NewInMemoryVault() *InMemoryVault {
	return &InMemoryVault{
		encryptionKey: "test123",
		secretMap:     map[string]string{},
	}
}

func NewFileVault(opts ...Option) *FileVault {
	f := &FileVault{
		encryptionKey: "test123",
		filePath:      "secrets.db",
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

type Option func(f *FileVault)

func WithVaultPath(path string) Option {
	return func(f *FileVault) {
		f.filePath = path
	}
}

func WithEncryptionKey(encryptionKey string) Option {
	return func(f *FileVault) {
		f.encryptionKey = encryptionKey
	}
}

func (f *InMemoryVault) Set(key, value string) error {
	encryptedVal, err := Encrypt(f.encryptionKey, value)
	if err != nil {
		return err
	}
	err = f.save(key, encryptedVal)
	return err
}

func (f *InMemoryVault) Get(key string) (string, error) {
	var encryptedVal string
	var ok bool
	if encryptedVal, ok = f.secretMap[key]; !ok {
		return "", errors.New("key does not exist")
	}
	ret, err := Decrypt(f.encryptionKey, encryptedVal)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (f *InMemoryVault) save(key, encryptedValue string) error {
	f.secretMap[key] = encryptedValue
	return nil
}

func (f *FileVault) Set(key, value string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	err := f.loadKeyValues()
	if err != nil {
		return err
	}
	f.secretMap[key] = value
	err = f.saveKeyValues()
	return err
}

func (f *FileVault) Get(key string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	err := f.loadKeyValues()
	if err != nil {
		return "", err
	}
	var ret string
	var ok bool
	if ret, ok = f.secretMap[key]; !ok {
		return "", errors.New("key does not exist")
	}
	return ret, nil
}

func (f *FileVault) loadKeyValues() error {
	file, err := os.Open(f.filePath)
	if err != nil {
		f.secretMap = map[string]string{}
		return nil
	}
	defer file.Close()

	decryptedReader, err := DecryptReader(f.encryptionKey, file)
	if err != nil {
		return err
	}
	return f.readKeyValues(decryptedReader)
}

func (f *FileVault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&f.secretMap)
}

func (f *FileVault) saveKeyValues() error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := EncryptWriter(f.encryptionKey, file)
	if err != nil {
		return err
	}
	return f.writeKeyValues(writer)
}

func (f *FileVault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(f.secretMap)
}
