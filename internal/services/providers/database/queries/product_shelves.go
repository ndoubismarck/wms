package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type LocationShelves struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationShelvesParams struct {
	ID         *string
	BayID      *string
	LocationID *string
}

func NewLocationShelves(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *LocationShelves {
	return &LocationShelves{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *LocationShelves) query() *gorm.DB {
	return q.conn.Model(&entities.LocationShelf{})
}

func (q *LocationShelves) queryParams(params LocationShelvesParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("locations_shelves.id = ?", *params.ID)
	}
	if params.BayID != nil {
		query = query.Where("locations_shelves.bay_id = ?", *params.BayID)
	}
	if params.LocationID != nil {
		query = query.
			Joins("JOIN location_bays bay ON bay.id = locations_shelves.bay_id").
			Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
			Where("ais.location_id = ?", *params.LocationID)
	}
	return query
}

func (q *LocationShelves) Create(data entities.LocationShelf) (entities.LocationShelf, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.LocationShelf{}, err
	}
	return data, nil
}

func (q *LocationShelves) Count(params LocationShelvesParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *LocationShelves) Update(params LocationShelvesParams, data entities.LocationShelf) (entities.LocationShelf, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LocationShelf{}, false, nil
		}
		return entities.LocationShelf{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *LocationShelves) Exists(params LocationShelvesParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *LocationShelves) FindOne(params LocationShelvesParams) (entities.LocationShelf, bool, error) {
	var result entities.LocationShelf
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *LocationShelves) FindMany(params LocationShelvesParams, pagination types.IPagination) ([]entities.LocationShelf, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.LocationShelf
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
			return []entities.LocationShelf{}, &paginationResult, nil
		}
		return []entities.LocationShelf{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *LocationShelves) Delete(params LocationShelvesParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.LocationShelf{ID: *params.ID})
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
