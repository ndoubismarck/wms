package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductVariantMedia struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductVariantMediaParams struct {
	ID               *string
	ProductID        *string
	ProductVariantID *string
	FileType         []entities.ProductFileTypeEnum
}

func NewProductVariantMedia(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductVariantMedia {
	return &ProductVariantMedia{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductVariantMedia) query() *gorm.DB {
	return q.conn.Model(&entities.ProductVariantMedia{})
}

func (q *ProductVariantMedia) queryParams(params ProductVariantMediaParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	if params.ProductVariantID != nil {
		query = query.Where("product_variant_id = ?", *params.ProductVariantID)
	}
	if len(params.FileType) > 0 {
		query = query.Where("file_type IN ?", params.FileType)
	}
	return query
}

func (q *ProductVariantMedia) Create(data entities.ProductVariantMedia) (entities.ProductVariantMedia, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductVariantMedia{}, err
	}
	return data, nil
}

func (q *ProductVariantMedia) Count(params ProductVariantMediaParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductVariantMedia) Update(params ProductVariantMediaParams, data entities.ProductVariantMedia) (entities.ProductVariantMedia, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductVariantMedia{}, false, nil
		}
		return entities.ProductVariantMedia{}, false, err
	}
	return data, true, nil
}

func (q *ProductVariantMedia) Exists(params ProductVariantMediaParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductVariantMedia) FindOne(params ProductVariantMediaParams) (entities.ProductVariantMedia, bool, error) {
	var result entities.ProductVariantMedia
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductVariantMedia) FindMany(params ProductVariantMediaParams, pagination types.IPagination) ([]entities.ProductVariantMedia, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductVariantMedia
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
			return []entities.ProductVariantMedia{}, &paginationResult, nil
		}
		return []entities.ProductVariantMedia{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductVariantMedia) Delete(params ProductVariantMediaParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductVariantMedia{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
