package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var (
	lock sync.Mutex
	//Marshal will marshal the object into an io reader
	//uses the json marshaller
	Marshal = func(v interface{}) (io.Reader, error) {
		b, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	}
	//UnMarshal will unmarshal r into v
	UnMarshal = func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}
)

//Save will save a representation of v to a file at path
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil
}

//Load will load data from a file at path into v
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return UnMarshal(f, v)
}
