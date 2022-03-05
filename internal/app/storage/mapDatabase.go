package app

import (
	"hash/crc32"
	"strconv"
)

type MapDatabase struct {
	urls map[string]string
}

func NewMapStorage() *MapDatabase {
	var md MapDatabase
	md.urls = make(map[string]string)
	return &md
}

func (m *MapDatabase) Find(id string) (string, error) {
	return m.urls[id], nil
}

func (m *MapDatabase) Save(url string) (string, error) {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))
	m.urls[checksum] = url

	return checksum, nil
}

func (m *MapDatabase) Close() error {
	return nil
}
