package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Users struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type UserFilters struct {
	Roles []string
}

func NewUsers(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Users {
	return &Users{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Users) query() *gorm.DB {
	return q.conn.Model(&entities.User{})
}

func (q *Users) queryWithFilters(filters UserFilters) *gorm.DB {
	query := q.query()
	if len(filters.Roles) > 0 {
		query = query.Where("role IN ?", filters.Roles)
	}
	return query
}

func (q *Users) Create(data entities.User) (entities.User, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.User{}, err
	}
	return data, nil
}

func (q *Users) Count(filters UserFilters) (uint64, error) {
	var result int64
	query := q.queryWithFilters(filters)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Users) Upsert(data entities.User) (entities.User, error) {
	if err := q.query().Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "email_address",
			},
		},
	}).Create(&data).Error; err != nil {
		return entities.User{}, err
	}
	return data, nil
}

func (q *Users) UpdateByID(id string, data entities.User) (entities.User, bool, error) {
	result := q.query().Where("id = ?", id).Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, false, nil
		}
		return entities.User{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Users) ExistsByID(id string) (bool, error) {
	var result int64
	if err := q.query().Where("id = ?", id).Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Users) ExistsByEmailAddress(emailAddress string) (bool, error) {
	var result int64
	if err := q.query().Where("email_address = ?", emailAddress).Count(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return result > 0, nil
}

func (q *Users) FindAll(filters UserFilters, pagination types.IPagination) ([]entities.User, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.User
	)
	query := q.queryWithFilters(filters)
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
			return []entities.User{}, &paginationResult, nil
		}
		return []entities.User{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Users) FindOneByID(id string) (entities.User, bool, error) {
	var result entities.User
	if err := q.query().
		Where("id = ?", id).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Users) FindOneByEmailAddress(email string) (entities.User, bool, error) {
	var result entities.User
	if err := q.query().
		Where("email_address = ?", email).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Users) DeleteByID(id string) (bool, error) {
	result := q.query().
		Delete(&entities.User{ID: id})
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return result.RowsAffected > 0, nil
}
