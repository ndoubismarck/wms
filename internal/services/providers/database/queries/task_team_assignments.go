package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type TaskTeamAssignments struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type TaskTeamAssignmentsParams struct {
	ID      *string
	TaskID  *string
	TaskIDs []string
	TeamID  *string
	TeamIDs []string
}

func NewTaskTeamAssignments(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *TaskTeamAssignments {
	return &TaskTeamAssignments{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *TaskTeamAssignments) query() *gorm.DB {
	return q.conn.Model(&entities.TaskTeamAssignment{})
}

func (q *TaskTeamAssignments) queryParams(params TaskTeamAssignmentsParams) *gorm.DB {
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
	if params.TeamID != nil {
		query = query.Where("team_id = ?", *params.TeamID)
	}
	if len(params.TeamIDs) > 0 {
		query = query.Where("team_id IN ?", params.TeamIDs)
	}
	return query
}

func (q *TaskTeamAssignments) CreateMany(data []entities.TaskTeamAssignment) error {
	if len(data) == 0 {
		return nil
	}
	if err := q.query().Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *TaskTeamAssignments) DeleteByTaskID(taskID string) error {
	result := q.query().Where("task_id = ?", taskID).Delete(&entities.TaskTeamAssignment{})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (q *TaskTeamAssignments) FindTeamIDs(params TaskTeamAssignmentsParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("team_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}

func (q *TaskTeamAssignments) FindTaskIDs(params TaskTeamAssignmentsParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("task_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}
