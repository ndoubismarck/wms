package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Customers struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type CustomersParams struct {
	ID         *string
	LocationID *string
	Email      *string
	Tier       []entities.CustomerTierEnum
	Status     []entities.CustomerStatusEnum
	Related    bool
}

func NewCustomers(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Customers {
	return &Customers{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Customers) query() *gorm.DB {
	return q.conn.Model(&entities.Customer{})
}

func (q *Customers) queryParams(params CustomersParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}
	if params.Email != nil {
		query = query.Where("email = ?", *params.Email)
	}
	if len(params.Tier) > 0 {
		query = query.Where("tier IN ?", params.Tier)
	}
	if len(params.Status) > 0 {
		query = query.Where("status IN ?", params.Status)
	}
	if params.Related {
		query = query.Preload("Location")
	}
	return query
}

func (q *Customers) Create(data entities.Customer) (entities.Customer, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Customer{}, err
	}
	return data, nil
}

func (q *Customers) Count(params CustomersParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Customers) Update(params CustomersParams, data entities.Customer) (entities.Customer, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Customer{}, false, nil
		}
		return entities.Customer{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Customers) Exists(params CustomersParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Customers) FindOne(params CustomersParams) (entities.Customer, bool, error) {
	var result entities.Customer
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Customers) FindMany(params CustomersParams, pagination types.IPagination) ([]entities.Customer, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Customer
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
			return []entities.Customer{}, &paginationResult, nil
		}
		return []entities.Customer{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Customers) Delete(params CustomersParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Customer{ID: *params.ID})
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
