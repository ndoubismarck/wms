package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type LocationShelfLevels struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationShelfLevelsParams struct {
	ID         *string
	ShelfID    *string
	LocationID *string
}

func NewLocationShelfLevels(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *LocationShelfLevels {
	return &LocationShelfLevels{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *LocationShelfLevels) query() *gorm.DB {
	return q.conn.Model(&entities.LocationShelfLevel{})
}

func (q *LocationShelfLevels) queryParams(params LocationShelfLevelsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("location_shelf_levels.id = ?", *params.ID)
	}
	if params.ShelfID != nil {
		query = query.Where("location_shelf_levels.shelf_id = ?", *params.ShelfID)
	}
	if params.LocationID != nil {
		query = query.
			Joins("JOIN locations_shelves shf ON shf.id = location_shelf_levels.shelf_id").
			Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
			Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
			Where("ais.location_id = ?", *params.LocationID)
	}
	return query
}

func (q *LocationShelfLevels) Create(data entities.LocationShelfLevel) (entities.LocationShelfLevel, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.LocationShelfLevel{}, err
	}
	return data, nil
}

func (q *LocationShelfLevels) FindCodes(params LocationShelfLevelsParams) ([]string, error) {
	var rows []struct {
		Code string `gorm:"column:code"`
	}
	query := q.queryParams(params).Select("location_shelf_levels.code")
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Code)
	}
	return result, nil
}

func (q *LocationShelfLevels) Count(params LocationShelfLevelsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *LocationShelfLevels) Update(params LocationShelfLevelsParams, data entities.LocationShelfLevel) (entities.LocationShelfLevel, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LocationShelfLevel{}, false, nil
		}
		return entities.LocationShelfLevel{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *LocationShelfLevels) Exists(params LocationShelfLevelsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *LocationShelfLevels) FindOne(params LocationShelfLevelsParams) (entities.LocationShelfLevel, bool, error) {
	var result entities.LocationShelfLevel
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *LocationShelfLevels) FindMany(params LocationShelfLevelsParams, pagination types.IPagination) ([]entities.LocationShelfLevel, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.LocationShelfLevel
	)
	query := q.queryParams(params)
	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}
	paginationResult := pagination.GetResult(count)
	if err := query.
		Order(pagination.GetOrderWithPrefix("location_shelf_levels")).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.LocationShelfLevel{}, &paginationResult, nil
		}
		return []entities.LocationShelfLevel{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *LocationShelfLevels) Delete(params LocationShelfLevelsParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.LocationShelfLevel{ID: *params.ID})
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
