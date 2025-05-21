package bufferspool_test

import (
	"bytes"
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
	bufferspool "github.com/sv-tools/buffers-pool"
)

func TestGlobalPool(t *testing.T) {
	t.Parallel()

	b1 := bufferspool.Get()
	b2 := bufferspool.Get()
	require.NotEqual(t, unsafe.Pointer(b1), unsafe.Pointer(b2)) //nolint: gosec

	b1.WriteString("foo")
	bufferspool.Put(b1)
	b2 = bufferspool.Get()
	require.Zero(t, b2.Len())

	bufferspool.Do(func(b *bytes.Buffer) {
		require.Zero(t, b.Len())
		b.WriteString("foo")
		require.Equal(t, "foo", b.String())
	})
}

func TestCustomPool(t *testing.T) {
	t.Parallel()

	p := bufferspool.New()
	b1 := p.Get()
	b2 := p.Get()
	require.NotEqual(t, unsafe.Pointer(b1), unsafe.Pointer(b2)) //nolint: gosec

	b1.WriteString("foo")
	p.Put(b1)
	b2 = p.Get()
	require.Zero(t, b2.Len())

	p.Do(func(b *bytes.Buffer) {
		require.Zero(t, b.Len())
		b.WriteString("foo")
		require.Equal(t, "foo", b.String())
	})
}

func TestSafety(t *testing.T) {
	t.Parallel()

	p := bufferspool.New()
	b1 := p.Get()
	b1.WriteString("foo42")
	b1AsString := b1.String()
	b1AsBytes := b1.Bytes()
	b1Data, err := io.ReadAll(b1)
	require.NoError(t, err)

	require.Equal(t, "foo42", b1AsString)
	require.Equal(t, []byte("foo42"), b1AsBytes, string(b1AsBytes))
	require.Equal(t, []byte("foo42"), b1Data, string(b1Data))

	// put the same buffer to the poll several times to increase the chance that the same object will be returned
	// this is the only way to test the Pools if we need already used object
	for range 10 {
		p.Put(b1)
	}
	b2 := p.Get()
	defer p.Put(b2)

	require.Equal(t, unsafe.Pointer(b1), unsafe.Pointer(b2)) //nolint: gosec
	require.Empty(t, b2.Bytes())

	b2.WriteString("bar")
	b2AsString := b2.String()
	b2AsBytes := b2.Bytes()
	b2Data, err := io.ReadAll(b2)
	require.NoError(t, err)

	require.Equal(t, "bar", b2AsString)
	require.Equal(t, []byte("bar"), b2AsBytes, string(b2AsBytes))
	require.Equal(t, []byte("bar"), b2Data, string(b2Data))

	// usage of string is safe and should not be changed
	require.Equal(t, "foo42", b1AsString)
	// usage of bytes is unsafe, the b1AsBytes must contain `bar42` instead of `foo42` or `bar`
	require.Equal(t, []byte("bar42"), b1AsBytes, string(b1AsBytes))
	// usage of Reader is safe and should not be changed
	require.Equal(t, []byte("foo42"), b1Data, string(b1Data))
}
