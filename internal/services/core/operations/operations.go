package operations

import (
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database/queries"
	"strings"
)

type Service struct {
	ctx       types.IContext
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func (s *Service) AddOrder(data AddOrderData) (*AddOrderResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}

	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddOrderResult{Code: types.ServiceResultCodeFailed}, nil
	}

	locationID := strings.TrimSpace(data.LocationID)
	if !s.locationExists(query, locationID) {
		return &AddOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	orderNumber := strings.TrimSpace(data.OrderNumber)
	customerName := strings.TrimSpace(data.CustomerName)
	channel := strings.TrimSpace(data.Channel)
	if orderNumber == "" {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"order_number": {"cannot be empty"}}}, nil
	}
	if customerName == "" {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"customer_name": {"cannot be empty"}}}, nil
	}
	if channel == "" {
		channel = "manual"
	}
	if data.TotalAmount < 0 {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"total_amount": {"must be greater than or equal to 0"}}}, nil
	}

	exists, err := query.Orders().Exists(queries.OrdersParams{LocationID: &locationID, OrderNumber: &orderNumber})
	if err != nil {
		return nil, err
	}
	if exists {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"order_number": {"already exists for this location"}}}, nil
	}

	status := data.Status
	if status == "" {
		status = entities.OrderStatusNew
	}
	if !validOrderStatus(status) {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
	}

	priority := data.Priority
	if priority == "" {
		priority = entities.OrderPriorityNormal
	}
	if !validOrderPriority(priority) {
		return &AddOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"priority": {"invalid value"}}}, nil
	}

	order, err := query.Orders().Create(entities.Order{
		LocationID:   locationID,
		OrderNumber:  orderNumber,
		CustomerName: customerName,
		Channel:      channel,
		Status:       status,
		Priority:     priority,
		ItemsCount:   data.ItemsCount,
		TotalAmount:  data.TotalAmount,
		DueAt:        data.DueAt,
		Notes:        trimStringPtr(data.Notes),
	})
	if err != nil {
		return nil, err
	}

	return &AddOrderResult{Code: types.ServiceResultCodeSuccess, Payload: AddOrderResultPayload{Order: order}}, nil
}

func (s *Service) UpdateOrder(data UpdateOrderData) (*UpdateOrderResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}

	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateOrderResult{Code: types.ServiceResultCodeFailed}, nil
	}

	locationID := strings.TrimSpace(data.LocationID)
	current, exists, err := query.Orders().FindOne(queries.OrdersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	hasChanges := false
	update := entities.Order{}

	if data.OrderNumber != nil {
		hasChanges = true
		orderNumber := strings.TrimSpace(*data.OrderNumber)
		if orderNumber == "" {
			return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"order_number": {"cannot be empty"}}}, nil
		}
		if orderNumber != current.OrderNumber {
			existing, exists, err := query.Orders().FindOne(queries.OrdersParams{LocationID: &locationID, OrderNumber: &orderNumber})
			if err != nil {
				return nil, err
			}
			if exists && existing.ID != current.ID {
				return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"order_number": {"already exists for this location"}}}, nil
			}
		}
		update.OrderNumber = orderNumber
	}

	if data.CustomerName != nil {
		hasChanges = true
		customerName := strings.TrimSpace(*data.CustomerName)
		if customerName == "" {
			return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"customer_name": {"cannot be empty"}}}, nil
		}
		update.CustomerName = customerName
	}

	if data.Channel != nil {
		hasChanges = true
		channel := strings.TrimSpace(*data.Channel)
		if channel == "" {
			channel = "manual"
		}
		update.Channel = channel
	}

	if data.Status != nil {
		hasChanges = true
		if !validOrderStatus(*data.Status) {
			return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
		}
		update.Status = *data.Status
	}

	if data.Priority != nil {
		hasChanges = true
		if !validOrderPriority(*data.Priority) {
			return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"priority": {"invalid value"}}}, nil
		}
		update.Priority = *data.Priority
	}

	if data.ItemsCount != nil {
		hasChanges = true
		update.ItemsCount = *data.ItemsCount
	}

	if data.TotalAmount != nil {
		hasChanges = true
		if *data.TotalAmount < 0 {
			return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"total_amount": {"must be greater than or equal to 0"}}}, nil
		}
		update.TotalAmount = *data.TotalAmount
	}

	if data.DueAt != nil {
		hasChanges = true
		update.DueAt = data.DueAt
	}

	if data.Notes != nil {
		hasChanges = true
		update.Notes = trimStringPtr(data.Notes)
	}

	if !hasChanges {
		return &UpdateOrderResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}

	_, updated, err := query.Orders().Update(queries.OrdersParams{ID: &data.ID, LocationID: &locationID}, update)
	if err != nil {
		return nil, err
	}
	if !updated {
		return &UpdateOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	order, exists, err := query.Orders().FindOne(queries.OrdersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateOrderResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateOrderResultPayload{Order: order}}, nil
}

func (s *Service) DeleteOrder(data DeleteOrderData) (*DeleteOrderResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteOrderResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}

	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteOrderResult{Code: types.ServiceResultCodeFailed}, nil
	}

	locationID := strings.TrimSpace(data.LocationID)
	_, exists, err := query.Orders().FindOne(queries.OrdersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	deleted, err := query.Orders().Delete(queries.OrdersParams{ID: &data.ID})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteOrderResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteOrderResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetOrders(data GetOrdersData) (*GetOrdersResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetOrdersResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)

	statusFilter, valid := sanitizeOrderStatuses(data.Status)
	if !valid {
		return &GetOrdersResult{Code: types.ServiceResultCodeInvalid}, nil
	}
	priorityFilter, valid := sanitizeOrderPriorities(data.Priority)
	if !valid {
		return &GetOrdersResult{Code: types.ServiceResultCodeInvalid}, nil
	}

	params := queries.OrdersParams{LocationID: &locationID, Status: statusFilter, Priority: priorityFilter}
	results, paging, err := query.Orders().FindMany(
		params,
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.fetchAllOrders(locationID, statusFilter, priorityFilter)
	if err != nil {
		return nil, err
	}
	summary := OrdersSummary{}
	for _, item := range rows {
		summary.TotalOrders++
		summary.TotalOrderValue += item.TotalAmount
		if item.Priority == entities.OrderPriorityUrgent {
			summary.UrgentOrders++
		}
		if item.Status != entities.OrderStatusShipped && item.Status != entities.OrderStatusCancelled {
			summary.OpenOrders++
		}
	}

	return &GetOrdersResult{Code: types.ServiceResultCodeSuccess, Payload: GetOrdersResultPayload{Orders: results, Summary: summary}, Pagination: paging}, nil
}

func (s *Service) AddShipment(data AddShipmentData) (*AddShipmentResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}

	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddShipmentResult{Code: types.ServiceResultCodeFailed}, nil
	}

	locationID := strings.TrimSpace(data.LocationID)
	if !s.locationExists(query, locationID) {
		return &AddShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	shipmentNumber := strings.TrimSpace(data.ShipmentNumber)
	if shipmentNumber == "" {
		return &AddShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"shipment_number": {"cannot be empty"}}}, nil
	}

	exists, err := query.Shipments().Exists(queries.ShipmentsParams{LocationID: &locationID, ShipmentNumber: &shipmentNumber})
	if err != nil {
		return nil, err
	}
	if exists {
		return &AddShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"shipment_number": {"already exists for this location"}}}, nil
	}

	status := data.Status
	if status == "" {
		status = entities.ShipmentStatusDraft
	}
	if !validShipmentStatus(status) {
		return &AddShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
	}

	shipment, err := query.Shipments().Create(entities.Shipment{
		LocationID:     locationID,
		ShipmentNumber: shipmentNumber,
		OrderNumber:    strings.TrimSpace(data.OrderNumber),
		Carrier:        strings.TrimSpace(data.Carrier),
		Service:        strings.TrimSpace(data.Service),
		Status:         status,
		Packages:       data.Packages,
		TrackingCode:   strings.TrimSpace(data.TrackingCode),
		ETA:            data.ETA,
		Destination:    strings.TrimSpace(data.Destination),
	})
	if err != nil {
		return nil, err
	}
	return &AddShipmentResult{Code: types.ServiceResultCodeSuccess, Payload: AddShipmentResultPayload{Shipment: shipment}}, nil
}

func (s *Service) UpdateShipment(data UpdateShipmentData) (*UpdateShipmentResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}

	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeFailed}, nil
	}

	locationID := strings.TrimSpace(data.LocationID)
	current, exists, err := query.Shipments().FindOne(queries.ShipmentsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	hasChanges := false
	update := entities.Shipment{}

	if data.ShipmentNumber != nil {
		hasChanges = true
		shipmentNumber := strings.TrimSpace(*data.ShipmentNumber)
		if shipmentNumber == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"shipment_number": {"cannot be empty"}}}, nil
		}
		if shipmentNumber != current.ShipmentNumber {
			existing, exists, err := query.Shipments().FindOne(queries.ShipmentsParams{LocationID: &locationID, ShipmentNumber: &shipmentNumber})
			if err != nil {
				return nil, err
			}
			if exists && existing.ID != current.ID {
				return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"shipment_number": {"already exists for this location"}}}, nil
			}
		}
		update.ShipmentNumber = shipmentNumber
	}

	if data.OrderNumber != nil {
		hasChanges = true
		orderNumber := strings.TrimSpace(*data.OrderNumber)
		if orderNumber == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"order_number": {"cannot be empty"}}}, nil
		}
		update.OrderNumber = orderNumber
	}

	if data.Carrier != nil {
		hasChanges = true
		carrier := strings.TrimSpace(*data.Carrier)
		if carrier == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"carrier": {"cannot be empty"}}}, nil
		}
		update.Carrier = carrier
	}

	if data.Service != nil {
		hasChanges = true
		service := strings.TrimSpace(*data.Service)
		if service == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"service": {"cannot be empty"}}}, nil
		}
		update.Service = service
	}

	if data.Status != nil {
		hasChanges = true
		if !validShipmentStatus(*data.Status) {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
		}
		update.Status = *data.Status
	}

	if data.Packages != nil {
		hasChanges = true
		update.Packages = *data.Packages
	}

	if data.TrackingCode != nil {
		hasChanges = true
		trackingCode := strings.TrimSpace(*data.TrackingCode)
		if trackingCode == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"tracking_code": {"cannot be empty"}}}, nil
		}
		update.TrackingCode = trackingCode
	}

	if data.ETA != nil {
		hasChanges = true
		update.ETA = data.ETA
	}

	if data.Destination != nil {
		hasChanges = true
		destination := strings.TrimSpace(*data.Destination)
		if destination == "" {
			return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"destination": {"cannot be empty"}}}, nil
		}
		update.Destination = destination
	}

	if !hasChanges {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}

	_, updated, err := query.Shipments().Update(queries.ShipmentsParams{ID: &data.ID, LocationID: &locationID}, update)
	if err != nil {
		return nil, err
	}
	if !updated {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	shipment, exists, err := query.Shipments().FindOne(queries.ShipmentsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateShipmentResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateShipmentResultPayload{Shipment: shipment}}, nil
}

func (s *Service) DeleteShipment(data DeleteShipmentData) (*DeleteShipmentResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteShipmentResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteShipmentResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	_, exists, err := query.Shipments().FindOne(queries.ShipmentsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	deleted, err := query.Shipments().Delete(queries.ShipmentsParams{ID: &data.ID})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteShipmentResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteShipmentResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetShipments(data GetShipmentsData) (*GetShipmentsResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetShipmentsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	statusFilter, valid := sanitizeShipmentStatuses(data.Status)
	if !valid {
		return &GetShipmentsResult{Code: types.ServiceResultCodeInvalid}, nil
	}
	params := queries.ShipmentsParams{LocationID: &locationID, Status: statusFilter}
	results, paging, err := query.Shipments().FindMany(
		params,
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.fetchAllShipments(locationID, statusFilter)
	if err != nil {
		return nil, err
	}
	summary := ShipmentsSummary{}
	for _, item := range rows {
		summary.TotalShipments++
		summary.TotalPackages += item.Packages
		switch item.Status {
		case entities.ShipmentStatusInTransit:
			summary.InTransit++
		case entities.ShipmentStatusDelivered:
			summary.Delivered++
		case entities.ShipmentStatusException:
			summary.Exceptions++
		}
	}

	return &GetShipmentsResult{Code: types.ServiceResultCodeSuccess, Payload: GetShipmentsResultPayload{Shipments: results, Summary: summary}, Pagination: paging}, nil
}

func (s *Service) AddCustomer(data AddCustomerData) (*AddCustomerResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddCustomerResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	if !s.locationExists(query, locationID) {
		return &AddCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	tier := data.Tier
	if tier == "" {
		tier = entities.CustomerTierStandard
	}
	if !validCustomerTier(tier) {
		return &AddCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"tier": {"invalid value"}}}, nil
	}

	status := data.Status
	if status == "" {
		status = entities.CustomerStatusActive
	}
	if !validCustomerStatus(status) {
		return &AddCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
	}

	if data.LifetimeValue < 0 {
		return &AddCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"lifetime_value": {"must be greater than or equal to 0"}}}, nil
	}

	email := strings.TrimSpace(strings.ToLower(data.Email))
	existing, err := query.Customers().Exists(queries.CustomersParams{LocationID: &locationID, Email: &email})
	if err != nil {
		return nil, err
	}
	if existing {
		return &AddCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"email": {"already exists for this location"}}}, nil
	}

	customer, err := query.Customers().Create(entities.Customer{
		LocationID:    locationID,
		Name:          strings.TrimSpace(data.Name),
		Email:         email,
		Phone:         strings.TrimSpace(data.Phone),
		City:          strings.TrimSpace(data.City),
		Tier:          tier,
		Status:        status,
		TotalOrders:   data.TotalOrders,
		LifetimeValue: data.LifetimeValue,
		LastOrderAt:   data.LastOrderAt,
	})
	if err != nil {
		return nil, err
	}
	return &AddCustomerResult{Code: types.ServiceResultCodeSuccess, Payload: AddCustomerResultPayload{Customer: customer}}, nil
}

func (s *Service) UpdateCustomer(data UpdateCustomerData) (*UpdateCustomerResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	current, exists, err := query.Customers().FindOne(queries.CustomersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	hasChanges := false
	update := entities.Customer{}

	if data.Name != nil {
		hasChanges = true
		name := strings.TrimSpace(*data.Name)
		if name == "" {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"cannot be empty"}}}, nil
		}
		update.Name = name
	}

	if data.Email != nil {
		hasChanges = true
		email := strings.TrimSpace(strings.ToLower(*data.Email))
		if email == "" {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"email": {"cannot be empty"}}}, nil
		}
		if email != current.Email {
			existing, exists, err := query.Customers().FindOne(queries.CustomersParams{LocationID: &locationID, Email: &email})
			if err != nil {
				return nil, err
			}
			if exists && existing.ID != current.ID {
				return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"email": {"already exists for this location"}}}, nil
			}
		}
		update.Email = email
	}

	if data.Phone != nil {
		hasChanges = true
		phone := strings.TrimSpace(*data.Phone)
		if phone == "" {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"phone": {"cannot be empty"}}}, nil
		}
		update.Phone = phone
	}

	if data.City != nil {
		hasChanges = true
		city := strings.TrimSpace(*data.City)
		if city == "" {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"city": {"cannot be empty"}}}, nil
		}
		update.City = city
	}

	if data.Tier != nil {
		hasChanges = true
		if !validCustomerTier(*data.Tier) {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"tier": {"invalid value"}}}, nil
		}
		update.Tier = *data.Tier
	}

	if data.Status != nil {
		hasChanges = true
		if !validCustomerStatus(*data.Status) {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
		}
		update.Status = *data.Status
	}

	if data.TotalOrders != nil {
		hasChanges = true
		update.TotalOrders = *data.TotalOrders
	}

	if data.LifetimeValue != nil {
		hasChanges = true
		if *data.LifetimeValue < 0 {
			return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"lifetime_value": {"must be greater than or equal to 0"}}}, nil
		}
		update.LifetimeValue = *data.LifetimeValue
	}

	if data.LastOrderAt != nil {
		hasChanges = true
		update.LastOrderAt = data.LastOrderAt
	}

	if !hasChanges {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}

	_, updated, err := query.Customers().Update(queries.CustomersParams{ID: &data.ID, LocationID: &locationID}, update)
	if err != nil {
		return nil, err
	}
	if !updated {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	customer, exists, err := query.Customers().FindOne(queries.CustomersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateCustomerResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateCustomerResultPayload{Customer: customer}}, nil
}

func (s *Service) DeleteCustomer(data DeleteCustomerData) (*DeleteCustomerResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteCustomerResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteCustomerResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	_, exists, err := query.Customers().FindOne(queries.CustomersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	deleted, err := query.Customers().Delete(queries.CustomersParams{ID: &data.ID})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteCustomerResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteCustomerResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetCustomers(data GetCustomersData) (*GetCustomersResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetCustomersResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	tierFilter, valid := sanitizeCustomerTiers(data.Tier)
	if !valid {
		return &GetCustomersResult{Code: types.ServiceResultCodeInvalid}, nil
	}
	statusFilter, valid := sanitizeCustomerStatuses(data.Status)
	if !valid {
		return &GetCustomersResult{Code: types.ServiceResultCodeInvalid}, nil
	}

	params := queries.CustomersParams{LocationID: &locationID, Tier: tierFilter, Status: statusFilter}
	results, paging, err := query.Customers().FindMany(
		params,
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.fetchAllCustomers(locationID, tierFilter, statusFilter)
	if err != nil {
		return nil, err
	}
	summary := CustomersSummary{}
	for _, item := range rows {
		summary.TotalCustomers++
		if item.Status == entities.CustomerStatusActive {
			summary.ActiveCount++
		}
		if item.Tier == entities.CustomerTierEnterprise {
			summary.EnterpriseCount++
		}
		summary.AvgLTV += item.LifetimeValue
	}
	if summary.TotalCustomers > 0 {
		summary.AvgLTV = summary.AvgLTV / float64(summary.TotalCustomers)
	}

	return &GetCustomersResult{Code: types.ServiceResultCodeSuccess, Payload: GetCustomersResultPayload{Customers: results, Summary: summary}, Pagination: paging}, nil
}

func (s *Service) AddSupplier(data AddSupplierData) (*AddSupplierResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddSupplierResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	if !s.locationExists(query, locationID) {
		return &AddSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	supplierType := data.SupplierType
	if supplierType == "" {
		supplierType = entities.SupplierTypeService
	}
	if !validSupplierType(supplierType) {
		return &AddSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"supplier_type": {"invalid value"}}}, nil
	}

	status := data.Status
	if status == "" {
		status = entities.SupplierStatusApproved
	}
	if !validSupplierStatus(status) {
		return &AddSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
	}

	if data.OnTimeRate < 0 || data.OnTimeRate > 100 {
		return &AddSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"on_time_rate": {"must be between 0 and 100"}}}, nil
	}

	name := strings.TrimSpace(data.Name)
	existing, err := query.Suppliers().Exists(queries.SuppliersParams{LocationID: &locationID, Name: &name})
	if err != nil {
		return nil, err
	}
	if existing {
		return &AddSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"already exists for this location"}}}, nil
	}

	supplier, err := query.Suppliers().Create(entities.Supplier{
		LocationID:         locationID,
		Name:               name,
		SupplierType:       supplierType,
		Status:             status,
		ContactName:        strings.TrimSpace(data.ContactName),
		Email:              strings.TrimSpace(strings.ToLower(data.Email)),
		LeadTimeDays:       data.LeadTimeDays,
		OnTimeRate:         data.OnTimeRate,
		OpenPurchaseOrders: data.OpenPurchaseOrders,
		City:               strings.TrimSpace(data.City),
	})
	if err != nil {
		return nil, err
	}
	return &AddSupplierResult{Code: types.ServiceResultCodeSuccess, Payload: AddSupplierResultPayload{Supplier: supplier}}, nil
}

func (s *Service) UpdateSupplier(data UpdateSupplierData) (*UpdateSupplierResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	current, exists, err := query.Suppliers().FindOne(queries.SuppliersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	hasChanges := false
	update := entities.Supplier{}

	if data.Name != nil {
		hasChanges = true
		name := strings.TrimSpace(*data.Name)
		if name == "" {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"cannot be empty"}}}, nil
		}
		if name != current.Name {
			existing, exists, err := query.Suppliers().FindOne(queries.SuppliersParams{LocationID: &locationID, Name: &name})
			if err != nil {
				return nil, err
			}
			if exists && existing.ID != current.ID {
				return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"already exists for this location"}}}, nil
			}
		}
		update.Name = name
	}

	if data.SupplierType != nil {
		hasChanges = true
		if !validSupplierType(*data.SupplierType) {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"supplier_type": {"invalid value"}}}, nil
		}
		update.SupplierType = *data.SupplierType
	}

	if data.Status != nil {
		hasChanges = true
		if !validSupplierStatus(*data.Status) {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
		}
		update.Status = *data.Status
	}

	if data.ContactName != nil {
		hasChanges = true
		contactName := strings.TrimSpace(*data.ContactName)
		if contactName == "" {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"contact_name": {"cannot be empty"}}}, nil
		}
		update.ContactName = contactName
	}

	if data.Email != nil {
		hasChanges = true
		email := strings.TrimSpace(strings.ToLower(*data.Email))
		if email == "" {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"email": {"cannot be empty"}}}, nil
		}
		update.Email = email
	}

	if data.LeadTimeDays != nil {
		hasChanges = true
		update.LeadTimeDays = *data.LeadTimeDays
	}

	if data.OnTimeRate != nil {
		hasChanges = true
		if *data.OnTimeRate < 0 || *data.OnTimeRate > 100 {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"on_time_rate": {"must be between 0 and 100"}}}, nil
		}
		update.OnTimeRate = *data.OnTimeRate
	}

	if data.OpenPurchaseOrders != nil {
		hasChanges = true
		update.OpenPurchaseOrders = *data.OpenPurchaseOrders
	}

	if data.City != nil {
		hasChanges = true
		city := strings.TrimSpace(*data.City)
		if city == "" {
			return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"city": {"cannot be empty"}}}, nil
		}
		update.City = city
	}

	if !hasChanges {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}

	_, updated, err := query.Suppliers().Update(queries.SuppliersParams{ID: &data.ID, LocationID: &locationID}, update)
	if err != nil {
		return nil, err
	}
	if !updated {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}

	supplier, exists, err := query.Suppliers().FindOne(queries.SuppliersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateSupplierResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateSupplierResultPayload{Supplier: supplier}}, nil
}

func (s *Service) DeleteSupplier(data DeleteSupplierData) (*DeleteSupplierResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteSupplierResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteSupplierResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	_, exists, err := query.Suppliers().FindOne(queries.SuppliersParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	deleted, err := query.Suppliers().Delete(queries.SuppliersParams{ID: &data.ID})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteSupplierResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteSupplierResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetSuppliers(data GetSuppliersData) (*GetSuppliersResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetSuppliersResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := strings.TrimSpace(data.LocationID)
	statusFilter, valid := sanitizeSupplierStatuses(data.Status)
	if !valid {
		return &GetSuppliersResult{Code: types.ServiceResultCodeInvalid}, nil
	}

	params := queries.SuppliersParams{LocationID: &locationID, Status: statusFilter}
	results, paging, err := query.Suppliers().FindMany(
		params,
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.fetchAllSuppliers(locationID, statusFilter)
	if err != nil {
		return nil, err
	}
	summary := SuppliersSummary{}
	for _, item := range rows {
		summary.TotalSuppliers++
		summary.TotalOpenPOs += item.OpenPurchaseOrders
		summary.AvgLeadTime += float64(item.LeadTimeDays)
		if item.Status == entities.SupplierStatusApproved {
			summary.ApprovedCount++
		}
	}
	if summary.TotalSuppliers > 0 {
		summary.AvgLeadTime = summary.AvgLeadTime / float64(summary.TotalSuppliers)
	}

	return &GetSuppliersResult{Code: types.ServiceResultCodeSuccess, Payload: GetSuppliersResultPayload{Suppliers: results, Summary: summary}, Pagination: paging}, nil
}

func (s *Service) locationExists(query interface{ Locations() *queries.Locations }, locationID string) bool {
	exists, err := query.Locations().Exists(queries.LocationsParams{ID: &locationID})
	if err != nil {
		s.ctx.Logger().Error(err)
		return false
	}
	return exists
}

func (s *Service) fetchAllOrders(locationID string, statuses []entities.OrderStatusEnum, priorities []entities.OrderPriorityEnum) ([]entities.Order, error) {
	query, _ := s.providers.Database().Query()
	rows := []entities.Order{}
	page := uint32(1)
	for {
		results, paging, err := query.Orders().FindMany(
			queries.OrdersParams{LocationID: &locationID, Status: statuses, Priority: priorities},
			pagination.New(types.PaginationParams{Page: page, Limit: 250, Order: "created_at.desc"}),
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, results...)
		if paging == nil || !paging.HasNext {
			break
		}
		page = paging.Next
	}
	return rows, nil
}

func (s *Service) fetchAllShipments(locationID string, statuses []entities.ShipmentStatusEnum) ([]entities.Shipment, error) {
	query, _ := s.providers.Database().Query()
	rows := []entities.Shipment{}
	page := uint32(1)
	for {
		results, paging, err := query.Shipments().FindMany(
			queries.ShipmentsParams{LocationID: &locationID, Status: statuses},
			pagination.New(types.PaginationParams{Page: page, Limit: 250, Order: "created_at.desc"}),
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, results...)
		if paging == nil || !paging.HasNext {
			break
		}
		page = paging.Next
	}
	return rows, nil
}

func (s *Service) fetchAllCustomers(locationID string, tiers []entities.CustomerTierEnum, statuses []entities.CustomerStatusEnum) ([]entities.Customer, error) {
	query, _ := s.providers.Database().Query()
	rows := []entities.Customer{}
	page := uint32(1)
	for {
		results, paging, err := query.Customers().FindMany(
			queries.CustomersParams{LocationID: &locationID, Tier: tiers, Status: statuses},
			pagination.New(types.PaginationParams{Page: page, Limit: 250, Order: "created_at.desc"}),
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, results...)
		if paging == nil || !paging.HasNext {
			break
		}
		page = paging.Next
	}
	return rows, nil
}

func (s *Service) fetchAllSuppliers(locationID string, statuses []entities.SupplierStatusEnum) ([]entities.Supplier, error) {
	query, _ := s.providers.Database().Query()
	rows := []entities.Supplier{}
	page := uint32(1)
	for {
		results, paging, err := query.Suppliers().FindMany(
			queries.SuppliersParams{LocationID: &locationID, Status: statuses},
			pagination.New(types.PaginationParams{Page: page, Limit: 250, Order: "created_at.desc"}),
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, results...)
		if paging == nil || !paging.HasNext {
			break
		}
		page = paging.Next
	}
	return rows, nil
}

func sanitizeOrderStatuses(values []entities.OrderStatusEnum) ([]entities.OrderStatusEnum, bool) {
	result := []entities.OrderStatusEnum{}
	for _, value := range values {
		if !validOrderStatus(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func sanitizeOrderPriorities(values []entities.OrderPriorityEnum) ([]entities.OrderPriorityEnum, bool) {
	result := []entities.OrderPriorityEnum{}
	for _, value := range values {
		if !validOrderPriority(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func sanitizeShipmentStatuses(values []entities.ShipmentStatusEnum) ([]entities.ShipmentStatusEnum, bool) {
	result := []entities.ShipmentStatusEnum{}
	for _, value := range values {
		if !validShipmentStatus(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func sanitizeCustomerTiers(values []entities.CustomerTierEnum) ([]entities.CustomerTierEnum, bool) {
	result := []entities.CustomerTierEnum{}
	for _, value := range values {
		if !validCustomerTier(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func sanitizeCustomerStatuses(values []entities.CustomerStatusEnum) ([]entities.CustomerStatusEnum, bool) {
	result := []entities.CustomerStatusEnum{}
	for _, value := range values {
		if !validCustomerStatus(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func sanitizeSupplierStatuses(values []entities.SupplierStatusEnum) ([]entities.SupplierStatusEnum, bool) {
	result := []entities.SupplierStatusEnum{}
	for _, value := range values {
		if !validSupplierStatus(value) {
			return nil, false
		}
		result = append(result, value)
	}
	return result, true
}

func validOrderStatus(status entities.OrderStatusEnum) bool {
	switch status {
	case entities.OrderStatusNew, entities.OrderStatusPicking, entities.OrderStatusPacked, entities.OrderStatusShipped, entities.OrderStatusBackorder, entities.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func validOrderPriority(priority entities.OrderPriorityEnum) bool {
	switch priority {
	case entities.OrderPriorityLow, entities.OrderPriorityNormal, entities.OrderPriorityHigh, entities.OrderPriorityUrgent:
		return true
	default:
		return false
	}
}

func validShipmentStatus(status entities.ShipmentStatusEnum) bool {
	switch status {
	case entities.ShipmentStatusDraft, entities.ShipmentStatusReady, entities.ShipmentStatusInTransit, entities.ShipmentStatusDelivered, entities.ShipmentStatusException:
		return true
	default:
		return false
	}
}

func validCustomerTier(tier entities.CustomerTierEnum) bool {
	switch tier {
	case entities.CustomerTierStandard, entities.CustomerTierGrowth, entities.CustomerTierEnterprise:
		return true
	default:
		return false
	}
}

func validCustomerStatus(status entities.CustomerStatusEnum) bool {
	switch status {
	case entities.CustomerStatusActive, entities.CustomerStatusPaused, entities.CustomerStatusChurnRisk:
		return true
	default:
		return false
	}
}

func validSupplierType(supplierType entities.SupplierTypeEnum) bool {
	switch supplierType {
	case entities.SupplierTypeRawMaterial, entities.SupplierTypePackaging, entities.SupplierTypeFinishedGood, entities.SupplierTypeService:
		return true
	default:
		return false
	}
}

func validSupplierStatus(status entities.SupplierStatusEnum) bool {
	switch status {
	case entities.SupplierStatusApproved, entities.SupplierStatusProbation, entities.SupplierStatusBlocked:
		return true
	default:
		return false
	}
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
