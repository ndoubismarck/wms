package app

import (
	"server/internal/core/shared/types"
	"server/internal/server/helpers"
	"server/internal/server/middleware"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type handler struct {
	ctx      types.IContext
	helpers  *helpers.Helper
	services *services.Services
}

type middlewares struct {
	auth         gin.HandlerFunc
	setup        gin.HandlerFunc
	owner        gin.HandlerFunc
	admin        gin.HandlerFunc
	ownerOrAdmin gin.HandlerFunc
}

type apiResponseBody types.ApiResponseBody

func Setup(ctx types.IContext, engine *gin.Engine, services *services.Services, helpers *helpers.Helper, middleware *middleware.Middleware, eventsChan chan *types.ServerEvent) error {
	h := &handler{
		ctx:      ctx,
		helpers:  helpers,
		services: services,
	}

	mw := middlewares{
		auth:  middleware.Auth(),
		setup: middleware.Setup(),
		owner: middleware.OnlyOwner(),
		admin: middleware.AllowedRoles([]string{"admin"}),
	}

	router := engine.Group("/app/v1")

	router.GET("/me", mw.setup, mw.auth, h.getAuthUser)
	router.PUT("/me", mw.setup, mw.auth, h.updateAuthUser)

	router.GET("/sse", h.sse)

	router.GET("setup", h.getSetup)
	router.POST("/setup", h.createSetup)

	router.POST("/auth/login", mw.setup, h.authLogin)

	router.GET("/:lid/stats", mw.setup, mw.auth, h.getStats)

	router.POST("/:lid/files/upload", mw.setup, mw.auth, h.uploadFile)

	router.GET("/:lid/users", mw.setup, mw.auth, h.getUsers)
	router.POST("/:lid/users", mw.setup, mw.auth, h.addUser)
	router.PUT("/:lid/users/:id", mw.setup, mw.auth, h.updateUser)
	router.DELETE("/:lid/users/:id", mw.setup, mw.auth, mw.admin, h.deleteUser)
	router.GET("/:lid/products", mw.setup, mw.auth, h.getProducts)
	router.POST("/:lid/products", mw.setup, mw.auth, h.addProduct)
	router.PUT("/:lid/products/:id", mw.setup, mw.auth, h.updateProduct)
	router.GET("/:lid/products/:id", mw.setup, mw.auth, h.getProduct)
	router.DELETE("/:lid/products/:id", mw.setup, mw.auth, h.deleteProduct)
	router.GET("/:lid/products/:id/barcode", mw.setup, h.getProductBarcode)
	router.GET("/:lid/products/categories", mw.setup, mw.auth, h.getProductCategories)
	router.POST("/:lid/products/categories", mw.setup, mw.auth, h.addProductCategory)
	router.PUT("/:lid/products/categories/:id", mw.setup, mw.auth, h.updateProductCategory)
	router.GET("/:lid/products/categories/:id", mw.setup, mw.auth, h.getProductCategory)
	router.DELETE("/:lid/products/categories/:id", mw.setup, mw.auth, h.deleteProductCategory)
	router.GET("/:lid/products/subcategories", mw.setup, mw.auth, h.getProductSubcategories)
	router.POST("/:lid/products/subcategories", mw.setup, mw.auth, h.addProductSubcategory)
	router.PUT("/:lid/products/subcategories/:id", mw.setup, mw.auth, h.updateProductSubcategory)
	router.GET("/:lid/products/subcategories/:id", mw.setup, mw.auth, h.getProductSubcategory)
	router.DELETE("/:lid/products/subcategories/:id", mw.setup, mw.auth, h.deleteProductSubcategory)
	router.GET("/:lid/products/brands", mw.setup, mw.auth, h.getProductBrands)
	router.POST("/:lid/products/brands", mw.setup, mw.auth, h.addProductBrand)
	router.PUT("/:lid/products/brands/:id", mw.setup, mw.auth, h.updateProductBrand)
	router.GET("/:lid/products/brands/:id", mw.setup, mw.auth, h.getProductBrand)
	router.DELETE("/:lid/products/brands/:id", mw.setup, mw.auth, h.deleteProductBrand)
	router.GET("/:lid/products/attributes", mw.setup, mw.auth, h.getProductAttributes)
	router.POST("/:lid/products/attributes", mw.setup, mw.auth, h.addProductAttribute)
	router.PUT("/:lid/products/attributes/:id", mw.setup, mw.auth, h.updateProductAttribute)
	router.GET("/:lid/products/attributes/:id", mw.setup, mw.auth, h.getProductAttribute)
	router.DELETE("/:lid/products/attributes/:id", mw.setup, mw.auth, h.deleteProductAttribute)

	router.GET("/locations", mw.setup, mw.auth, h.getLocations)
	router.GET("/:lid/locations/aisles", mw.setup, mw.auth, h.getLocationAisles)
	router.GET("/:lid/locations/aisles/next-code", mw.setup, mw.auth, h.getLocationNextAisleCode)
	router.POST("/:lid/locations/aisles", mw.setup, mw.auth, h.addLocationAisle)
	router.PUT("/:lid/locations/aisles/:id", mw.setup, mw.auth, h.updateLocationAisle)
	router.GET("/:lid/locations/aisles/:id", mw.setup, mw.auth, h.getLocationAisle)
	router.DELETE("/:lid/locations/aisles/:id", mw.setup, mw.auth, h.deleteLocationAisle)

	router.GET("/:lid/locations/bays", mw.setup, mw.auth, h.getLocationBays)
	router.GET("/:lid/locations/bays/next-code", mw.setup, mw.auth, h.getLocationNextBayCode)
	router.POST("/:lid/locations/bays", mw.setup, mw.auth, h.addLocationBay)
	router.PUT("/:lid/locations/bays/:id", mw.setup, mw.auth, h.updateLocationBay)
	router.GET("/:lid/locations/bays/:id", mw.setup, mw.auth, h.getLocationBay)
	router.DELETE("/:lid/locations/bays/:id", mw.setup, mw.auth, h.deleteLocationBay)

	router.GET("/:lid/locations/shelves", mw.setup, mw.auth, h.getLocationShelves)
	router.GET("/:lid/locations/shelves/next-code", mw.setup, mw.auth, h.getLocationNextShelfCode)
	router.POST("/:lid/locations/shelves", mw.setup, mw.auth, h.addLocationShelf)
	router.GET("/:lid/locations/shelves/levels", mw.setup, mw.auth, h.getLocationShelfLevels)
	router.GET("/:lid/locations/shelves/levels/next-code", mw.setup, mw.auth, h.getLocationNextShelfLevelCode)
	router.POST("/:lid/locations/shelves/levels", mw.setup, mw.auth, h.addLocationShelfLevel)
	router.PUT("/:lid/locations/shelves/levels/:id", mw.setup, mw.auth, h.updateLocationShelfLevel)
	router.GET("/:lid/locations/shelves/levels/:id", mw.setup, mw.auth, h.getLocationShelfLevel)
	router.DELETE("/:lid/locations/shelves/levels/:id", mw.setup, mw.auth, h.deleteLocationShelfLevel)
	router.PUT("/:lid/locations/shelves/:id", mw.setup, mw.auth, h.updateLocationShelf)
	router.GET("/:lid/locations/shelves/:id", mw.setup, mw.auth, h.getLocationShelf)
	router.DELETE("/:lid/locations/shelves/:id", mw.setup, mw.auth, h.deleteLocationShelf)

	router.GET("/:lid/locations/bins", mw.setup, mw.auth, h.getLocationBins)
	router.GET("/:lid/locations/bins/next-code", mw.setup, mw.auth, h.getLocationNextBinCode)
	router.POST("/:lid/locations/bins", mw.setup, mw.auth, h.addLocationBin)
	router.PUT("/:lid/locations/bins/:id", mw.setup, mw.auth, h.updateLocationBin)
	router.GET("/:lid/locations/bins/:id", mw.setup, mw.auth, h.getLocationBin)
	router.DELETE("/:lid/locations/bins/:id", mw.setup, mw.auth, h.deleteLocationBin)

	router.GET("/:lid/inventory", mw.setup, mw.auth, h.getInventory)
	router.POST("/:lid/inventory", mw.setup, mw.auth, h.addInventory)
	router.GET("/:lid/inventory/summary", mw.setup, mw.auth, h.getInventorySummary)
	router.GET("/:lid/inventory/movements", mw.setup, mw.auth, h.getInventoryMovements)
	router.POST("/:lid/inventory/movements", mw.setup, mw.auth, h.addInventoryMovement)
	router.PUT("/:lid/inventory/movements/:id", mw.setup, mw.auth, h.updateInventoryMovement)
	router.GET("/:lid/inventory/movements/:id", mw.setup, mw.auth, h.getInventoryMovement)
	router.DELETE("/:lid/inventory/movements/:id", mw.setup, mw.auth, h.deleteInventoryMovement)
	router.PUT("/:lid/inventory/:id", mw.setup, mw.auth, h.updateInventory)
	router.GET("/:lid/inventory/:id", mw.setup, mw.auth, h.getInventoryItem)
	router.DELETE("/:lid/inventory/:id", mw.setup, mw.auth, h.deleteInventory)

	go func(eventsChan chan *types.ServerEvent) {
		for event := range eventsChan {
			h.helpers.SSE().Publish(event)
		}
	}(eventsChan)
	return nil
}
