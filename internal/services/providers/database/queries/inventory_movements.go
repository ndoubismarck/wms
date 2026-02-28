package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type InventoryMovements struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type InventoryMovementsParams struct {
	ID           *string
	InventoryID  *string
	ReferenceID  *string
	LocationID   *string
	MovementType []entities.InventoryMovementTypeEnum
}

func NewInventoryMovements(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *InventoryMovements {
	return &InventoryMovements{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *InventoryMovements) query() *gorm.DB {
	return q.conn.Model(&entities.InventoryMovement{})
}

func (q *InventoryMovements) queryParams(params InventoryMovementsParams) *gorm.DB {
	query := q.query()
	if params.ID != nil {
		query = query.Where("inventory_movements.id = ?", *params.ID)
	}
	if params.InventoryID != nil {
		query = query.Where("inventory_movements.inventory_id = ?", *params.InventoryID)
	}
	if params.ReferenceID != nil {
		query = query.Where("inventory_movements.reference_id = ?", *params.ReferenceID)
	}
	if params.LocationID != nil {
		query = query.
			Joins("JOIN inventory inv ON inv.id = inventory_movements.inventory_id").
			Joins("JOIN location_bins bin ON bin.id = inv.bin_id").
			Joins("JOIN location_shelf_levels shl ON shl.id = bin.shelf_level_id").
			Joins("JOIN locations_shelves shf ON shf.id = shl.shelf_id").
			Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
			Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
			Where("ais.location_id = ?", *params.LocationID)
	}
	if len(params.MovementType) > 0 {
		query = query.Where("inventory_movements.movement_type IN ?", params.MovementType)
	}
	return query
}

func (q *InventoryMovements) Create(data entities.InventoryMovement) (entities.InventoryMovement, error) {
	if err := q.query().Create(&data).Error; err != nil {
		return entities.InventoryMovement{}, err
	}
	return data, nil
}

func (q *InventoryMovements) Count(params InventoryMovementsParams) (uint64, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *InventoryMovements) Update(params InventoryMovementsParams, data entities.InventoryMovement) (entities.InventoryMovement, bool, error) {
	query := q.queryParams(params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.InventoryMovement{}, false, nil
		}
		return entities.InventoryMovement{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *InventoryMovements) Exists(params InventoryMovementsParams) (bool, error) {
	var result int64
	query := q.queryParams(params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *InventoryMovements) FindOne(params InventoryMovementsParams) (entities.InventoryMovement, bool, error) {
	var result entities.InventoryMovement
	query := q.queryParams(params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *InventoryMovements) FindMany(params InventoryMovementsParams, pagination types.IPagination) ([]entities.InventoryMovement, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.InventoryMovement
	)
	query := q.queryParams(params)
	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}
	paginationResult := pagination.GetResult(count)
	if err := query.
		Order(pagination.GetOrderWithPrefix("inventory_movements")).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.InventoryMovement{}, &paginationResult, nil
		}
		return []entities.InventoryMovement{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *InventoryMovements) Delete(params InventoryMovementsParams) (bool, error) {
	if params.ID != nil {
		result := q.query().
			Delete(&entities.InventoryMovement{ID: *params.ID})
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
