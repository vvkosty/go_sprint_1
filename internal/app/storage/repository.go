package app

type Repository interface {
	Find(id string) string
	Save(url string) string
}
