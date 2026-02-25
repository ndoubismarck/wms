package database

import (
	"server/internal/core/shared/types"
	"server/internal/services/providers/database/queries"

	"gorm.io/gorm"
)

type Query struct {
	users                       *queries.Users
	locations                   *queries.Locations
	products                    *queries.Products
	productVariants             *queries.ProductVariants
	productVariantMedia         *queries.ProductVariantMedia
	inventory                   *queries.Inventory
	inventoryMovements          *queries.InventoryMovements
	locationBins                *queries.LocationBins
	locationBays                *queries.LocationBays
	locationAisles              *queries.LocationAisles
	locationShelves             *queries.LocationShelves
	locationShelfLevels         *queries.LocationShelfLevels
	productCategories           *queries.ProductCategories
	productSubcategories        *queries.ProductSubcategories
	productBrands               *queries.ProductBrands
	productAttributes           *queries.ProductAttributes
	productAttributeDefinitions *queries.ProductAttributeDefinitions
	productAttributeOptions     *queries.ProductAttributeOptions
}

func newQuery(ctx types.IContext, conn *gorm.DB, dialect types.Dialect) *Query {
	return &Query{
		users:                       queries.NewUsers(ctx, conn, dialect),
		locations:                   queries.NewLocations(ctx, conn, dialect),
		products:                    queries.NewProducts(ctx, conn, dialect),
		productVariants:             queries.NewProductVariants(ctx, conn, dialect),
		productVariantMedia:         queries.NewProductVariantMedia(ctx, conn, dialect),
		inventory:                   queries.NewInventory(ctx, conn, dialect),
		inventoryMovements:          queries.NewInventoryMovements(ctx, conn, dialect),
		locationBins:                queries.NewLocationBins(ctx, conn, dialect),
		locationBays:                queries.NewLocationBays(ctx, conn, dialect),
		locationAisles:              queries.NewLocationAisles(ctx, conn, dialect),
		locationShelves:             queries.NewLocationShelves(ctx, conn, dialect),
		locationShelfLevels:         queries.NewLocationShelfLevels(ctx, conn, dialect),
		productCategories:           queries.NewProductCategories(ctx, conn, dialect),
		productSubcategories:        queries.NewProductSubcategories(ctx, conn, dialect),
		productBrands:               queries.NewProductBrands(ctx, conn, dialect),
		productAttributes:           queries.NewProductAttributes(ctx, conn, dialect),
		productAttributeDefinitions: queries.NewProductAttributeDefinitions(ctx, conn, dialect),
		productAttributeOptions:     queries.NewProductAttributeOptions(ctx, conn, dialect),
	}
}

func (q *Query) Users() *queries.Users {
	return q.users
}

func (q *Query) Locations() *queries.Locations {
	return q.locations
}

func (q *Query) Products() *queries.Products {
	return q.products
}

func (q *Query) ProductVariants() *queries.ProductVariants {
	return q.productVariants
}

func (q *Query) ProductVariantMedia() *queries.ProductVariantMedia {
	return q.productVariantMedia
}

func (q *Query) Inventory() *queries.Inventory {
	return q.inventory
}

func (q *Query) InventoryMovements() *queries.InventoryMovements {
	return q.inventoryMovements
}

func (q *Query) LocationAisles() *queries.LocationAisles {
	return q.locationAisles
}

func (q *Query) LocationBays() *queries.LocationBays {
	return q.locationBays
}

func (q *Query) LocationShelves() *queries.LocationShelves {
	return q.locationShelves
}

func (q *Query) LocationShelfLevels() *queries.LocationShelfLevels {
	return q.locationShelfLevels
}

func (q *Query) LocationBins() *queries.LocationBins {
	return q.locationBins
}

func (q *Query) ProductCategories() *queries.ProductCategories {
	return q.productCategories
}

func (q *Query) ProductSubcategories() *queries.ProductSubcategories {
	return q.productSubcategories
}

func (q *Query) ProductBrands() *queries.ProductBrands {
	return q.productBrands
}

func (q *Query) ProductAttributes() *queries.ProductAttributes {
	return q.productAttributes
}

func (q *Query) ProductAttributeDefinitions() *queries.ProductAttributeDefinitions {
	return q.productAttributeDefinitions
}

func (q *Query) ProductAttributeOptions() *queries.ProductAttributeOptions {
	return q.productAttributeOptions
}
