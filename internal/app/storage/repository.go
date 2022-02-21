package app

type Repository interface {
	Find(id string) (string, error)
	Save(url string) (string, error)
	Close() error
}
