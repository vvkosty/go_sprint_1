package storage

import (
	"hash/crc32"
	"strconv"
)

var DB *MapDatabase

func init() {
	if DB == nil {
		DB = newMapDatabase()
	}
}

type MapDatabase struct {
	urls map[string]string
}

func newMapDatabase() *MapDatabase {
	var md MapDatabase
	md.urls = make(map[string]string)
	return &md
}

func (m MapDatabase) Find(id string) string {
	return m.urls[id]
}

func (m MapDatabase) Save(url string) string {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))
	m.urls[checksum] = url

	return checksum
}
