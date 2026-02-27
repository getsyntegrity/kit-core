// Package errorschain provides utilities for chaining and wrapping errors
// with additional context and metadata.
package errorschain

import (
	"go.uber.org/multierr"
)

// Chain defines an error chain.
type Chain struct {
	returnFirst bool
	errs        []error
}

// ChainOption configures a validation chain at creation time.
type ChainOption func(*Chain)

// New creates a new error chain. All errors will be evaluated respectively
// according to their insertion order.
func New(opts ...ChainOption) *Chain {
	chain := &Chain{
		errs: make([]error, 0),
	}

	for _, opt := range opts {
		opt(chain)
	}

	return chain
}

// AddError add an error to the chain.
func (c *Chain) AddError(err error) *Chain {
	if c.returnFirst {
		if len(c.errs) == 0 {
			if err != nil {
				c.errs = append(c.errs, err)
				return c
			}
		}
		return c
	}

	if err != nil {
		c.errs = append(c.errs, err)
		return c
	}

	return c
}

// AddErrorFn add an error to the chain.
func (c *Chain) AddErrorFn(errorFunction func() error) *Chain {
	if c.returnFirst {
		if len(c.errs) == 0 {
			if err := errorFunction(); err != nil {
				c.errs = append(c.errs, err)
				return c
			}
		}
		return c
	}

	if err := errorFunction(); err != nil {
		c.errs = append(c.errs, err)
		return c
	}

	return c
}

// AddErrorFns add a slice of error functions to the chain. Remember the slice order does matter here.
func (c *Chain) AddErrorFns(fn ...func() error) *Chain {
	for _, f := range fn {
		c = c.AddErrorFn(f)
	}
	return c
}

// Error returns the error.
func (c *Chain) Error() error {
	if c.returnFirst {
		if len(c.errs) == 0 {
			return nil
		}
		return c.errs[0]
	}

	var err error
	for _, errorItem := range c.errs {
		if errorItem != nil {
			err = multierr.Append(err, errorItem)
		}
	}
	return err
}

// ReturnFirst sets whether a chain should stop validation on first error.
func ReturnFirst() ChainOption {
	return func(c *Chain) { c.returnFirst = true }
}

// ReturnAll sets whether a chain should return all errors.
func ReturnAll() ChainOption {
	return func(c *Chain) { c.returnFirst = false }
}
