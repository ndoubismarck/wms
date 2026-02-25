package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductCategories struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductCategoriesParams struct {
	ID      *string
	BrandID *string
}

func NewProductCategories(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductCategories {
	return &ProductCategories{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductCategories) query() *gorm.DB {
	return q.conn.Model(&entities.ProductCategory{})
}

func (q *ProductCategories) queryParams(params ProductCategoriesParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.BrandID != nil {
		query = query.Where("brand_id = ?", *params.BrandID)
	}
	return query
}

func (q *ProductCategories) Create(data entities.ProductCategory) (entities.ProductCategory, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductCategory{}, err
	}
	return data, nil
}

func (q *ProductCategories) Count(params ProductCategoriesParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductCategories) Update(params ProductCategoriesParams, data entities.ProductCategory) (entities.ProductCategory, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductCategory{}, false, nil
		}
		return entities.ProductCategory{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *ProductCategories) Exists(params ProductCategoriesParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductCategories) FindOne(params ProductCategoriesParams) (entities.ProductCategory, bool, error) {
	var result entities.ProductCategory
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductCategories) FindMany(params ProductCategoriesParams, pagination types.IPagination) ([]entities.ProductCategory, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductCategory
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
			return []entities.ProductCategory{}, &paginationResult, nil
		}
		return []entities.ProductCategory{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductCategories) Delete(params ProductCategoriesParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.ProductCategory{ID: *params.ID})
		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return result.RowsAffected > 0, nil
	}
	return false, nil
}
