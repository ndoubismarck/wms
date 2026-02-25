package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductVariants struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductVariantsParams struct {
	ID        *string
	ProductID *string
	SKU       *string
	Related   bool
}

func NewProductVariants(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductVariants {
	return &ProductVariants{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductVariants) query() *gorm.DB {
	return q.conn.Model(&entities.ProductVariant{})
}

func (q *ProductVariants) queryParams(params ProductVariantsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	if params.SKU != nil {
		query = query.Where("sku = ?", *params.SKU)
	}
	if params.Related {
		query = query.
			Preload("Product").
			Preload("Attributes.ProductAttributeDefinition").
			Preload("Attributes.ValueOption").
			Preload("Media")
	}
	return query
}

func (q *ProductVariants) Create(data entities.ProductVariant) (entities.ProductVariant, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductVariant{}, err
	}
	return data, nil
}

func (q *ProductVariants) Count(params ProductVariantsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductVariants) Update(params ProductVariantsParams, data entities.ProductVariant) (entities.ProductVariant, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductVariant{}, false, nil
		}
		return entities.ProductVariant{}, false, err
	}
	return data, true, nil
}

func (q *ProductVariants) Exists(params ProductVariantsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductVariants) FindOne(params ProductVariantsParams) (entities.ProductVariant, bool, error) {
	var result entities.ProductVariant
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductVariants) FindMany(params ProductVariantsParams, pagination types.IPagination) ([]entities.ProductVariant, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductVariant
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
			return []entities.ProductVariant{}, &paginationResult, nil
		}
		return []entities.ProductVariant{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductVariants) Delete(params ProductVariantsParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductVariant{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
