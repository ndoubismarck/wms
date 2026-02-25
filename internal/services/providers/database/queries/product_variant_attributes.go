package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductAttributes struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductAttributeParams struct {
	ID        *string
	IDs       []string
	ProductID *string
}

func NewProductAttributes(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductAttributes {
	return &ProductAttributes{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductAttributes) query() *gorm.DB {
	return q.conn.Model(&entities.ProductAttribute{})
}

func (q *ProductAttributes) queryParams(params ProductAttributeParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.IDs != nil && len(params.IDs) > 0 {
		query = query.Where("id IN ?", params.IDs)
	}
	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	return query
}

func (q *ProductAttributes) Create(data entities.ProductAttribute) (entities.ProductAttribute, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductAttribute{}, err
	}
	return data, nil
}

func (q *ProductAttributes) Count(params ProductAttributeParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductAttributes) Update(params ProductAttributeParams, data entities.ProductAttribute) (entities.ProductAttribute, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductAttribute{}, false, nil
		}
		return entities.ProductAttribute{}, false, err
	}
	return data, true, nil
}

func (q *ProductAttributes) Exists(params ProductAttributeParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductAttributes) FindOne(params ProductAttributeParams) (entities.ProductAttribute, bool, error) {
	var result entities.ProductAttribute
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductAttributes) FindMany(params ProductAttributeParams, pagination types.IPagination) ([]entities.ProductAttribute, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductAttribute
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
			return []entities.ProductAttribute{}, &paginationResult, nil
		}
		return []entities.ProductAttribute{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductAttributes) Delete(params ProductAttributeParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductAttribute{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
