package tokenstore

import (
	"bufio"
	"github.com/99designs/keyring"
	"os"
)

type LocalKeyring struct {
	fileDir string
	keyPair map[string]keyring.Item
}

func (k *LocalKeyring) Get(key string) (keyring.Item, error) {
	if val, ok := k.keyPair[key]; ok {
		return val, nil
	}
	file := k.fileDir + "/" + key
	_, err := os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		return keyring.Item{}, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return keyring.Item{}, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	line, _, err := reader.ReadLine()
	if err != nil {
		return keyring.Item{}, err
	}
	item := keyring.Item{
		Key:  key,
		Data: line,
	}
	k.keyPair[key] = item
	return item, nil
}

// Returns the non-secret parts of an Item
func (k *LocalKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

// Stores an Item on the keyring
func (k *LocalKeyring) Set(item keyring.Item) error {
	file := k.fileDir + "/" + item.Key
	_, err := os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		return err

	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteAt(item.Data, 0) // Write at 0 beginning
	if err != nil {
		return err
	}
	k.keyPair[item.Key] = item
	return nil
}

// Removes the item with matching key
func (k *LocalKeyring) Remove(key string) error {
	return nil
}

// Provides a slice of all keys stored on the keyring
func (k *LocalKeyring) Keys() ([]string, error) {
	return nil, nil
}
