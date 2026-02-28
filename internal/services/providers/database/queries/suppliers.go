package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Suppliers struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type SuppliersParams struct {
	ID         *string
	LocationID *string
	Name       *string
	Status     []entities.SupplierStatusEnum
	Types      []entities.SupplierTypeEnum
	Related    bool
}

func NewSuppliers(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Suppliers {
	return &Suppliers{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Suppliers) query() *gorm.DB {
	return q.conn.Model(&entities.Supplier{})
}

func (q *Suppliers) queryParams(params SuppliersParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}
	if params.Name != nil {
		query = query.Where("name = ?", *params.Name)
	}
	if len(params.Status) > 0 {
		query = query.Where("status IN ?", params.Status)
	}
	if len(params.Types) > 0 {
		query = query.Where("supplier_type IN ?", params.Types)
	}
	if params.Related {
		query = query.Preload("Location")
	}
	return query
}

func (q *Suppliers) Create(data entities.Supplier) (entities.Supplier, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Supplier{}, err
	}
	return data, nil
}

func (q *Suppliers) Count(params SuppliersParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Suppliers) Update(params SuppliersParams, data entities.Supplier) (entities.Supplier, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Supplier{}, false, nil
		}
		return entities.Supplier{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Suppliers) Exists(params SuppliersParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Suppliers) FindOne(params SuppliersParams) (entities.Supplier, bool, error) {
	var result entities.Supplier
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Suppliers) FindMany(params SuppliersParams, pagination types.IPagination) ([]entities.Supplier, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Supplier
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
			return []entities.Supplier{}, &paginationResult, nil
		}
		return []entities.Supplier{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Suppliers) Delete(params SuppliersParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Supplier{ID: *params.ID})
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
