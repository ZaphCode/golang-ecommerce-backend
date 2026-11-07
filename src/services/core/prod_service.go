package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type prodService struct {
	prodRepo domain.ProductRepository
	catRepo  domain.CategoryRepository
}

func NewProductService(
	prodRepo domain.ProductRepository,
	catRepo domain.CategoryRepository,
) domain.ProductService {
	return &prodService{
		prodRepo: prodRepo,
		catRepo:  catRepo,
	}
}

func (s *prodService) Create(prod *domain.Product) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	c, err := s.catRepo.FindByField("Name", prod.Category)

	if err != nil {
		return err
	}

	if c == nil {
		return fmt.Errorf("category %q does not exist", prod.Category)
	}

	prod.ID = ID
	prod.CreatedAt = time.Now().Unix()
	prod.UpdatedAt = time.Now().Unix()

	s.prodRepo.Save(prod)

	return nil
}

func (s *prodService) GetAll() ([]domain.Product, error) {
	return s.prodRepo.Find()
}

func (s *prodService) GetByID(ID uuid.UUID) (*domain.Product, error) {
	return s.prodRepo.FindByID(ID)
}

func (s *prodService) GetLatestProds(lim ...int) ([]domain.Product, error) {
	prods, err := s.prodRepo.FindOrderBy("CreatedAt", "DESC")
	if err != nil {
		return nil, err
	}

	if len(lim) <= 0 || lim[0] >= len(prods) {
		return prods, nil
	}
	return prods[:lim[0]], nil
}

func (s *prodService) GetByTags(tags ...string) ([]domain.Product, error) {
	return s.prodRepo.FindWhere("Tags", "array-contains-any", tags)
}

func (s *prodService) GetByCategory(c string) ([]domain.Product, error) {
	return s.prodRepo.FindWhere("Category", "==", c)
}

func (s *prodService) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	p, err := s.prodRepo.FindByID(ID)

	if err != nil || p == nil {
		return fmt.Errorf("error getting product")
	}

	if v, ok := uf["Category"]; ok {
		c, err := s.catRepo.FindByField("Name", v)

		if err != nil {
			return err
		}

		if c == nil {
			return fmt.Errorf("invalid category field. that category does not exist")
		}
	}

	return s.prodRepo.Update(ID, uf)
}

func (s *prodService) SetAvailable(ID uuid.UUID, avl bool) error {
	p, err := s.prodRepo.FindByID(ID)

	if err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("product not found")
	}

	return s.prodRepo.UpdateField(ID, "Available", avl)
}

func (s *prodService) Delete(ID uuid.UUID) error {
	p, err := s.prodRepo.FindByID(ID)

	if err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("product not found")
	}

	return s.prodRepo.Remove(ID)
}

func (s *prodService) CalculateTotalPrice(ops []domain.OrderProduct) (int64, error) {
	if len(ops) == 0 || ops == nil {
		return 0, fmt.Errorf("missing products")
	}

	var total float32 = 0

	for _, op := range ops {
		p, err := s.prodRepo.FindByID(op.ID)

		if err != nil {
			return 0, err
		}

		if p == nil {
			return 0, fmt.Errorf("product %s not found", op.ID.String())
		}

		var price float32 = float32(p.Price) * ((float32(100) - float32(p.DiscountRate)) / 100) * float32(op.Quantity)

		total += price
	}

	return int64(total), nil
}
