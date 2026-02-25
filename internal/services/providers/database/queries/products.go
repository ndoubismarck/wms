package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductParams struct {
	ID            *string
	LocationID    *string
	BinID         *string
	Related       bool
	Status        []entities.ProductStatusEnum
	ToUpdatedAt   *time.Time
	FromUpdatedAt *time.Time
}

func NewProducts(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Products {
	return &Products{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Products) query() *gorm.DB {
	return q.conn.Model(&entities.Product{})
}

func (q *Products) queryParams(params ProductParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_d = ?", *params.LocationID)
	}
	if params.BinID != nil {
		query = query.Where("bin_id = ?", *params.BinID)
	}
	if params.Status != nil && len(params.Status) > 0 {
		query = query.Where("status IN ?", params.Status)
	}
	if params.FromUpdatedAt != nil && !params.FromUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("updated_at >= ?", *params.FromUpdatedAt)
	}
	if params.ToUpdatedAt != nil && !params.ToUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("updated_at <= ?", *params.ToUpdatedAt)
	}
	if params.Related {
		query = query.
			Preload("Subcategory.Category.Brand").
			Preload("Attributes.ValueOption")
	}
	return query
}

func (q *Products) Create(data entities.Product) (entities.Product, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Product{}, err
	}
	return data, nil
}

func (q *Products) Count(params ProductParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Products) Update(params ProductParams, data entities.Product) (entities.Product, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Product{}, false, nil
		}
		return entities.Product{}, false, err
	}
	return data, true, nil
}

func (q *Products) Exists(params ProductParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Products) FindOne(params ProductParams) (entities.Product, bool, error) {
	var result entities.Product
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Products) FindMany(params ProductParams, pagination types.IPagination) ([]entities.Product, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Product
	)
	query := q.queryParams(params)
	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}
	paginationResult := pagination.GetResult(count)
	if err := query.
		Order(pagination.GetOrder()).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.Product{}, &paginationResult, nil
		}
		return []entities.Product{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Products) Delete(params ProductParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.Product{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
