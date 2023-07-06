package cache

type Cache interface {
	Get(key string) (value Value, ok bool)
	Add(key string, value Value)
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}
