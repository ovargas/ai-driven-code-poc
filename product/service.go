package product

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(p *Product) (*Product, error) {
	if err := validateProductInput(p); err != nil {
		return nil, err
	}
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate UUID")
	}
	p.ID = uid.String()
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(id string, p *Product) (*Product, error) {
	if err := validateProductInput(p); err != nil {
		return nil, err
	}
	if err := s.repo.Update(id, p); err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) GetByID(id string) (*Product, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetAll() ([]*Product, error) {
	return s.repo.GetAll()
}

func validateProductInput(p *Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return &validationError{"Name is required and cannot be empty"}
	}
	if len(p.Name) > 128 {
		return &validationError{"Name cannot exceed 128 characters"}
	}
	if p.Price <= 0 {
		return &validationError{"Price must be greater than 0"}
	}
	return nil
}

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}
