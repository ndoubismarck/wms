package database

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"gorm.io/gorm"
)

type Provider struct {
	ctx     types.IContext
	conn    *gorm.DB
	queries *Query
}

func New(ctx types.IContext) *Provider {
	p := &Provider{ctx: ctx}
	ctx.Hooks().OnStart(p.start)
	return p
}

func (p *Provider) start() error {
	conn, ok, err := p.ctx.Database().MySQL().Conn()
	if err != nil {
		return err
	}
	if ok {
		p.conn = conn
		p.queries = newQuery(p.ctx, p.conn, types.MySQLDialect)
		if err = p.migrate(); err != nil {
			return err
		}
		p.ctx.Events().Emit(types.EventTypeDatabaseMigrationComplete)
		return nil
	}

	conn, ok, err = p.ctx.Database().SQLite().Conn()
	if err != nil {
		return err
	}
	if ok {
		p.conn = conn
		p.queries = newQuery(p.ctx, p.conn, types.SQLiteDialect)
		if err = p.migrate(); err != nil {
			return err
		}
		p.ctx.Events().Emit(types.EventTypeDatabaseMigrationComplete)
		return nil
	}

	conn, ok, err = p.ctx.Database().PostgreSQL().Conn()
	if err != nil {
		return err
	}
	if ok {
		p.conn = conn
		p.queries = newQuery(p.ctx, p.conn, types.PostgreSQLDialect)
		if err = p.migrate(); err != nil {
			return err
		}
		p.ctx.Events().Emit(types.EventTypeDatabaseMigrationComplete)
		return nil
	}
	return nil
}

func (p *Provider) migrate() error {
	if err := p.migrateLegacyTables(); err != nil {
		return err
	}
	if err := p.migrateLegacyLocationCodeIndexes(); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.User{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Location{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Customer{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Supplier{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Order{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Shipment{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Team{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Task{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.TeamMember{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.TaskUserAssignment{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.TaskTeamAssignment{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Product{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductVariant{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductVariantMedia{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductBrand{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductCategory{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductSubcategory{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductAttributeDefinition{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductAttribute{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.ProductAttributeOption{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.LocationAisle{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.LocationBay{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.LocationShelf{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.LocationShelfLevel{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.LocationBin{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.Inventory{}); err != nil {
		return err
	}
	if err := p.conn.AutoMigrate(&entities.InventoryMovement{}); err != nil {
		return err
	}
	return nil
}

func (p *Provider) dropLegacyIndex(entity interface{}, names ...string) error {
	for _, name := range names {
		if p.conn.Migrator().HasIndex(entity, name) {
			if err := p.conn.Migrator().DropIndex(entity, name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Provider) migrateLegacyTables() error {
	legacyInventoryTable := "inventories"
	inventoryTable := (&entities.Inventory{}).TableName()

	if p.conn.Migrator().HasTable(legacyInventoryTable) && !p.conn.Migrator().HasTable(inventoryTable) {
		if err := p.conn.Migrator().RenameTable(legacyInventoryTable, inventoryTable); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) migrateLegacyLocationCodeIndexes() error {
	if err := p.dropLegacyIndex(
		&entities.LocationAisle{},
		"idx_location_aisles_code",
		"uni_location_aisles_code",
		"location_aisles_code_key",
	); err != nil {
		return err
	}
	if err := p.dropLegacyIndex(
		&entities.LocationBay{},
		"idx_location_bays_code",
		"location_bays_code_key",
	); err != nil {
		return err
	}
	if err := p.dropLegacyIndex(
		&entities.LocationShelf{},
		"idx_locations_shelves_code",
		"locations_shelves_code_key",
	); err != nil {
		return err
	}
	if err := p.dropLegacyIndex(
		&entities.LocationShelfLevel{},
		"idx_location_shelf_levels_code",
		"location_shelf_levels_code_key",
	); err != nil {
		return err
	}
	if err := p.dropLegacyIndex(
		&entities.LocationBin{},
		"idx_location_bins_code",
		"location_bins_code_key",
	); err != nil {
		return err
	}
	return nil
}

func (p *Provider) Query() (*Query, bool) {
	return p.queries, p.queries != nil
}
