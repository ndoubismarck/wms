package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type TeamMembers struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type TeamMembersParams struct {
	ID      *string
	TeamID  *string
	TeamIDs []string
	UserID  *string
	UserIDs []string
}

func NewTeamMembers(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *TeamMembers {
	return &TeamMembers{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *TeamMembers) query() *gorm.DB {
	return q.conn.Model(&entities.TeamMember{})
}

func (q *TeamMembers) queryParams(params TeamMembersParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.TeamID != nil {
		query = query.Where("team_id = ?", *params.TeamID)
	}
	if len(params.TeamIDs) > 0 {
		query = query.Where("team_id IN ?", params.TeamIDs)
	}
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if len(params.UserIDs) > 0 {
		query = query.Where("user_id IN ?", params.UserIDs)
	}
	return query
}

func (q *TeamMembers) CreateMany(data []entities.TeamMember) error {
	if len(data) == 0 {
		return nil
	}
	if err := q.query().Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *TeamMembers) DeleteByTeamID(teamID string) error {
	result := q.query().Where("team_id = ?", teamID).Delete(&entities.TeamMember{})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (q *TeamMembers) FindUserIDs(params TeamMembersParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("user_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}

func (q *TeamMembers) FindTeamIDs(params TeamMembersParams) ([]string, error) {
	results := []string{}
	if err := q.queryParams(params).Pluck("team_id", &results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	return results, nil
}
