package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestEntity struct {
	ID   string
	Name string
	Age  int
}

func TestMemoryRepository(t *testing.T) {
	repo := NewMemoryRepository[TestEntity]()
	ctx := context.Background()
	tenantID := "test-tenant"

	entity := TestEntity{ID: "1", Name: "John", Age: 30}
	err := repo.SaveWithID(ctx, "1", entity)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, tenantID, "1")
	require.NoError(t, err)
	assert.Equal(t, entity, retrieved)

	exists, err := repo.Exists(ctx, tenantID, "1")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.Exists(ctx, tenantID, "2")
	require.NoError(t, err)
	assert.False(t, exists)

	count, err := repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	entities, err := repo.List(ctx, tenantID, 0, 10)
	require.NoError(t, err)
	assert.Len(t, entities, 1)
	assert.Equal(t, entity, entities[0])

	err = repo.Delete(ctx, tenantID, "1")
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, tenantID, "1")
	assert.Error(t, err)

	count, err = repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMemoryRepositoryMultipleEntities(t *testing.T) {
	repo := NewMemoryRepository[TestEntity]()
	ctx := context.Background()
	tenantID := "test-tenant"

	entities := []TestEntity{
		{ID: "1", Name: "John", Age: 30},
		{ID: "2", Name: "Jane", Age: 25},
		{ID: "3", Name: "Bob", Age: 35},
	}

	for _, entity := range entities {
		err := repo.SaveWithID(ctx, entity.ID, entity)
		require.NoError(t, err)
	}

	count, err := repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	all, err := repo.List(ctx, tenantID, 0, 10)
	require.NoError(t, err)
	assert.Len(t, all, 3)

	limited, err := repo.List(ctx, tenantID, 0, 2)
	require.NoError(t, err)
	assert.Len(t, limited, 2)

	offset, err := repo.List(ctx, tenantID, 1, 2)
	require.NoError(t, err)
	assert.Len(t, offset, 2)
}

func TestMemoryRepositoryClear(t *testing.T) {
	repo := NewMemoryRepository[TestEntity]()
	ctx := context.Background()
	tenantID := "test-tenant"

	entity := TestEntity{ID: "1", Name: "John", Age: 30}
	err := repo.SaveWithID(ctx, "1", entity)
	require.NoError(t, err)

	count, err := repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	repo.Clear()

	count, err = repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMemoryRepositoryGetAll(t *testing.T) {
	repo := NewMemoryRepository[TestEntity]()
	ctx := context.Background()

	entities := []TestEntity{
		{ID: "1", Name: "John", Age: 30},
		{ID: "2", Name: "Jane", Age: 25},
	}

	for _, entity := range entities {
		err := repo.SaveWithID(ctx, entity.ID, entity)
		require.NoError(t, err)
	}

	all := repo.GetAll()
	assert.Len(t, all, 2)
}

func TestMemoryRepositoryConcurrency(t *testing.T) {
	repo := NewMemoryRepository[TestEntity]()
	ctx := context.Background()
	tenantID := "test-tenant"

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			entity := TestEntity{ID: string(rune(id)), Name: "Test", Age: id}
			err := repo.SaveWithID(ctx, string(rune(id)), entity)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	count, err := repo.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}
