package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Orders struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type OrdersParams struct {
	ID          *string
	LocationID  *string
	OrderNumber *string
	Status      []entities.OrderStatusEnum
	Priority    []entities.OrderPriorityEnum
	Related     bool
}

func NewOrders(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Orders {
	return &Orders{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Orders) query() *gorm.DB {
	return q.conn.Model(&entities.Order{})
}

func (q *Orders) queryParams(params OrdersParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}
	if params.OrderNumber != nil {
		query = query.Where("order_number = ?", *params.OrderNumber)
	}
	if len(params.Status) > 0 {
		query = query.Where("status IN ?", params.Status)
	}
	if len(params.Priority) > 0 {
		query = query.Where("priority IN ?", params.Priority)
	}
	if params.Related {
		query = query.Preload("Location")
	}
	return query
}

func (q *Orders) Create(data entities.Order) (entities.Order, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Order{}, err
	}
	return data, nil
}

func (q *Orders) Count(params OrdersParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Orders) Update(params OrdersParams, data entities.Order) (entities.Order, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Order{}, false, nil
		}
		return entities.Order{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Orders) Exists(params OrdersParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Orders) FindOne(params OrdersParams) (entities.Order, bool, error) {
	var result entities.Order
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Orders) FindMany(params OrdersParams, pagination types.IPagination) ([]entities.Order, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Order
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
			return []entities.Order{}, &paginationResult, nil
		}
		return []entities.Order{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Orders) Delete(params OrdersParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Order{ID: *params.ID})
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
