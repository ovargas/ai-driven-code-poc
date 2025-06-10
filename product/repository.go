package product

import (
	"errors"
	"log"
	"sync"

	"gorm.io/gorm"
)

type Repository interface {
	Create(p *Product) error
	Update(id string, p *Product) error
	Delete(id string) error
	GetByID(id string) (*Product, error)
	GetAll() ([]*Product, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	products map[string]*Product
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		products: make(map[string]*Product),
	}
}

func (r *InMemoryRepository) Create(p *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[p.ID]; exists {
		return errors.New("product already exists")
	}
	r.products[p.ID] = p
	return nil
}

func (r *InMemoryRepository) Update(id string, p *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}
	p.ID = id
	r.products[id] = p
	return nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(r.products, id)
	return nil
}

func (r *InMemoryRepository) GetByID(id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (r *InMemoryRepository) GetAll() ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	products := make([]*Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, p)
	}
	return products, nil
}

// GORM-based MySQL repository implementation
type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	// Auto-migrate the Product model
	err := db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(p *Product) error {
	if err := r.db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) Update(id string, p *Product) error {
	var existing Product
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return errors.New("product not found")
	}
	p.ID = id
	if err := r.db.Model(&existing).Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) Delete(id string) error {
	if err := r.db.Delete(&Product{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) GetByID(id string) (*Product, error) {
	var p Product
	if err := r.db.First(&p, "id = ?", id).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &p, nil
}

func (r *GormRepository) GetAll() ([]*Product, error) {
	var products []*Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
