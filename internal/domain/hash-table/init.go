package hash_table

import "sync"

type HashTable struct {
	table map[string]string
	mutex *sync.RWMutex
}

func New() *HashTable {
	return &HashTable{
		mutex: &sync.RWMutex{},
		table: make(map[string]string),
	}
}
