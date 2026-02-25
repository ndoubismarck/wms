package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type LocationBays struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationBaysParams struct {
	ID         *string
	AisleID    *string
	LocationID *string
}

func NewLocationBays(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *LocationBays {
	return &LocationBays{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *LocationBays) query() *gorm.DB {
	return q.conn.Model(&entities.LocationBay{})
}

func (q *LocationBays) queryParams(params LocationBaysParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("location_bays.id = ?", *params.ID)
	}
	if params.AisleID != nil {
		query = query.Where("location_bays.aisle_id = ?", *params.AisleID)
	}
	if params.LocationID != nil {
		query = query.
			Joins("JOIN location_aisles ais ON ais.id = location_bays.aisle_id").
			Where("ais.location_id = ?", *params.LocationID)
	}
	return query
}

func (q *LocationBays) Create(data entities.LocationBay) (entities.LocationBay, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.LocationBay{}, err
	}
	return data, nil
}

func (q *LocationBays) Count(params LocationBaysParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *LocationBays) Update(params LocationBaysParams, data entities.LocationBay) (entities.LocationBay, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LocationBay{}, false, nil
		}
		return entities.LocationBay{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *LocationBays) Exists(params LocationBaysParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *LocationBays) FindOne(params LocationBaysParams) (entities.LocationBay, bool, error) {
	var result entities.LocationBay
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *LocationBays) FindMany(params LocationBaysParams, pagination types.IPagination) ([]entities.LocationBay, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.LocationBay
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
			return []entities.LocationBay{}, &paginationResult, nil
		}
		return []entities.LocationBay{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *LocationBays) Delete(params LocationBaysParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.LocationBay{ID: *params.ID})
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
