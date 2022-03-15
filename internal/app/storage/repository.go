package app

type Repository interface {
	Find(id string) (string, error)
	Save(url string, userId string, correlationId string) (string, error)
	List(userId string) map[string]string
	Close() error
}
