package hello

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type Messages struct {
	sync.RWMutex

	messages map[string]string
}

func (m *Messages) Load(path string) (err error) {
	m.Lock()
	defer m.Unlock()

	// Load the messages from the JSON file
	m.messages = make(map[string]string)

	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()

	var data []byte
	if data, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	if err = json.Unmarshal(data, &m.messages); err != nil {
		return err
	}

	return nil
}

func (m *Messages) Get(key string) (value string, err error) {
	m.RLock()
	defer m.RUnlock()

	if value, ok := m.messages[key]; ok {
		return value, nil
	}
	return "", errors.New("language not found")
}
