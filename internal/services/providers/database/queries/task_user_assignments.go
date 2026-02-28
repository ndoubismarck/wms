package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type TaskUserAssignments struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type TaskUserAssignmentsParams struct {
	ID      *string
	TaskID  *string
	TaskIDs []string
	UserID  *string
	UserIDs []string
}

func NewTaskUserAssignments(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *TaskUserAssignments {
	return &TaskUserAssignments{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *TaskUserAssignments) query() *gorm.DB {
	return q.conn.Model(&entities.TaskUserAssignment{})
}

func (q *TaskUserAssignments) queryParams(params TaskUserAssignmentsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.TaskID != nil {
		query = query.Where("task_id = ?", *params.TaskID)
	}
	if len(params.TaskIDs) > 0 {
		query = query.Where("task_id IN ?", params.TaskIDs)
	}
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if len(params.UserIDs) > 0 {
		query = query.Where("user_id IN ?", params.UserIDs)
	}
	return query
}

func (q *TaskUserAssignments) CreateMany(data []entities.TaskUserAssignment) error {
	if len(data) == 0 {
		return nil
	}
	if err := q.query().Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *TaskUserAssignments) DeleteByTaskID(taskID string) error {
	result := q.query().Where("task_id = ?", taskID).Delete(&entities.TaskUserAssignment{})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (q *TaskUserAssignments) FindUserIDs(params TaskUserAssignmentsParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("user_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}

func (q *TaskUserAssignments) FindTaskIDs(params TaskUserAssignmentsParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("task_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}
