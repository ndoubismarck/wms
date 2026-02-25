package queries

import (
	"errors"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ctx     types.IContext
	conn    *gorm.DB
	dialect types.Dialect
}

type InventoryParams struct {
	ID            *string
	BinID         *string
	LocationID    *string
	Related       bool
	ToUpdatedAt   *time.Time
	FromUpdatedAt *time.Time
}

func NewInventory(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Inventory {
	return &Inventory{
		ctx:     ctx,
		conn:    conn,
		dialect: dialect,
	}
}

func (q *Inventory) query(joins bool) *gorm.DB {
	query := q.conn.Model(&entities.Inventory{})
	query = query.Set("ctx", q.ctx)
	if joins {
		expr := "CONCAT(loc.code,' ',ais.code,' ',bay.code,' ',shf.code,' ',shl.code,' ',bin.code)"
		if q.dialect == types.SQLiteDialect {
			expr = "loc.code || ' ' || ais.code || ' ' || bay.code || ' ' || shf.code || ' ' || shl.code || ' ' || bin.code"
		}
		query = query.Select("inventory.*, " + expr + " AS location").
			Joins("JOIN location_bins bin ON bin.id = inventory.bin_id").
			Joins("JOIN location_shelf_levels shl ON shl.id = bin.shelf_level_id").
			Joins("JOIN locations_shelves shf ON shf.id = shl.shelf_id").
			Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
			Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
			Joins("JOIN locations loc ON loc.id = ais.location_id")
	}
	return query
}

func (q *Inventory) queryParams(joins bool, params InventoryParams) *gorm.DB {
	query := q.query(joins)
	if params.ID != nil {
		query = query.Where("inventory.id = ?", *params.ID)
	}
	if params.BinID != nil {
		query = query.Where("inventory.bin_id = ?", *params.BinID)
	}
	if params.LocationID != nil {
		if !joins {
			query = query.Select("inventory.*").
				Joins("JOIN location_bins bin ON bin.id = inventory.bin_id").
				Joins("JOIN location_shelf_levels shl ON shl.id = bin.shelf_level_id").
				Joins("JOIN locations_shelves shf ON shf.id = shl.shelf_id").
				Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
				Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
				Joins("JOIN locations loc ON loc.id = ais.location_id")
		}
		query = query.Where("loc.id = ?", *params.LocationID)
	}
	if params.FromUpdatedAt != nil && !params.FromUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("inventory.updated_at >= ?", *params.FromUpdatedAt)
	}
	if params.ToUpdatedAt != nil && !params.ToUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("inventory.updated_at <= ?", *params.ToUpdatedAt)
	}
	if params.Related {
		query = query.
			Preload("Variant.Media").
			Preload("Variant.Product.Category").
			Preload("Bin.ShelfLevel.Shelf.Bay.Aisle.Location")
	}
	return query
}

func (q *Inventory) Create(data entities.Inventory) (entities.Inventory, error) {
	if err := q.query(false).Create(&data).Error; err != nil {
		return entities.Inventory{}, err
	}
	return data, nil
}

func (q *Inventory) Count(params InventoryParams) (uint64, error) {
	var result int64
	query := q.queryParams(false, params)
	if err := query.Count(&result).Error; err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (q *Inventory) Update(params InventoryParams, data entities.Inventory) (entities.Inventory, bool, error) {
	query := q.queryParams(false, params)
	result := query.Updates(&data)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Inventory{}, false, nil
		}
		return entities.Inventory{}, false, err
	}
	return data, result.RowsAffected > 0, nil
}

func (q *Inventory) Exists(params InventoryParams) (bool, error) {
	var result int64
	query := q.queryParams(false, params)
	if err := query.Count(&result).Error; err != nil {
		return false, err
	}
	return result > 0, nil
}

func (q *Inventory) FindOne(params InventoryParams) (entities.Inventory, bool, error) {
	var result entities.Inventory
	query := q.queryParams(true, params)
	if err := query.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	return result, true, nil
}

func (q *Inventory) FindMany(params InventoryParams, pagination types.IPagination) ([]entities.Inventory, *types.PaginationResult, error) {
	var (
		count   int64
		results []entities.Inventory
	)
	query := q.queryParams(true, params)
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
			return []entities.Inventory{}, &paginationResult, nil
		}
		return []entities.Inventory{}, &paginationResult, err
	}
	return results, &paginationResult, nil
}

func (q *Inventory) GenerateSummary(params InventoryParams) ([]entities.InventorySummary, error) {
	var rows []entities.InventorySummary
	query := q.query(false).
		Joins("JOIN location_bins bin ON bin.id = inventory.bin_id").
		Joins("JOIN location_shelf_levels shl ON shl.id = bin.shelf_level_id").
		Joins("JOIN locations_shelves shf ON shf.id = shl.shelf_id").
		Joins("JOIN location_bays bay ON bay.id = shf.bay_id").
		Joins("JOIN location_aisles ais ON ais.id = bay.aisle_id").
		Joins("JOIN locations loc ON loc.id = ais.location_id")
	if params.ID != nil {
		query = query.Where("inventory.id = ?", *params.ID)
	}
	if params.BinID != nil {
		query = query.Where("inventory.bin_id = ?", *params.BinID)
	}
	if params.LocationID != nil {
		query = query.Where("loc.id = ?", *params.LocationID)
	}
	if params.FromUpdatedAt != nil && !params.FromUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("inventory.updated_at >= ?", *params.FromUpdatedAt)
	}
	if params.ToUpdatedAt != nil && !params.ToUpdatedAt.UTC().After(time.Now().UTC()) {
		query = query.Where("inventory.updated_at <= ?", *params.ToUpdatedAt)
	}
	if err := query.
		Select(`
			loc.id AS location_id,
			inventory.variant_id AS variant_id,
			SUM(inventory.quantity_on_hand) AS quantity_on_hand,
			SUM(inventory.quantity_reserved) AS quantity_reserved,
			SUM(inventory.quantity_incoming) AS quantity_incoming,
			SUM(inventory.quantity_outgoing) AS quantity_outgoing
		`).
		Group("loc.id, inventory.variant_id").
		Scan(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rows, nil
		}
		return nil, err
	}
	return rows, nil
}

func (q *Inventory) Delete(params InventoryParams) (bool, error) {
	if params.ID != nil {
		result := q.query(false).
			Delete(&entities.Inventory{ID: *params.ID})
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
