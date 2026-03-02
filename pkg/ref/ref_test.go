package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefTypes(t *testing.T) {
	require.Equal(t, "test", SafeDeref(Of("test")))
	require.Equal(t, int32(123), SafeDeref(Of(int32(123))))
	require.Equal(t, int64(123), SafeDeref(Of(int64(123))))
}

func TestOf(t *testing.T) {
	p := Of("v")
	require.NotNil(t, p)
	assert.Equal(t, "v", *p)
}

func TestSafeDeref_Nil(t *testing.T) {
	var p *string
	assert.Equal(t, "", SafeDeref(p))
}

func TestNilIfNil(t *testing.T) {
	assert.Nil(t, NilIfNil[string](nil))
	v := "x"
	assert.Equal(t, "x", NilIfNil(&v))
}

func TestDerefSlice(t *testing.T) {
	assert.Nil(t, DerefSlice[string](nil))
	assert.Empty(t, DerefSlice([]*string{}))
	a, b := "a", "b"
	ptrs := []*string{&a, nil, &b}
	out := DerefSlice(ptrs)
	assert.Equal(t, []string{"a", "b"}, out)
}
