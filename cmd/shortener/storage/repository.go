package storage

type Repository interface {
	Find(id string) string
	Save(url string) string
}
