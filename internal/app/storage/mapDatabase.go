package app

import (
	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
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

func (m *MapDatabase) Save(url string, userID string) (string, error) {
	checksum := helpers.GenerateChecksum(url)
	m.urls[checksum] = url
	m.usersUrls[userID] = append(m.usersUrls[userID], checksum)

	return checksum, nil
}

func (m *MapDatabase) List(userID string) map[string]string {
	result := make(map[string]string)

	if urls, found := m.usersUrls[userID]; found {
		for _, checksum := range urls {
			result[checksum] = m.urls[checksum]
		}
	}
	return result
}

func (m *MapDatabase) Close() error {
	return nil
}
