package errorschain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorsChain(t *testing.T) {
	t.Run("With ReturnFirst", func(t *testing.T) {
		e1 := errors.New("err1")
		e2 := errors.New("err2")
		e3 := errors.New("err3")

		chain := New(ReturnFirst()).AddError(e1).AddError(e2).AddError(e3)
		actual := chain.Error()
		require.True(t, errors.Is(actual, e1))
	})

	t.Run("With Error", func(t *testing.T) {
		chain := New(ReturnFirst()).AddError(nil)
		actual := chain.Error()
		require.NoError(t, actual)
	})

	t.Run("With AddErrorFn ReturnFirst", func(t *testing.T) {
		var (
			calledFn1 = false
			calledFn2 = false
			calledFn3 = false
		)

		fn1 := func() error { calledFn1 = true; return errors.New("err1") }
		fn2 := func() error { calledFn2 = true; return errors.New("err2") }
		fn3 := func() error { calledFn3 = true; return errors.New("err3") }

		chain := New(ReturnFirst()).
			AddErrorFn(fn1).
			AddErrorFn(fn2).
			AddErrorFn(fn3)
		actual := chain.Error()

		require.EqualError(t, actual, "err1")
		require.True(t, calledFn1)
		require.False(t, calledFn2)
		require.False(t, calledFn3)
	})

	t.Run("With AddErrorFns ReturnFirst", func(t *testing.T) {
		var (
			calledFn1 = false
			calledFn2 = false
			calledFn3 = false
		)

		fn1 := func() error { calledFn1 = true; return errors.New("err1") }
		fn2 := func() error { calledFn2 = true; return errors.New("err2") }
		fn3 := func() error { calledFn3 = true; return errors.New("err3") }

		chain := New(ReturnFirst()).AddErrorFns(fn1, fn2, fn3)
		actual := chain.Error()

		require.EqualError(t, actual, "err1")
		require.True(t, calledFn1)
		require.False(t, calledFn2)
		require.False(t, calledFn3)
	})

	t.Run("With AddErrorFn ReturnAll", func(t *testing.T) {
		var (
			calledFn1 = false
			calledFn2 = false
			calledFn3 = false
		)

		fn1 := func() error { calledFn1 = true; return errors.New("err1") }
		fn2 := func() error { calledFn2 = true; return errors.New("err2") }
		fn3 := func() error { calledFn3 = true; return nil }

		chain := New(ReturnAll()).
			AddErrorFn(fn1).
			AddErrorFn(fn2).
			AddErrorFn(fn3)
		actual := chain.Error()

		require.EqualError(t, actual, "err1; err2")
		require.True(t, calledFn1)
		require.True(t, calledFn2)
		require.True(t, calledFn3)
	})

	t.Run("With ReturnAll", func(t *testing.T) {
		e1 := errors.New("err1")
		e2 := errors.New("err2")
		e3 := errors.New("err3")

		chain := New(ReturnAll()).
			AddError(e1).
			AddError(e2).
			AddError(e3).
			AddError(nil)
		actual := chain.Error()
		require.EqualError(t, actual, "err1; err2; err3")
	})
}
