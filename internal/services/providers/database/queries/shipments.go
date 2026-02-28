package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Shipments struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ShipmentsParams struct {
	ID             *string
	LocationID     *string
	ShipmentNumber *string
	Status         []entities.ShipmentStatusEnum
	Related        bool
}

func NewShipments(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Shipments {
	return &Shipments{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Shipments) query() *gorm.DB {
	return q.conn.Model(&entities.Shipment{})
}

func (q *Shipments) queryParams(params ShipmentsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}
	if params.ShipmentNumber != nil {
		query = query.Where("shipment_number = ?", *params.ShipmentNumber)
	}
	if len(params.Status) > 0 {
		query = query.Where("status IN ?", params.Status)
	}
	if params.Related {
		query = query.Preload("Location")
	}
	return query
}

func (q *Shipments) Create(data entities.Shipment) (entities.Shipment, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Shipment{}, err
	}
	return data, nil
}

func (q *Shipments) Count(params ShipmentsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Shipments) Update(params ShipmentsParams, data entities.Shipment) (entities.Shipment, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Shipment{}, false, nil
		}
		return entities.Shipment{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Shipments) Exists(params ShipmentsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Shipments) FindOne(params ShipmentsParams) (entities.Shipment, bool, error) {
	var result entities.Shipment
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Shipments) FindMany(params ShipmentsParams, pagination types.IPagination) ([]entities.Shipment, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Shipment
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
			return []entities.Shipment{}, &paginationResult, nil
		}
		return []entities.Shipment{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Shipments) Delete(params ShipmentsParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Shipment{ID: *params.ID})
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
