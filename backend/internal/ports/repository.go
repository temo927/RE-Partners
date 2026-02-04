package ports

type PackSizeRepository interface {
	GetAllActive() ([]int, error)
	Create(sizes []int) error
}
