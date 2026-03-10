package idgen

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestIDGen(t *testing.T) {
	specs.Describe(t, "idgen", func(s *specs.Spec) {
		s.When("IsValidULID", func(s *specs.Spec) {
			s.It("returns false for empty string", func(ctx *specs.Context) {
				assert.False(ctx.T, IsValidULID(""))
			})
			s.It("returns true for valid ULID", func(ctx *specs.Context) {
				assert.True(ctx.T, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
			})
			s.It("returns false for invalid", func(ctx *specs.Context) {
				assert.False(ctx.T, IsValidULID("invalid"))
			})
			s.It("returns false for too short", func(ctx *specs.Context) {
				assert.False(ctx.T, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FA"))
			})
		})

		s.When("ParseULID", func(s *specs.Spec) {
			s.It("parses valid ULID and returns string", func(ctx *specs.Context) {
				u, err := ParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
			})
			s.It("returns error for bad input", func(ctx *specs.Context) {
				_, err := ParseULID("bad")
				assert.Error(ctx.T, err)
			})
		})

		s.When("MustParseULID", func(s *specs.Spec) {
			s.It("returns ULID for valid input", func(ctx *specs.Context) {
				u := MustParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
				assert.Equal(ctx.T, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
			})
			s.It("panics on invalid input", func(ctx *specs.Context) {
				assert.Panics(ctx.T, func() { MustParseULID("invalid") })
			})
		})

		s.When("NewULID", func(s *specs.Spec) {
			s.It("returns valid ULID string", func(ctx *specs.Context) {
				u := NewULID()
				assert.True(ctx.T, IsValidULID(u))
			})
			s.It("returns string of length 26", func(ctx *specs.Context) {
				u := NewULID()
				assert.Len(ctx.T, u, 26)
			})
			s.It("multiple calls generate unique IDs", func(ctx *specs.Context) {
				u1, u2 := NewULID(), NewULID()
				assert.NotEqual(ctx.T, u1, u2)
			})
		})

		s.When("MustNewULID", func(s *specs.Spec) {
			s.It("returns valid ULID and does not panic", func(ctx *specs.Context) {
				u := MustNewULID()
				assert.True(ctx.T, IsValidULID(u))
				assert.Len(ctx.T, u, 26)
			})
		})

		s.When("GenerateID", func(s *specs.Spec) {
			s.It("returns ID that can be converted to string", func(ctx *specs.Context) {
				id := GenerateID()
				str := ToString(id)
				assert.NotEmpty(ctx.T, str)
			})
			s.It("multiple calls generate unique IDs when both succeed", func(ctx *specs.Context) {
				id1, id2 := GenerateID(), GenerateID()
				if id1 != 0 && id2 != 0 {
					assert.NotEqual(ctx.T, id1, id2)
				}
			})
		})

		s.When("ToString", func(s *specs.Spec) {
			s.It("converts zero ID to string", func(ctx *specs.Context) {
				assert.Equal(ctx.T, "0", ToString(0))
			})
			s.It("converts non-zero ID to string", func(ctx *specs.Context) {
				assert.Equal(ctx.T, "123", ToString(123))
			})
			s.It("roundtrip: ID to string to parsed equals same value", func(ctx *specs.Context) {
				id := ID(12345)
				str := ToString(id)
				parsed, err := StringToSonyFlake(str)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, uint64(id), parsed)
			})
		})

		s.When("StringToSonyFlake", func(s *specs.Spec) {
			s.It("parses valid decimal string", func(ctx *specs.Context) {
				v, err := StringToSonyFlake("0")
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, uint64(0), v)
			})
			s.It("parses large valid string", func(ctx *specs.Context) {
				v, err := StringToSonyFlake("18446744073709551615")
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, uint64(18446744073709551615), v)
			})
			s.It("returns error for invalid string", func(ctx *specs.Context) {
				_, err := StringToSonyFlake("abc")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "failed to parse")
			})
			s.It("returns error for negative-looking input", func(ctx *specs.Context) {
				_, err := StringToSonyFlake("-1")
				assert.Error(ctx.T, err)
			})
			s.It("returns error for empty string", func(ctx *specs.Context) {
				_, err := StringToSonyFlake("")
				assert.Error(ctx.T, err)
			})
		})
	})
}
