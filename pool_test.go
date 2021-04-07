package bufferspool_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	bufferspool "github.com/sv-tools/buffers-pool"
)

func TestGlobalPool(t *testing.T) {
	b1 := bufferspool.Get()
	b2 := bufferspool.Get()
	require.NotSame(t, b1, b2)

	b1.WriteString("foo")
	bufferspool.Put(b1)
	b2 = bufferspool.Get()
	require.Zero(t, b2.Len())
}

func TestCustomPool(t *testing.T) {
	p := bufferspool.New()
	b1 := p.Get()
	b2 := p.Get()
	require.NotSame(t, b1, b2)

	b1.WriteString("foo")
	p.Put(b1)
	b2 = p.Get()
	require.Zero(t, b2.Len())
}
