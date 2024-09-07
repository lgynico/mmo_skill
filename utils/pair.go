package utils

type Pair[K any, V any] struct {
	key   K
	value V
}

func NewPair[K any, V any](key K, value V) *Pair[K, V] {
	return &Pair[K, V]{
		key:   key,
		value: value,
	}
}

func (p *Pair[K, V]) Key() K {
	return p.key
}

func (p *Pair[K, V]) Value() V {
	return p.value
}

func (p *Pair[K, V]) SetValue(value V) {
	p.value = value
}
