package app

type Repository interface {
	Find(id string) (string, error)
	Save(url string, userID string) (string, error)
	List(userID string) map[string]string
	DeleteBatchByChecksums(checksums []string) error
	Close() error
}
