package caches

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type dump struct {
	Data map[string]*value
	Options Options
	Status *Status
}

func newEmptyDump() *dump {
	return &dump{}
}

func newDump(c *Cache) *dump {
	return &dump{
		Data: c.data,
		Options: c.options,
		Status: c.status,
	}
}

func nowSuffix() string {
	return "."+time.Now().Format("20060102150405")
}

func (d *dump) to(dumpFile string) error {
	newDumpFile := dumpFile + nowSuffix()
	file, err := os.OpenFile(newDumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(d)
	if err != nil {
		file.Close()
		os.Remove(newDumpFile)
		return err
	}

	os.Remove(dumpFile)
	file.Close()
	return os.Rename(newDumpFile, dumpFile)
}

func (d *dump) from(dumpFile string) (*Cache, error) {
	file, err := os.Open(dumpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(d); err != nil {
		return nil, err
	}

	return &Cache{
		data: d.Data,
		options: d.Options,
		status: d.Status,
		lock: &sync.RWMutex{},
	}, nil
}
