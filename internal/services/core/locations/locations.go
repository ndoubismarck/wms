package locations

import (
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database/queries"
)

type Service struct {
	ctx types.IContext
	tmp struct {
		admin *entities.User
	}
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func (s *Service) GetLocations(data GetLocationsData) (*GetLocationsResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetLocationsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.Locations().FindMany(
		queries.LocationsParams{},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetLocationsResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetLocationsResultPayload{
			Locations: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddAisle(data AddAisleData) (*AddAisleResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddAisleResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddAisleResult{Code: types.ServiceResultCodeFailed}, nil
	}
	parentExists, err := query.Locations().Exists(queries.LocationsParams{
		ID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddAisleResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationAisles().Create(entities.LocationAisle{
		LocationID:  data.LocationID,
		Code:        data.Code,
		Name:        data.Name,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &AddAisleResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddAisleResultPayload{
			Aisle: result,
		},
	}, nil
}

func (s *Service) UpdateAisle(data UpdateAisleData) (*UpdateAisleResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateAisleResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateAisleResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationAisles().Exists(queries.LocationAislesParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateAisleResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.LocationAisles().Update(queries.LocationAislesParams{
		ID: &data.ID,
	}, entities.LocationAisle{
		Code:        data.Code,
		Name:        data.Name,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateAisleResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateAisleResultPayload{
			Aisle: result,
		},
	}, nil
}

func (s *Service) DeleteAisle(data DeleteAisleData) (*DeleteAisleResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteAisleResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteAisleResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationAisles().Exists(queries.LocationAislesParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteAisleResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationAisles().Delete(queries.LocationAislesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteAisleResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteAisleResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetAisle(data GetAisleData) (*GetAisleResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetAisleResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.LocationAisles().FindOne(queries.LocationAislesParams{
		ID:         data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetAisleResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetAisleResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetAisleResultPayload{
			Aisle: result,
		},
	}, nil
}

func (s *Service) GetAisles(data GetAislesData) (*GetAislesResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetAislesResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.LocationAisles().FindMany(
		queries.LocationAislesParams{
			LocationID: &locationID,
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
	return &GetAislesResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetAislesResultPayload{
			Aisles: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddBay(data AddBayData) (*AddBayResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddBayResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddBayResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.LocationAisles().Exists(queries.LocationAislesParams{
		ID:         &data.AisleID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddBayResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationBays().Create(entities.LocationBay{
		Code:        data.Code,
		Name:        data.Name,
		AisleID:     data.AisleID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &AddBayResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddBayResultPayload{
			Bay: result,
		},
	}, nil
}

func (s *Service) UpdateBay(data UpdateBayData) (*UpdateBayResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateBayResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateBayResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationBays().Exists(queries.LocationBaysParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateBayResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.LocationBays().Update(queries.LocationBaysParams{
		ID: &data.ID,
	}, entities.LocationBay{
		Code:        data.Code,
		Name:        data.Name,
		AisleID:     data.AisleID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateBayResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateBayResultPayload{
			Bay: result,
		},
	}, nil
}

func (s *Service) DeleteBay(data DeleteBayData) (*DeleteBayResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteBayResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteBayResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationBays().Exists(queries.LocationBaysParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteBayResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationBays().Delete(queries.LocationBaysParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteBayResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteBayResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetBay(data GetBayData) (*GetBayResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBayResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.LocationBays().FindOne(queries.LocationBaysParams{
		ID:         data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetBayResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetBayResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBayResultPayload{
			Bay: result,
		},
	}, nil
}

func (s *Service) GetBays(data GetBaysData) (*GetBaysResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBaysResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.LocationBays().FindMany(
		queries.LocationBaysParams{
			AisleID:    data.AisleID,
			LocationID: &locationID,
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
	return &GetBaysResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBaysResultPayload{
			Bays: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddShelf(data AddShelfData) (*AddShelfResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddShelfResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddShelfResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.LocationBays().Exists(queries.LocationBaysParams{
		ID:         &data.BayID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddShelfResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationShelves().Create(entities.LocationShelf{
		Code:        data.Code,
		Name:        data.Name,
		BayID:       data.BayID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &AddShelfResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddShelfResultPayload{
			Shelf: result,
		},
	}, nil
}

func (s *Service) UpdateShelf(data UpdateShelfData) (*UpdateShelfResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateShelfResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateShelfResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationShelves().Exists(queries.LocationShelvesParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateShelfResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.LocationShelves().Update(queries.LocationShelvesParams{
		ID: &data.ID,
	}, entities.LocationShelf{
		Code:        data.Code,
		Name:        data.Name,
		BayID:       data.BayID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateShelfResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateShelfResultPayload{
			Shelf: result,
		},
	}, nil
}

func (s *Service) DeleteShelf(data DeleteShelfData) (*DeleteShelfResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteShelfResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteShelfResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationShelves().Exists(queries.LocationShelvesParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteShelfResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationShelves().Delete(queries.LocationShelvesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteShelfResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteShelfResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetShelf(data GetShelfData) (*GetShelfResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetShelfResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.LocationShelves().FindOne(queries.LocationShelvesParams{
		ID:         data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetShelfResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetShelfResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetShelfResultPayload{
			Shelf: result,
		},
	}, nil
}

func (s *Service) GetShelves(data GetShelvesData) (*GetShelvesResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetShelvesResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.LocationShelves().FindMany(
		queries.LocationShelvesParams{
			BayID:      data.BayID,
			LocationID: &locationID,
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
	return &GetShelvesResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetShelvesResultPayload{
			Shelves: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddShelfLevel(data AddShelfLevelData) (*AddShelfLevelResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddShelfLevelResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddShelfLevelResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.LocationShelves().Exists(queries.LocationShelvesParams{
		ID:         &data.ShelfID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddShelfLevelResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationShelfLevels().Create(entities.LocationShelfLevel{
		Code:        data.Code,
		Name:        data.Name,
		ShelfID:     data.ShelfID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &AddShelfLevelResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddShelfLevelResultPayload{
			ShelfLevel: result,
		},
	}, nil
}

func (s *Service) UpdateShelfLevel(data UpdateShelfLevelData) (*UpdateShelfLevelResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateShelfLevelResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateShelfLevelResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationShelfLevels().Exists(queries.LocationShelfLevelsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateShelfLevelResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.LocationShelfLevels().Update(queries.LocationShelfLevelsParams{
		ID: &data.ID,
	}, entities.LocationShelfLevel{
		Code:        data.Code,
		Name:        data.Name,
		ShelfID:     data.ShelfID,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateShelfLevelResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateShelfLevelResultPayload{
			ShelfLevel: result,
		},
	}, nil
}

func (s *Service) DeleteShelfLevel(data DeleteShelfLevelData) (*DeleteShelfLevelResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteShelfLevelResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteShelfLevelResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationShelfLevels().Exists(queries.LocationShelfLevelsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteShelfLevelResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationShelfLevels().Delete(queries.LocationShelfLevelsParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteShelfLevelResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteShelfLevelResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetShelfLevel(data GetShelfLevelData) (*GetShelfLevelResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetShelfLevelResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.LocationShelfLevels().FindOne(queries.LocationShelfLevelsParams{
		ID:         data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetShelfLevelResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetShelfLevelResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetShelfLevelResultPayload{
			ShelfLevel: result,
		},
	}, nil
}

func (s *Service) GetShelfLevels(data GetShelfLevelsData) (*GetShelfLevelsResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetShelfLevelsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.LocationShelfLevels().FindMany(
		queries.LocationShelfLevelsParams{
			ShelfID:    data.ShelfID,
			LocationID: &locationID,
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
	return &GetShelfLevelsResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetShelfLevelsResultPayload{
			ShelfLevels: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddBin(data AddBinData) (*AddBinResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddBinResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddBinResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	parentExists, err := query.LocationShelfLevels().Exists(queries.LocationShelfLevelsParams{
		ID:         &data.ShelfLevelID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return &AddBinResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationBins().Create(entities.LocationBin{
		Code:         data.Code,
		Name:         data.Name,
		ShelfLevelID: data.ShelfLevelID,
		Description:  data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &AddBinResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddBinResultPayload{
			Bin: result,
		},
	}, nil
}

func (s *Service) UpdateBin(data UpdateBinData) (*UpdateBinResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateBinResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateBinResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationBins().Exists(queries.LocationBinsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateBinResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, _, err := query.LocationBins().Update(queries.LocationBinsParams{
		ID: &data.ID,
	}, entities.LocationBin{
		Code:         data.Code,
		Name:         data.Name,
		ShelfLevelID: data.ShelfLevelID,
		Description:  data.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateBinResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateBinResultPayload{
			Bin: result,
		},
	}, nil
}

func (s *Service) DeleteBin(data DeleteBinData) (*DeleteBinResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteBinResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteBinResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.LocationBins().Exists(queries.LocationBinsParams{
		ID:         &data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteBinResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.LocationBins().Delete(queries.LocationBinsParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteBinResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteBinResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetBin(data GetBinData) (*GetBinResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBinResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.LocationBins().FindOne(queries.LocationBinsParams{
		ID:         data.ID,
		LocationID: &locationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetBinResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetBinResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBinResultPayload{
			Bin: result,
		},
	}, nil
}

func (s *Service) GetBins(data GetBinsData) (*GetBinsResult, error) {
	locationID := data.LocationID
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBinsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.LocationBins().FindMany(
		queries.LocationBinsParams{
			ShelfLevelID: data.ShelfLevelID,
			LocationID:   &locationID,
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
	return &GetBinsResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBinsResultPayload{
			Bins: results,
		},
		Pagination: paging,
	}, nil
}
