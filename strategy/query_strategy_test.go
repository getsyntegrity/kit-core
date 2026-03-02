package strategy

import (
	"context"
	"errors"
	"testing"

	"github.com/getsyntegrity/kit-core/repository"
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

func TestSimpleQueryStrategy_Execute(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions()
	repo := new(mockQueryableRepo)
	repo.On("Query", ctx, "tid", opts).Return([]string{"a"}, nil)

	s := &SimpleQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	repo.AssertExpectations(t)
}

func TestSimpleQueryStrategy_GetName(t *testing.T) {
	s := &SimpleQueryStrategy[string]{}
	assert.Equal(t, "simple", s.GetName())
}

func TestPaginatedQueryStrategy_Execute(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions().WithPagination(0, 0)
	repo := new(mockQueryableRepo)
	repo.On("Query", ctx, "tid", mock.Anything).Return([]string{"a"}, nil)

	s := &PaginatedQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	assert.Equal(t, 10, opts.Limit) // defaulted
	repo.AssertExpectations(t)
}

func TestPaginatedQueryStrategy_GetName(t *testing.T) {
	s := &PaginatedQueryStrategy[string]{}
	assert.Equal(t, "paginated", s.GetName())
}

func TestSearchQueryStrategy_Execute_NoSearch(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions()
	repo := new(mockQueryableRepo)
	repo.On("Query", ctx, "tid", opts).Return([]string{"a"}, nil)

	s := &SearchQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	repo.AssertExpectations(t)
}

func TestSearchQueryStrategy_Execute_WithSearch(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions().WithSearch("q", []string{"name"})
	repo := new(mockQueryableRepo)
	repo.On("Search", ctx, "tid", "q", []string{"name"}, 10).Return([]string{"b"}, nil)

	s := &SearchQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"b"}, result)
	repo.AssertExpectations(t)
}

func TestSearchQueryStrategy_GetName(t *testing.T) {
	s := &SearchQueryStrategy[string]{}
	assert.Equal(t, "search", s.GetName())
}

func TestFilteredQueryStrategy_Execute_NoFilters(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions()
	repo := new(mockQueryableRepo)
	repo.On("Query", ctx, "tid", opts).Return([]string{"a"}, nil)

	s := &FilteredQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	repo.AssertExpectations(t)
}

func TestFilteredQueryStrategy_Execute_WithFilters(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions().WithFilter("status", "active")
	repo := new(mockQueryableRepo)
	repo.On("FindByFields", ctx, "tid", map[string]interface{}{"status": "active"}).Return([]string{"c"}, nil)

	s := &FilteredQueryStrategy[string]{}
	result, err := s.Execute(ctx, "tid", repo, opts)
	assert.NoError(t, err)
	assert.Equal(t, []string{"c"}, result)
	repo.AssertExpectations(t)
}

func TestFilteredQueryStrategy_GetName(t *testing.T) {
	s := &FilteredQueryStrategy[string]{}
	assert.Equal(t, "filtered", s.GetName())
}

func TestSimpleQueryStrategy_Execute_RepoError(t *testing.T) {
	ctx := context.Background()
	opts := repository.NewQueryOptions()
	repo := new(mockQueryableRepo)
	repo.On("Query", ctx, "tid", opts).Return(nil, errors.New("db error"))

	s := &SimpleQueryStrategy[string]{}
	_, err := s.Execute(ctx, "tid", repo, opts)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	repo.AssertExpectations(t)
}
