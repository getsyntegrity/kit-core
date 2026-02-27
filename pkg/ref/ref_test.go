package ref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefTypes(t *testing.T) {
	require.Equal(t, "test", SafeDeref(Of("test")))
	require.Equal(t, int32(123), SafeDeref(Of(int32(123))))
	require.Equal(t, int64(123), SafeDeref(Of(int64(123))))
}
