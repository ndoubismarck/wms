package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type ProductAttributeDefinitions struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type ProductAttributeDefinitionsParams struct {
	ID         *string
	LocationID *string
	Key        *string
	IsSystem   *bool
	AppliesTo  []entities.ProductAttributeAppliesToEnum
	DataType   []entities.ProductAttributeDataTypeEnum
}

func NewProductAttributeDefinitions(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *ProductAttributeDefinitions {
	return &ProductAttributeDefinitions{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *ProductAttributeDefinitions) query() *gorm.DB {
	return q.conn.Model(&entities.ProductAttributeDefinition{})
}

func (q *ProductAttributeDefinitions) queryParams(params ProductAttributeDefinitionsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.LocationID != nil {
		query = query.Where("location_d = ?", *params.LocationID)
	}
	if params.Key != nil {
		query = query.Where("key = ?", *params.Key)
	}
	if params.IsSystem != nil {
		query = query.Where("is_system = ?", *params.IsSystem)
	}
	if len(params.AppliesTo) > 0 {
		query = query.Where("applies_to IN ?", params.AppliesTo)
	}
	if len(params.DataType) > 0 {
		query = query.Where("data_type IN ?", params.DataType)
	}
	return query
}

func (q *ProductAttributeDefinitions) Create(data entities.ProductAttributeDefinition) (entities.ProductAttributeDefinition, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.ProductAttributeDefinition{}, err
	}
	return data, nil
}

func (q *ProductAttributeDefinitions) Count(params ProductAttributeDefinitionsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *ProductAttributeDefinitions) Update(params ProductAttributeDefinitionsParams, data entities.ProductAttributeDefinition) (entities.ProductAttributeDefinition, bool, error) {
	query := q.queryParams(params)
	if err := query.Updates(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ProductAttributeDefinition{}, false, nil
		}
		return entities.ProductAttributeDefinition{}, false, err
	}
	return data, true, nil
}

func (q *ProductAttributeDefinitions) Exists(params ProductAttributeDefinitionsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *ProductAttributeDefinitions) FindOne(params ProductAttributeDefinitionsParams) (entities.ProductAttributeDefinition, bool, error) {
	var result entities.ProductAttributeDefinition
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *ProductAttributeDefinitions) FindMany(params ProductAttributeDefinitionsParams, pagination types.IPagination) ([]entities.ProductAttributeDefinition, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.ProductAttributeDefinition
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
			return []entities.ProductAttributeDefinition{}, &paginationResult, nil
		}
		return []entities.ProductAttributeDefinition{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *ProductAttributeDefinitions) Delete(params ProductAttributeDefinitionsParams) (bool, error) {
	if params.ID != nil {
		if err := q.query().
			Delete(&entities.ProductAttributeDefinition{ID: *params.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}
