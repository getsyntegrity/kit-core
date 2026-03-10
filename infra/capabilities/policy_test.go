package capabilities

import (
	"go/build"
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilities(t *testing.T) {
	specs.Describe(t, "capabilities", func(s *specs.Spec) {
		s.When("Allowed", func(s *specs.Spec) {
			s.It("always returns false for any capability", func(ctx *specs.Context) {
				allowed := Allowed(Config{}, Capability("any"))
				assert.False(ctx.T, allowed)
			})
		})
		s.When("package imports", func(s *specs.Spec) {
			s.It("must not depend on fflags", func(ctx *specs.Context) {
				pkg, err := build.Import("github.com/getsyntegrity/kit-core/infra/capabilities", "", 0)
				require.NoError(ctx.T, err)
				for _, imp := range pkg.Imports {
					assert.NotContains(ctx.T, imp, "fflags", "infra must not depend on fflags; infra defines only hard platform rules (defaults, non-bypassable capabilities)")
				}
			})
		})
	})
}
