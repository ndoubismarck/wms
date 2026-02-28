package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Teams struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type TeamsParams struct {
	ID         *string
	IDs        []string
	LocationID *string
	Name       *string
}

func NewTeams(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Teams {
	return &Teams{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Teams) query() *gorm.DB {
	return q.conn.Model(&entities.Team{})
}

func (q *Teams) queryParams(params TeamsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if len(params.IDs) > 0 {
		query = query.Where("id IN ?", params.IDs)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}
	if params.Name != nil {
		query = query.Where("name = ?", *params.Name)
	}
	return query
}

func (q *Teams) Create(data entities.Team) (entities.Team, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Team{}, err
	}
	return data, nil
}

func (q *Teams) Count(params TeamsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Teams) Update(params TeamsParams, data entities.Team) (entities.Team, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Team{}, false, nil
		}
		return entities.Team{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Teams) Exists(params TeamsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Teams) FindOne(params TeamsParams) (entities.Team, bool, error) {
	var result entities.Team
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Teams) FindAll(params TeamsParams) ([]entities.Team, error) {
	results := []entities.Team{}
	if err := q.queryParams(params).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.Team{}, nil
		}
		return nil, err
	}
	return results, nil
}

func (q *Teams) FindMany(params TeamsParams, pagination types.IPagination) ([]entities.Team, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Team
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
			return []entities.Team{}, &paginationResult, nil
		}
		return []entities.Team{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Teams) FindIDs(params TeamsParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}

func (q *Teams) Delete(params TeamsParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Team{ID: *params.ID})
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
