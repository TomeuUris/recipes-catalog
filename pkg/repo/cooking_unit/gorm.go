package cooking_unit

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

// Database model
type CookingUnit struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func (i *CookingUnit) ToEntity() *entity.CookingUnit {
	return &entity.CookingUnit{
		ID:   int64(i.ID),
		Name: i.Name,
	}
}

func (i *CookingUnit) FromEntity(unit *entity.CookingUnit) {
	i.ID = uint(unit.ID)
	i.Name = unit.Name
}

// Repository implementation
type RepoGorm struct {
	db *gorm.DB
}

// Utility functions
func NewGormRepo(db *gorm.DB) *RepoGorm {
	return &RepoGorm{
		db: db,
	}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&CookingUnit{})
}

// CRUD functions
func (r *RepoGorm) FindByID(ctx context.Context, id int) (*entity.CookingUnit, error) {
	unit := &CookingUnit{}
	if err := r.db.WithContext(ctx).First(unit, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return unit.ToEntity(), nil
}

func (r *RepoGorm) FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.CookingUnit, error) {
	var units []*CookingUnit
	if err := r.db.WithContext(ctx).Where(f).Find(&units).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.CookingUnit, len(units))
	for i, unit := range units {
		result[i] = unit.ToEntity()
	}

	return result, nil
}

func (r *RepoGorm) CountByFilter(f *FindFilter) (int, error) {
	var count int64
	if err := r.db.Model(&CookingUnit{}).Where(f).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RepoGorm) Add(ctx context.Context, unit *entity.CookingUnit) error {
	i := &CookingUnit{}
	i.FromEntity(unit)
	err := r.db.WithContext(ctx).Create(i).Error
	if err != nil {
		return err
	}
	unit.ID = int64(i.ID)
	return nil
}

func (r *RepoGorm) Edit(ctx context.Context, unit *entity.CookingUnit) error {
	i := &CookingUnit{}
	i.FromEntity(unit)
	err := r.db.WithContext(ctx).Save(i).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *RepoGorm) Delete(ctx context.Context, unit *entity.CookingUnit) error {
	i := &CookingUnit{}
	i.FromEntity(unit)
	err := r.db.WithContext(ctx).Delete(i).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrNotFound
		}
		return err
	}
	return nil
}
