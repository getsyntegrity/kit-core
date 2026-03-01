// Package errorschain provides utilities for chaining and wrapping errors.
package errorschain

import "go.uber.org/multierr"

// Chain defines an error chain.
type Chain struct {
	returnFirst bool
	errs        []error
}

// ChainOption configures a validation chain at creation time.
type ChainOption func(*Chain)

// New creates a new error chain.
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
		if len(c.errs) == 0 && err != nil {
			c.errs = append(c.errs, err)
		}
		return c
	}
	if err != nil {
		c.errs = append(c.errs, err)
	}
	return c
}

// AddErrorFn add an error to the chain.
func (c *Chain) AddErrorFn(fn func() error) *Chain {
	if c.returnFirst {
		if len(c.errs) == 0 {
			if err := fn(); err != nil {
				c.errs = append(c.errs, err)
			}
		}
		return c
	}
	if err := fn(); err != nil {
		c.errs = append(c.errs, err)
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
	for _, e := range c.errs {
		if e != nil {
			err = multierr.Append(err, e)
		}
	}
	return err
}

// ReturnFirst sets whether a chain should stop on first error.
func ReturnFirst() ChainOption {
	return func(c *Chain) { c.returnFirst = true }
}

// ReturnAll sets whether a chain should return all errors.
func ReturnAll() ChainOption {
	return func(c *Chain) { c.returnFirst = false }
}
