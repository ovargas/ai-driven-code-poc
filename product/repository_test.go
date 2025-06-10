package product

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryRepository_CRUD(t *testing.T) {
	repo := NewInMemoryRepository()
	p := &Product{
		ID:          "test-id",
		Name:        "Test",
		Description: "Desc",
		Price:       10,
	}

	// Create
	if err := repo.Create(p); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Duplicate Create
	if err := repo.Create(p); err == nil {
		t.Fatalf("Expected error on duplicate create")
	}

	// GetByID
	got, err := repo.GetByID("test-id")
	if err != nil || got.Name != "Test" {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Update
	p2 := &Product{Name: "Updated", Description: "D", Price: 20}
	if err := repo.Update("test-id", p2); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	got, _ = repo.GetByID("test-id")
	if got.Name != "Updated" {
		t.Fatalf("Update did not persist")
	}

	// Delete
	if err := repo.Delete("test-id"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := repo.GetByID("test-id"); err == nil {
		t.Fatalf("Expected error after delete")
	}
}

func TestInMemoryRepository_GetAll(t *testing.T) {
	repo := NewInMemoryRepository()
	assert.NoError(t, repo.Create(&Product{ID: "1", Name: "A", Price: 1}))
	assert.NoError(t, repo.Create(&Product{ID: "2", Name: "B", Price: 2}))
	all, err := repo.GetAll()
	if err != nil || len(all) != 2 {
		t.Fatalf("GetAll failed: %v", err)
	}
}
