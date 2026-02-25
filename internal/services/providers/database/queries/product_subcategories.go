package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductSubcategories struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductSubcategoriesParams struct {
	ID         *string
	CategoryID *string
}

func NewProductSubcategories(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductSubcategories {
	return &ProductSubcategories{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductSubcategories) query() *gorm.DB {
	return q.conn.Model(&entities.ProductSubcategory{})
}

func (q *ProductSubcategories) queryParams(params ProductSubcategoriesParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	return query
}

func (q *ProductSubcategories) Create(data entities.ProductSubcategory) (entities.ProductSubcategory, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductSubcategory{}, err
	}
	return data, nil
}

func (q *ProductSubcategories) Count(params ProductSubcategoriesParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductSubcategories) Update(params ProductSubcategoriesParams, data entities.ProductSubcategory) (entities.ProductSubcategory, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductSubcategory{}, false, nil
		}
		return entities.ProductSubcategory{}, false, err
	}
	return data, true, nil
}

func (q *ProductSubcategories) Exists(params ProductSubcategoriesParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductSubcategories) FindOne(params ProductSubcategoriesParams) (entities.ProductSubcategory, bool, error) {
	var result entities.ProductSubcategory
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductSubcategories) FindMany(params ProductSubcategoriesParams, pagination types.IPagination) ([]entities.ProductSubcategory, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductSubcategory
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
			return []entities.ProductSubcategory{}, &paginationResult, nil
		}
		return []entities.ProductSubcategory{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductSubcategories) Delete(params ProductSubcategoriesParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductSubcategory{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
