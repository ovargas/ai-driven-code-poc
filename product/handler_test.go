package product

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupHandler() *Handler {
	repo := NewInMemoryRepository()
	return NewHandler(repo)
}

func TestHandler_CreateProduct(t *testing.T) {
	handler := setupHandler()
	reqBody := `{"name":"Test","description":"Desc","price":10}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var p Product
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if p.Name != "Test" {
		t.Fatalf("unexpected product: %+v", p)
	}
}

func TestHandler_CreateProduct_Validation(t *testing.T) {
	handler := setupHandler()
	reqBody := `{"name":"","price":0}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_GetAllProducts(t *testing.T) {
	handler := setupHandler()
	// Create a product first
	reqBody := `{"name":"Test","description":"Desc","price":10}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)

	req = httptest.NewRequest(http.MethodGet, "/product", nil)
	w = httptest.NewRecorder()
	handler.handleProducts(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(products) == 0 {
		t.Fatalf("expected at least one product")
	}
}

func TestHandler_ProductByID(t *testing.T) {
	handler := setupHandler()
	// Create a product
	reqBody := `{"name":"Test","description":"Desc","price":10}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)
	var p Product
	json.NewDecoder(w.Result().Body).Decode(&p)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, "/product/"+p.ID, nil)
	w = httptest.NewRecorder()
	handler.handleProductByID(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got Product
	json.NewDecoder(resp.Body).Decode(&got)
	if got.ID != p.ID {
		t.Fatalf("expected ID %s, got %s", p.ID, got.ID)
	}
}

func TestHandler_UpdateProduct(t *testing.T) {
	handler := setupHandler()
	// Create a product
	reqBody := `{"name":"Test","description":"Desc","price":10}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)
	var p Product
	json.NewDecoder(w.Result().Body).Decode(&p)

	// Update
	updateBody := `{"name":"Updated","description":"New","price":20}`
	req = httptest.NewRequest(http.MethodPatch, "/product/"+p.ID, strings.NewReader(updateBody))
	w = httptest.NewRecorder()
	handler.handleProductByID(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var updated Product
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Name != "Updated" {
		t.Fatalf("expected name Updated, got %s", updated.Name)
	}
}

func TestHandler_DeleteProduct(t *testing.T) {
	handler := setupHandler()
	// Create a product
	reqBody := `{"name":"Test","description":"Desc","price":10}`
	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.handleProducts(w, req)
	var p Product
	json.NewDecoder(w.Result().Body).Decode(&p)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/product/"+p.ID, nil)
	w = httptest.NewRecorder()
	handler.handleProductByID(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}
