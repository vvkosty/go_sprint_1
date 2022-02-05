package storage

import (
	"hash/crc32"
	"strconv"
)

type MapDatabase struct {
	urls map[string]string
}

func NewMapDatabase() *MapDatabase {
	var md MapDatabase
	md.urls = make(map[string]string)
	return &md
}

func (m *MapDatabase) Find(id string) string {
	return m.urls[id]
}

func (m *MapDatabase) Save(url string) string {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))
	m.urls[checksum] = url

	return checksum
}
