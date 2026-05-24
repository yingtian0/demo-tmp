package memory

import (
	"fmt"
	"sync/atomic"
)

type IncrementalIDGenerator struct {
	prefix  string
	counter atomic.Uint64
}

func NewIncrementalIDGenerator(prefix string) *IncrementalIDGenerator {
	return &IncrementalIDGenerator{prefix: prefix}
}

func (g *IncrementalIDGenerator) NewID() string {
	n := g.counter.Add(1)
	return fmt.Sprintf("%s%d", g.prefix, n)
}
