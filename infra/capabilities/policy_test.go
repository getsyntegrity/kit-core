package capabilities

import (
	"go/build"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowed_AlwaysReturnsFalse(t *testing.T) {
	// Kit does not define any capabilities; Allowed is fail-closed for all inputs.
	allowed := Allowed(Config{}, Capability("any"))
	assert.False(t, allowed)
}

// TestInfraMustNotDependOnFflags ensures this package does not import fflags.
// Infra defines only hard platform rules; it must not call feature flags or depend on runtime providers.
func TestInfraMustNotDependOnFflags(t *testing.T) {
	pkg, err := build.Import("github.com/getsyntegrity/kit-core/infra/capabilities", "", 0)
	require.NoError(t, err)
	for _, imp := range pkg.Imports {
		assert.NotContains(t, imp, "fflags", "infra must not depend on fflags; infra defines only hard platform rules (defaults, non-bypassable capabilities)")
	}
}
