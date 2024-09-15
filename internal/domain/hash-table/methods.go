package hash_table

func (h *HashTable) Get(key string) (string, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	value, ok := h.table[key]

	return value, ok
}

func (h *HashTable) Set(key string, value string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.table[key] = value
}

func (h *HashTable) Delete(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.table, key)
}
