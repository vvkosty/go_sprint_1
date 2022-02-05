package storage

type Repository interface {
	Get(id string) string
	Save(url string)
}
