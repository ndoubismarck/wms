package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type LocationAisles struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationAislesParams struct {
	ID         *string
	LocationID *string
}

func NewLocationAisles(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *LocationAisles {
	return &LocationAisles{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *LocationAisles) query() *gorm.DB {
	return q.conn.Model(&entities.LocationAisle{})
}

func (q *LocationAisles) queryParams(params LocationAislesParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("location_aisles.id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_aisles.location_id = ?", *params.LocationID)
	}
	return query
}

func (q *LocationAisles) Create(data entities.LocationAisle) (entities.LocationAisle, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.LocationAisle{}, err
	}
	return data, nil
}

func (q *LocationAisles) FindCodes(params LocationAislesParams) ([]string, error) {
	var rows []struct {
		Code string `gorm:"column:code"`
	}
	query := q.queryParams(params).Select("location_aisles.code")
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Code)
	}
	return result, nil
}

func (q *LocationAisles) Count(params LocationAislesParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *LocationAisles) Update(params LocationAislesParams, data entities.LocationAisle) (entities.LocationAisle, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LocationAisle{}, false, nil
		}
		return entities.LocationAisle{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *LocationAisles) Exists(params LocationAislesParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *LocationAisles) FindOne(params LocationAislesParams) (entities.LocationAisle, bool, error) {
	var result entities.LocationAisle
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *LocationAisles) FindMany(params LocationAislesParams, pagination types.IPagination) ([]entities.LocationAisle, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.LocationAisle
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
			return []entities.LocationAisle{}, &paginationResult, nil
		}
		return []entities.LocationAisle{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *LocationAisles) Delete(params LocationAislesParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.LocationAisle{ID: *params.ID})
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
