package app

import (
	"hash/crc32"
	"strconv"
)

type MapDatabase struct {
	urls      map[string]string
	usersUrls map[string][]string
}

func NewMapStorage() *MapDatabase {
	var md MapDatabase
	md.urls = make(map[string]string)
	md.usersUrls = make(map[string][]string)
	return &md
}

func (m *MapDatabase) Find(id string) (string, error) {
	return m.urls[id], nil
}

func (m *MapDatabase) Save(url string, userId string, correlationId string) (string, error) {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))
	m.urls[checksum] = url
	m.usersUrls[userId] = append(m.usersUrls[userId], checksum)

	return checksum, nil
}

func (m *MapDatabase) List(userId string) map[string]string {
	result := make(map[string]string)

	if urls, found := m.usersUrls[userId]; found {
		for _, checksum := range urls {
			result[checksum] = m.urls[checksum]
		}
	}
	return result
}

func (m *MapDatabase) Close() error {
	return nil
}
