package product

import (
	"testing"
)

func TestService_Create_Validation(t *testing.T) {
	s := NewService(NewInMemoryRepository())

	// Empty name
	_, err := s.Create(&Product{Name: "", Price: 1})
	if err == nil {
		t.Error("expected error for empty name")
	}

	// Name too long
	longName := ""
	for i := 0; i < 129; i++ {
		longName += "a"
	}
	_, err = s.Create(&Product{Name: longName, Price: 1})
	if err == nil {
		t.Error("expected error for long name")
	}

	// Price zero
	_, err = s.Create(&Product{Name: "Valid", Price: 0})
	if err == nil {
		t.Error("expected error for zero price")
	}

	// Price negative
	_, err = s.Create(&Product{Name: "Valid", Price: -1})
	if err == nil {
		t.Error("expected error for negative price")
	}

	// Valid
	p, err := s.Create(&Product{Name: "Valid", Price: 1})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_CRUD(t *testing.T) {
	s := NewService(NewInMemoryRepository())
	p, _ := s.Create(&Product{Name: "A", Price: 1})

	// GetByID
	got, err := s.GetByID(p.ID)
	if err != nil || got.Name != "A" {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Update
	updated, err := s.Update(p.ID, &Product{Name: "B", Price: 2})
	if err != nil || updated.Name != "B" {
		t.Fatalf("Update failed: %v", err)
	}

	// Delete
	if err := s.Delete(p.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := s.GetByID(p.ID); err == nil {
		t.Fatalf("Expected error after delete")
	}
}

func TestService_GetAll(t *testing.T) {
	s := NewService(NewInMemoryRepository())
	s.Create(&Product{Name: "A", Price: 1})
	s.Create(&Product{Name: "B", Price: 2})
	all, err := s.GetAll()
	if err != nil || len(all) != 2 {
		t.Fatalf("GetAll failed: %v", err)
	}
}
