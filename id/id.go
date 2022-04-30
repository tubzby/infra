package id

// Generater .
type Generater interface {
	NextID() int64
	NextIDStr() string
}

// New create a id generater
func New(node int64) Generater {
	return newSnowFlakeID(node)
}
