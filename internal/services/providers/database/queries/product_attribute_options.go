package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductAttributeOptions struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductAttributeOptionsParams struct {
	ID                 *string
	ProductAttributeID *string
}

func NewProductAttributeOptions(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductAttributeOptions {
	return &ProductAttributeOptions{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductAttributeOptions) query() *gorm.DB {
	return q.conn.Model(&entities.ProductAttributeOption{})
}

func (q *ProductAttributeOptions) queryParams(params ProductAttributeOptionsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.ProductAttributeID != nil {
		query = query.Where("product_attribute_id = ?", *params.ProductAttributeID)
	}
	return query
}

func (q *ProductAttributeOptions) Create(data entities.ProductAttributeOption) (entities.ProductAttributeOption, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductAttributeOption{}, err
	}
	return data, nil
}

func (q *ProductAttributeOptions) Count(params ProductAttributeOptionsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductAttributeOptions) Update(params ProductAttributeOptionsParams, data entities.ProductAttributeOption) (entities.ProductAttributeOption, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductAttributeOption{}, false, nil
		}
		return entities.ProductAttributeOption{}, false, err
	}
	return data, true, nil
}

func (q *ProductAttributeOptions) Exists(params ProductAttributeOptionsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductAttributeOptions) FindOne(params ProductAttributeOptionsParams) (entities.ProductAttributeOption, bool, error) {
	var result entities.ProductAttributeOption
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductAttributeOptions) FindMany(params ProductAttributeOptionsParams, pagination types.IPagination) ([]entities.ProductAttributeOption, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductAttributeOption
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
			return []entities.ProductAttributeOption{}, &paginationResult, nil
		}
		return []entities.ProductAttributeOption{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductAttributeOptions) Delete(params ProductAttributeOptionsParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductAttributeOption{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
