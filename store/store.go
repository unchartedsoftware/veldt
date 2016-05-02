package store

// Set adds data to the store under the provided hash.
func Set(id string, hash string, data []byte) error {
	conn, err := getConnection(id)
	if err != nil {
		return err
	}
	// compress data
	data, err = compress(data[0:])
	if err != nil {
		return err
	}
	// get tile data from store
	err = conn.Set(addHash(hash), data[0:])
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// Get returns data from the store under the provided hash.
func Get(id string, hash string) ([]byte, error) {
	conn, err := getConnection(id)
	if err != nil {
		return nil, err
	}
	// get tile data from store
	data, err := conn.Get(addHash(hash))
	if err != nil {
		return nil, err
	}
	conn.Close()
	// decompress and return data
	return decompress(data[0:])
}

// Exists returns true if data exists in the store under the provided hash.
func Exists(id string, hash string) (bool, error) {
	conn, err := getConnection(id)
	if err != nil {
		return false, err
	}
	exists, err := conn.Exists(addHash(hash))
	conn.Close()
	if err != nil {
		return false, err
	}
	return exists, nil
}
