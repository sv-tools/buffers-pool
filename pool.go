package bufferspool

import (
	"bytes"
	"sync"
)

var globalPool Pool

func init() {
	globalPool = New()
}

// Pool is an interface to work with the Buffers Pool
type Pool interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

type pool struct {
	pool sync.Pool
}

// New returns a new object of the Buffers Pool
func New() Pool {
	p := pool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	return &p
}

// Get returns a buffer from the pool or creates a new one
func Get() *bytes.Buffer {
	return globalPool.Get()
}

func (p *pool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put resets and puts back a given buffer to the pool
func Put(buf *bytes.Buffer) {
	globalPool.Put(buf)
}

func (p *pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}
