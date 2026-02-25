package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductBrands struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductBrandsParams struct {
	ID         *string
	LocationID *string
}

func NewProductBrands(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductBrands {
	return &ProductBrands{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductBrands) query() *gorm.DB {
	return q.conn.Model(&entities.ProductBrand{})
}

func (q *ProductBrands) queryParams(params ProductBrandsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_d = ?", *params.LocationID)
	}
	return query
}

func (q *ProductBrands) Create(data entities.ProductBrand) (entities.ProductBrand, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductBrand{}, err
	}
	return data, nil
}

func (q *ProductBrands) Count(params ProductBrandsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductBrands) Update(params ProductBrandsParams, data entities.ProductBrand) (entities.ProductBrand, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductBrand{}, false, nil
		}
		return entities.ProductBrand{}, false, err
	}
	return data, true, nil
}

func (q *ProductBrands) Exists(params ProductBrandsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductBrands) FindOne(params ProductBrandsParams) (entities.ProductBrand, bool, error) {
	var result entities.ProductBrand
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductBrands) FindMany(params ProductBrandsParams, pagination types.IPagination) ([]entities.ProductBrand, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductBrand
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
			return []entities.ProductBrand{}, &paginationResult, nil
		}
		return []entities.ProductBrand{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductBrands) Delete(params ProductBrandsParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductBrand{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
