package evict

type Entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

type Evicter interface {
	Get(key string) (value Value, err error)
	Add(key string, value Value)
}
