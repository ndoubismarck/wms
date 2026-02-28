package operations

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"time"
)

type OrdersSummary struct {
	TotalOrders     uint64  `json:"total_orders"`
	OpenOrders      uint64  `json:"open_orders"`
	UrgentOrders    uint64  `json:"urgent_orders"`
	TotalOrderValue float64 `json:"total_order_value"`
}

type ShipmentsSummary struct {
	TotalShipments uint64 `json:"total_shipments"`
	InTransit      uint64 `json:"in_transit"`
	Delivered      uint64 `json:"delivered"`
	Exceptions     uint64 `json:"exceptions"`
	TotalPackages  uint64 `json:"total_packages"`
}

type CustomersSummary struct {
	TotalCustomers  uint64  `json:"total_customers"`
	ActiveCount     uint64  `json:"active_count"`
	EnterpriseCount uint64  `json:"enterprise_count"`
	AvgLTV          float64 `json:"avg_ltv"`
}

type SuppliersSummary struct {
	TotalSuppliers uint64  `json:"total_suppliers"`
	ApprovedCount  uint64  `json:"approved_count"`
	AvgLeadTime    float64 `json:"avg_lead_time"`
	TotalOpenPOs   uint64  `json:"total_open_pos"`
}

type (
	AddOrderData struct {
		LocationID   string                     `json:"-" validate:"required"`
		OrderNumber  string                     `json:"order_number" validate:"required"`
		CustomerName string                     `json:"customer_name" validate:"required"`
		Channel      string                     `json:"channel"`
		Status       entities.OrderStatusEnum   `json:"status"`
		Priority     entities.OrderPriorityEnum `json:"priority"`
		ItemsCount   uint64                     `json:"items_count"`
		TotalAmount  float64                    `json:"total_amount"`
		DueAt        *time.Time                 `json:"due_at,omitempty"`
		Notes        *string                    `json:"notes,omitempty"`
	}
	AddOrderResult struct {
		Code       types.ServiceResultCode
		Payload    AddOrderResultPayload
		Validation types.ValidationResult
	}
	AddOrderResultPayload struct {
		Order entities.Order `json:"order"`
	}
)

type (
	UpdateOrderData struct {
		LocationID   string                      `json:"-" validate:"required"`
		ID           string                      `json:"id" validate:"required"`
		OrderNumber  *string                     `json:"order_number,omitempty"`
		CustomerName *string                     `json:"customer_name,omitempty"`
		Channel      *string                     `json:"channel,omitempty"`
		Status       *entities.OrderStatusEnum   `json:"status,omitempty"`
		Priority     *entities.OrderPriorityEnum `json:"priority,omitempty"`
		ItemsCount   *uint64                     `json:"items_count,omitempty"`
		TotalAmount  *float64                    `json:"total_amount,omitempty"`
		DueAt        *time.Time                  `json:"due_at,omitempty"`
		Notes        *string                     `json:"notes,omitempty"`
	}
	UpdateOrderResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateOrderResultPayload
		Validation types.ValidationResult
	}
	UpdateOrderResultPayload struct {
		Order entities.Order `json:"order"`
	}
)

type (
	DeleteOrderData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteOrderResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetOrdersData struct {
		types.PaginationParams
		LocationID string                       `json:"-" form:"-" validate:"required"`
		Status     []entities.OrderStatusEnum   `json:"status" form:"status"`
		Priority   []entities.OrderPriorityEnum `json:"priority" form:"priority"`
	}
	GetOrdersResult struct {
		Code       types.ServiceResultCode
		Payload    GetOrdersResultPayload
		Pagination *types.PaginationResult
	}
	GetOrdersResultPayload struct {
		Orders  []entities.Order `json:"orders"`
		Summary OrdersSummary    `json:"summary"`
	}
)

type (
	AddShipmentData struct {
		LocationID     string                      `json:"-" validate:"required"`
		ShipmentNumber string                      `json:"shipment_number" validate:"required"`
		OrderNumber    string                      `json:"order_number" validate:"required"`
		Carrier        string                      `json:"carrier" validate:"required"`
		Service        string                      `json:"service" validate:"required"`
		Status         entities.ShipmentStatusEnum `json:"status"`
		Packages       uint64                      `json:"packages"`
		TrackingCode   string                      `json:"tracking_code" validate:"required"`
		ETA            *time.Time                  `json:"eta,omitempty"`
		Destination    string                      `json:"destination" validate:"required"`
	}
	AddShipmentResult struct {
		Code       types.ServiceResultCode
		Payload    AddShipmentResultPayload
		Validation types.ValidationResult
	}
	AddShipmentResultPayload struct {
		Shipment entities.Shipment `json:"shipment"`
	}
)

type (
	UpdateShipmentData struct {
		LocationID     string                       `json:"-" validate:"required"`
		ID             string                       `json:"id" validate:"required"`
		ShipmentNumber *string                      `json:"shipment_number,omitempty"`
		OrderNumber    *string                      `json:"order_number,omitempty"`
		Carrier        *string                      `json:"carrier,omitempty"`
		Service        *string                      `json:"service,omitempty"`
		Status         *entities.ShipmentStatusEnum `json:"status,omitempty"`
		Packages       *uint64                      `json:"packages,omitempty"`
		TrackingCode   *string                      `json:"tracking_code,omitempty"`
		ETA            *time.Time                   `json:"eta,omitempty"`
		Destination    *string                      `json:"destination,omitempty"`
	}
	UpdateShipmentResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateShipmentResultPayload
		Validation types.ValidationResult
	}
	UpdateShipmentResultPayload struct {
		Shipment entities.Shipment `json:"shipment"`
	}
)

type (
	DeleteShipmentData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteShipmentResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetShipmentsData struct {
		types.PaginationParams
		LocationID string                        `json:"-" form:"-" validate:"required"`
		Status     []entities.ShipmentStatusEnum `json:"status" form:"status"`
	}
	GetShipmentsResult struct {
		Code       types.ServiceResultCode
		Payload    GetShipmentsResultPayload
		Pagination *types.PaginationResult
	}
	GetShipmentsResultPayload struct {
		Shipments []entities.Shipment `json:"shipments"`
		Summary   ShipmentsSummary    `json:"summary"`
	}
)

type (
	AddCustomerData struct {
		LocationID    string                      `json:"-" validate:"required"`
		Name          string                      `json:"name" validate:"required"`
		Email         string                      `json:"email" validate:"required|email"`
		Phone         string                      `json:"phone" validate:"required"`
		City          string                      `json:"city" validate:"required"`
		Tier          entities.CustomerTierEnum   `json:"tier"`
		Status        entities.CustomerStatusEnum `json:"status"`
		TotalOrders   uint64                      `json:"total_orders"`
		LifetimeValue float64                     `json:"lifetime_value"`
		LastOrderAt   *time.Time                  `json:"last_order_at,omitempty"`
	}
	AddCustomerResult struct {
		Code       types.ServiceResultCode
		Payload    AddCustomerResultPayload
		Validation types.ValidationResult
	}
	AddCustomerResultPayload struct {
		Customer entities.Customer `json:"customer"`
	}
)

type (
	UpdateCustomerData struct {
		LocationID    string                       `json:"-" validate:"required"`
		ID            string                       `json:"id" validate:"required"`
		Name          *string                      `json:"name,omitempty"`
		Email         *string                      `json:"email,omitempty"`
		Phone         *string                      `json:"phone,omitempty"`
		City          *string                      `json:"city,omitempty"`
		Tier          *entities.CustomerTierEnum   `json:"tier,omitempty"`
		Status        *entities.CustomerStatusEnum `json:"status,omitempty"`
		TotalOrders   *uint64                      `json:"total_orders,omitempty"`
		LifetimeValue *float64                     `json:"lifetime_value,omitempty"`
		LastOrderAt   *time.Time                   `json:"last_order_at,omitempty"`
	}
	UpdateCustomerResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateCustomerResultPayload
		Validation types.ValidationResult
	}
	UpdateCustomerResultPayload struct {
		Customer entities.Customer `json:"customer"`
	}
)

type (
	DeleteCustomerData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteCustomerResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetCustomersData struct {
		types.PaginationParams
		LocationID string                        `json:"-" form:"-" validate:"required"`
		Tier       []entities.CustomerTierEnum   `json:"tier" form:"tier"`
		Status     []entities.CustomerStatusEnum `json:"status" form:"status"`
	}
	GetCustomersResult struct {
		Code       types.ServiceResultCode
		Payload    GetCustomersResultPayload
		Pagination *types.PaginationResult
	}
	GetCustomersResultPayload struct {
		Customers []entities.Customer `json:"customers"`
		Summary   CustomersSummary    `json:"summary"`
	}
)

type (
	AddSupplierData struct {
		LocationID         string                      `json:"-" validate:"required"`
		Name               string                      `json:"name" validate:"required"`
		SupplierType       entities.SupplierTypeEnum   `json:"supplier_type"`
		Status             entities.SupplierStatusEnum `json:"status"`
		ContactName        string                      `json:"contact_name" validate:"required"`
		Email              string                      `json:"email" validate:"required|email"`
		LeadTimeDays       uint64                      `json:"lead_time_days"`
		OnTimeRate         float64                     `json:"on_time_rate"`
		OpenPurchaseOrders uint64                      `json:"open_purchase_orders"`
		City               string                      `json:"city" validate:"required"`
	}
	AddSupplierResult struct {
		Code       types.ServiceResultCode
		Payload    AddSupplierResultPayload
		Validation types.ValidationResult
	}
	AddSupplierResultPayload struct {
		Supplier entities.Supplier `json:"supplier"`
	}
)

type (
	UpdateSupplierData struct {
		LocationID         string                       `json:"-" validate:"required"`
		ID                 string                       `json:"id" validate:"required"`
		Name               *string                      `json:"name,omitempty"`
		SupplierType       *entities.SupplierTypeEnum   `json:"supplier_type,omitempty"`
		Status             *entities.SupplierStatusEnum `json:"status,omitempty"`
		ContactName        *string                      `json:"contact_name,omitempty"`
		Email              *string                      `json:"email,omitempty"`
		LeadTimeDays       *uint64                      `json:"lead_time_days,omitempty"`
		OnTimeRate         *float64                     `json:"on_time_rate,omitempty"`
		OpenPurchaseOrders *uint64                      `json:"open_purchase_orders,omitempty"`
		City               *string                      `json:"city,omitempty"`
	}
	UpdateSupplierResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateSupplierResultPayload
		Validation types.ValidationResult
	}
	UpdateSupplierResultPayload struct {
		Supplier entities.Supplier `json:"supplier"`
	}
)

type (
	DeleteSupplierData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteSupplierResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetSuppliersData struct {
		types.PaginationParams
		LocationID string                        `json:"-" form:"-" validate:"required"`
		Status     []entities.SupplierStatusEnum `json:"status" form:"status"`
	}
	GetSuppliersResult struct {
		Code       types.ServiceResultCode
		Payload    GetSuppliersResultPayload
		Pagination *types.PaginationResult
	}
	GetSuppliersResultPayload struct {
		Suppliers []entities.Supplier `json:"suppliers"`
		Summary   SuppliersSummary    `json:"summary"`
	}
)
