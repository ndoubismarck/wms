package inventory

import (
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database/queries"
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

func (s *Service) Add(data AddData) (*AddResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.LocationBins().Exists(queries.LocationBinsParams{
		ID:         &data.BinID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.Inventory().Create(entities.Inventory{
		BinID:            data.BinID,
		VariantID:        data.VariantID,
		QuantityOnHand:   data.QuantityOnHand,
		QuantityDamaged:  data.QuantityDamaged,
		QuantityReserved: data.QuantityReserved,
		QuantityIncoming: data.QuantityIncoming,
		QuantityOutgoing: data.QuantityOutgoing,
	})
	if err != nil {
		return nil, err
	}
	return &AddResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddResultPayload{
			Inventory: result,
		},
	}, nil
}

func (s *Service) Update(data UpdateData) (*UpdateResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.Inventory().Exists(queries.InventoryParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	if data.BinID != "" {
		parentExists, err := query.LocationBins().Exists(queries.LocationBinsParams{
			ID:         &data.BinID,
			LocationID: &locationID,
		})
		if err != nil {
			return nil, err
		}
		if !parentExists {
			return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
		}
	}
	result, _, err := query.Inventory().Update(queries.InventoryParams{
		ID: &data.ID,
	}, entities.Inventory{
		BinID:            data.BinID,
		VariantID:        data.VariantID,
		QuantityOnHand:   data.QuantityOnHand,
		QuantityDamaged:  data.QuantityDamaged,
		QuantityReserved: data.QuantityReserved,
		QuantityIncoming: data.QuantityIncoming,
		QuantityOutgoing: data.QuantityOutgoing,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateResultPayload{
			Inventory: result,
		},
	}, nil
}

func (s *Service) Delete(data DeleteData) (*DeleteResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.Inventory().Exists(queries.InventoryParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.Inventory().Delete(queries.InventoryParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetOne(data GetOneData) (*GetOneResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetOneResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.Inventory().FindOne(queries.InventoryParams{
		ID:            data.ID,
		BinID:         data.BinID,
		LocationID:    &locationID,
		Related:       data.Related,
		ToUpdatedAt:   data.ToUpdatedAt,
		FromUpdatedAt: data.FromUpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetOneResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetOneResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetOneResultPayload{
			Inventory: result,
		},
	}, nil
}

func (s *Service) GetMany(data GetManyData) (*GetManyResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetManyResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.Inventory().FindMany(
		queries.InventoryParams{
			BinID:         data.BinID,
			LocationID:    &locationID,
			Related:       data.Related,
			ToUpdatedAt:   data.ToUpdatedAt,
			FromUpdatedAt: data.FromUpdatedAt,
		},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetManyResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetManyResultPayload{
			Inventory: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) GetSummary(data GetSummaryData) (*GetSummaryResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetSummaryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.Inventory().GenerateSummary(queries.InventoryParams{
		ID:            data.ID,
		BinID:         data.BinID,
		LocationID:    &locationID,
		ToUpdatedAt:   data.ToUpdatedAt,
		FromUpdatedAt: data.FromUpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	return &GetSummaryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetSummaryResultPayload{
			Summary: result,
		},
	}, nil
}

func (s *Service) AddMovement(data AddMovementData) (*AddMovementResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddMovementResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddMovementResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.Inventory().Exists(queries.InventoryParams{
		ID:         &data.InventoryID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddMovementResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.InventoryMovements().Create(entities.InventoryMovement{
		InventoryID:  data.InventoryID,
		ReferenceID:  data.ReferenceID,
		MovementType: data.MovementType,
		Quantity:     data.Quantity,
		Reason:       data.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &AddMovementResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddMovementResultPayload{
			Movement: result,
		},
	}, nil
}

func (s *Service) UpdateMovement(data UpdateMovementData) (*UpdateMovementResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateMovementResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateMovementResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.InventoryMovements().Exists(queries.InventoryMovementsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateMovementResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.InventoryMovements().Update(queries.InventoryMovementsParams{
		ID: &data.ID,
	}, entities.InventoryMovement{
		InventoryID:  data.InventoryID,
		ReferenceID:  data.ReferenceID,
		MovementType: data.MovementType,
		Quantity:     data.Quantity,
		Reason:       data.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateMovementResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateMovementResultPayload{
			Movement: result,
		},
	}, nil
}

func (s *Service) DeleteMovement(data DeleteMovementData) (*DeleteMovementResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteMovementResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteMovementResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.InventoryMovements().Exists(queries.InventoryMovementsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteMovementResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.InventoryMovements().Delete(queries.InventoryMovementsParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteMovementResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteMovementResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetMovement(data GetMovementData) (*GetMovementResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetMovementResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.InventoryMovements().FindOne(queries.InventoryMovementsParams{
		ID:           data.ID,
		InventoryID:  data.InventoryID,
		ReferenceID:  data.ReferenceID,
		LocationID:   &locationID,
		MovementType: data.MovementType,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetMovementResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetMovementResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetMovementResultPayload{
			Movement: result,
		},
	}, nil
}

func (s *Service) GetMovements(data GetMovementsData) (*GetMovementsResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetMovementsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.InventoryMovements().FindMany(
		queries.InventoryMovementsParams{
			InventoryID:  data.InventoryID,
			ReferenceID:  data.ReferenceID,
			LocationID:   &locationID,
			MovementType: data.MovementType,
		},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetMovementsResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetMovementsResultPayload{
			Movements: results,
		},
		Pagination: paging,
	}, nil
}
