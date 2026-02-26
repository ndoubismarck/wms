package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type LocationBins struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationBinsParams struct {
	ID           *string
	ShelfLevelID *string
	LocationID   *string
}

func NewLocationBins(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *LocationBins {
	return &LocationBins{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *LocationBins) query() *gorm.DB {
	return q.conn.Model(&entities.LocationBin{})
}

func (q *LocationBins) queryParams(params LocationBinsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("location_bins.id = ?", *params.ID)
	}
	if params.ShelfLevelID != nil {
		query = query.Where("location_bins.shelf_level_id = ?", *params.ShelfLevelID)
	}
	if params.LocationID != nil {
		query = query.
			Joins("JOIN location_shelf_levels shl ON shl.id = location_bins.shelf_level_id").
			Joins("JOIN locations_shelves shf ON shf.id = shl.shelf_id").
			Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
			Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
			Where("ais.location_id = ?", *params.LocationID)
	}
	return query
}

func (q *LocationBins) Create(data entities.LocationBin) (entities.LocationBin, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.LocationBin{}, err
	}
	return data, nil
}

func (q *LocationBins) FindCodes(params LocationBinsParams) ([]string, error) {
	var rows []struct {
		Code string `gorm:"column:code"`
	}
	query := q.queryParams(params).Select("location_bins.code")
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Code)
	}
	return result, nil
}

func (q *LocationBins) Count(params LocationBinsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *LocationBins) Update(params LocationBinsParams, data entities.LocationBin) (entities.LocationBin, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LocationBin{}, false, nil
		}
		return entities.LocationBin{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *LocationBins) Exists(params LocationBinsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *LocationBins) FindOne(params LocationBinsParams) (entities.LocationBin, bool, error) {
	var result entities.LocationBin
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *LocationBins) FindMany(params LocationBinsParams, pagination types.IPagination) ([]entities.LocationBin, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.LocationBin
	)
	query := q.queryParams(params)
	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}
	paginationResult := pagination.GetResult(count)
	if err := query.
		Order(pagination.GetOrderWithPrefix("location_bins")).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.LocationBin{}, &paginationResult, nil
		}
		return []entities.LocationBin{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *LocationBins) Delete(params LocationBinsParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.LocationBin{ID: *params.ID})
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
