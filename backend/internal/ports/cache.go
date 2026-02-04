package ports

type Cache interface {
	Get(key string) ([]int, error)
	Set(key string, value []int, ttl int) error
	Delete(key string) error
}
