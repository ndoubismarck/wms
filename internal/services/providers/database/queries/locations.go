package queries

import (
	"errors"
	"fmt"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Locations struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type LocationsParams struct {
	ID *string
}

func NewLocations(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Locations {
	return &Locations{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Locations) query() *gorm.DB {
	return q.conn.Model(&entities.Location{})
}

func (q *Locations) queryParams(params LocationsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	return query
}

func (q *Locations) Create(data entities.Location) (entities.Location, error) {
	if err := q.query().Transaction(func(tx *gorm.DB) error {
		var lastErr error
		for i := 0; i < 10; i++ {
			lastErr = tx.Create(&data).Error
			if lastErr == nil {
				return nil
			}
			if errors.Is(lastErr, gorm.ErrDuplicatedKey) {
				data.ID = ""
				data.Code = 0
				continue
			}
			return lastErr
		}
		return fmt.Errorf("could not create record after retries: %w", lastErr)
	}); err != nil {
		return entities.Location{}, err
	}
	return data, nil
}

func (q *Locations) Count(params LocationsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Locations) Update(params LocationsParams, data entities.Location) (entities.Location, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Location{}, false, nil
		}
		return entities.Location{}, false, err
	}
	return data, true, nil
}

func (q *Locations) Exists(params LocationsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Locations) FindOne(params LocationsParams) (entities.Location, bool, error) {
	var result entities.Location
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Locations) FindMany(params LocationsParams, pagination types.IPagination) ([]entities.Location, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Location
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
			return []entities.Location{}, &paginationResult, nil
		}
		return []entities.Location{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Locations) Delete(params LocationsParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.Location{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
