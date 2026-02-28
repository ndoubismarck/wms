package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Tasks struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type TasksParams struct {
	ID                *string
	IDs               []string
	LocationID        *string
	Statuses          []entities.TaskStatusEnum
	Priorities        []entities.TaskPriorityEnum
	ScheduledByUserID *string
}

func NewTasks(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Tasks {
	return &Tasks{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Tasks) query() *gorm.DB {
	return q.conn.Model(&entities.Task{})
}

func (q *Tasks) queryParams(params TasksParams) *gorm.DB {
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
	if len(params.Statuses) > 0 {
		query = query.Where("status IN ?", params.Statuses)
	}
	if len(params.Priorities) > 0 {
		query = query.Where("priority IN ?", params.Priorities)
	}
	if params.ScheduledByUserID != nil {
		query = query.Where("scheduled_by_user_id = ?", *params.ScheduledByUserID)
	}
	return query
}

func (q *Tasks) Create(data entities.Task) (entities.Task, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.Task{}, err
	}
	return data, nil
}

func (q *Tasks) Count(params TasksParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Tasks) Update(params TasksParams, data entities.Task) (entities.Task, bool, error) {
	result := q.queryParams(params).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Task{}, false, nil
		}
		return entities.Task{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Tasks) Exists(params TasksParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Tasks) FindOne(params TasksParams) (entities.Task, bool, error) {
	var result entities.Task
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Tasks) FindMany(params TasksParams, pagination types.IPagination) ([]entities.Task, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Task
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
			return []entities.Task{}, &paginationResult, nil
		}
		return []entities.Task{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Tasks) Delete(params TasksParams) (bool, error) {
	if params.ID != nil {
		result := q.query().Delete(&entities.Task{ID: *params.ID})
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
