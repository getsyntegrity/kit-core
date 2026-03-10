package strategy

import (
	"context"
	"errors"
	"testing"

	"github.com/getsyntegrity/kit-core/repository"
	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockQueryableRepo is a testify mock for QueryableRepository for use in strategy tests.
type mockQueryableRepo struct {
	mock.Mock
}

func (m *mockQueryableRepo) Save(ctx context.Context, entity string) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *mockQueryableRepo) GetByID(ctx context.Context, tenantID, id string) (string, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *mockQueryableRepo) Exists(ctx context.Context, tenantID, id string) (bool, error) {
	args := m.Called(ctx, tenantID, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockQueryableRepo) List(ctx context.Context, tenantID string, offset, limit int) ([]string, error) {
	args := m.Called(ctx, tenantID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) Count(ctx context.Context, tenantID string) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockQueryableRepo) Delete(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

func (m *mockQueryableRepo) Query(ctx context.Context, tenantID string, options *repository.QueryOptions) ([]string, error) {
	args := m.Called(ctx, tenantID, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) QueryWithCount(ctx context.Context, tenantID string, options *repository.QueryOptions) ([]string, int64, error) {
	args := m.Called(ctx, tenantID, options)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]string), args.Get(1).(int64), args.Error(2)
}

func (m *mockQueryableRepo) QueryWithCursor(ctx context.Context, tenantID string, pagination *repository.CursorPagination, filters map[string]interface{}) (*repository.CursorResult[string], error) {
	args := m.Called(ctx, tenantID, pagination, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CursorResult[string]), args.Error(1)
}

func (m *mockQueryableRepo) FindByField(ctx context.Context, tenantID, field string, value interface{}) ([]string, error) {
	args := m.Called(ctx, tenantID, field, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) FindByFields(ctx context.Context, tenantID string, filters map[string]interface{}) ([]string, error) {
	args := m.Called(ctx, tenantID, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) Search(ctx context.Context, tenantID, term string, fields []string, limit int) ([]string, error) {
	args := m.Called(ctx, tenantID, term, fields, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) GetByIDs(ctx context.Context, tenantID string, ids []string) ([]string, error) {
	args := m.Called(ctx, tenantID, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockQueryableRepo) BulkSave(ctx context.Context, entities []string) error {
	args := m.Called(ctx, entities)
	return args.Error(0)
}

func (m *mockQueryableRepo) BulkDelete(ctx context.Context, tenantID string, ids []string) error {
	args := m.Called(ctx, tenantID, ids)
	return args.Error(0)
}

func (m *mockQueryableRepo) SoftDelete(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

func (m *mockQueryableRepo) Restore(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

func TestStrategy(t *testing.T) {
	specs.Describe(t, "strategy", func(s *specs.Spec) {
		s.When("SimpleQueryStrategy", func(s *specs.Spec) {
			s.It("Execute returns repo query result", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions()
				repo := new(mockQueryableRepo)
				repo.On("Query", bg, "tid", opts).Return([]string{"a"}, nil)

				strat := &SimpleQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"a"}, result)
				repo.AssertExpectations(ctx.T)
			})
			s.It("Execute returns repo error", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions()
				repo := new(mockQueryableRepo)
				repo.On("Query", bg, "tid", opts).Return(nil, errors.New("db error"))

				strat := &SimpleQueryStrategy[string]{}
				_, err := strat.Execute(bg, "tid", repo, opts)
				assert.Error(ctx.T, err)
				assert.Equal(ctx.T, "db error", err.Error())
				repo.AssertExpectations(ctx.T)
			})
			s.It("GetName returns simple", func(ctx *specs.Context) {
				strat := &SimpleQueryStrategy[string]{}
				assert.Equal(ctx.T, "simple", strat.GetName())
			})
		})

		s.When("PaginatedQueryStrategy", func(s *specs.Spec) {
			s.It("Execute uses repo and default limit", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions().WithPagination(0, 0)
				repo := new(mockQueryableRepo)
				repo.On("Query", bg, "tid", mock.Anything).Return([]string{"a"}, nil)

				strat := &PaginatedQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"a"}, result)
				assert.Equal(ctx.T, 10, opts.Limit)
				repo.AssertExpectations(ctx.T)
			})
			s.It("GetName returns paginated", func(ctx *specs.Context) {
				strat := &PaginatedQueryStrategy[string]{}
				assert.Equal(ctx.T, "paginated", strat.GetName())
			})
		})

		s.When("SearchQueryStrategy", func(s *specs.Spec) {
			s.It("Execute with no search uses Query", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions()
				repo := new(mockQueryableRepo)
				repo.On("Query", bg, "tid", opts).Return([]string{"a"}, nil)

				strat := &SearchQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"a"}, result)
				repo.AssertExpectations(ctx.T)
			})
			s.It("Execute with search uses Search", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions().WithSearch("q", []string{"name"})
				repo := new(mockQueryableRepo)
				repo.On("Search", bg, "tid", "q", []string{"name"}, 10).Return([]string{"b"}, nil)

				strat := &SearchQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"b"}, result)
				repo.AssertExpectations(ctx.T)
			})
			s.It("GetName returns search", func(ctx *specs.Context) {
				strat := &SearchQueryStrategy[string]{}
				assert.Equal(ctx.T, "search", strat.GetName())
			})
		})

		s.When("FilteredQueryStrategy", func(s *specs.Spec) {
			s.It("Execute with no filters uses Query", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions()
				repo := new(mockQueryableRepo)
				repo.On("Query", bg, "tid", opts).Return([]string{"a"}, nil)

				strat := &FilteredQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"a"}, result)
				repo.AssertExpectations(ctx.T)
			})
			s.It("Execute with filters uses FindByFields", func(ctx *specs.Context) {
				bg := context.Background()
				opts := repository.NewQueryOptions().WithFilter("status", "active")
				repo := new(mockQueryableRepo)
				repo.On("FindByFields", bg, "tid", map[string]interface{}{"status": "active"}).Return([]string{"c"}, nil)

				strat := &FilteredQueryStrategy[string]{}
				result, err := strat.Execute(bg, "tid", repo, opts)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, []string{"c"}, result)
				repo.AssertExpectations(ctx.T)
			})
			s.It("GetName returns filtered", func(ctx *specs.Context) {
				strat := &FilteredQueryStrategy[string]{}
				assert.Equal(ctx.T, "filtered", strat.GetName())
			})
		})
	})
}
