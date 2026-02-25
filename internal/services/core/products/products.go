package products

import (
	"encoding/json"
	"os"
	"path/filepath"
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/fileutil"
	"server/internal/core/shared/utils/imageutil"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database/queries"
	"strings"
	"time"

	"gorm.io/datatypes"
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

func parseAttributeDate(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		result, err := time.Parse(layout, value)
		if err == nil {
			return result, true
		}
	}

	return time.Time{}, false
}

func marshalAttributeJSON(value any) (*datatypes.JSON, bool) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, false
	}
	result := datatypes.JSON(payload)
	return &result, true
}

func (s *Service) getAttributes(data map[string]attribute) []entities.ProductAttribute {
	attributes := make([]entities.ProductAttribute, 0, len(data))
	for key, val := range data {
		if key == "" {
			continue
		}

		attr := entities.ProductAttribute{
			ProductAttributeDefinitionID: key,
		}

		if len(val.Options) > 0 {
			options := make([]string, 0, len(val.Options))
			for _, item := range val.Options {
				item = strings.TrimSpace(item)
				if item != "" {
					options = append(options, item)
				}
			}
			if len(options) == 0 {
				continue
			}

			payload, ok := marshalAttributeJSON(options)
			if !ok {
				s.ctx.Logger().Warnf("failed to marshal attribute options for %s", key)
				continue
			}

			attr.ValueType = entities.ProductAttributeTypeOptions
			attr.ValueJSON = payload
			firstOptionID := options[0]
			attr.ProductAttributeOptionID = &firstOptionID
			attr.Value = options
			attributes = append(attributes, attr)
			continue
		}

		if val.Option != nil {
			optionID := strings.TrimSpace(*val.Option)
			if optionID == "" {
				continue
			}

			attr.ValueType = entities.ProductAttributeTypeOption
			attr.ProductAttributeOptionID = &optionID
			attr.Value = optionID
			attributes = append(attributes, attr)
			continue
		}

		if val.String != nil {
			attr.ValueType = entities.ProductAttributeTypeString
			attr.ValueString = val.String
			attr.Value = *val.String
			attributes = append(attributes, attr)
			continue
		}

		if val.Int != nil {
			attr.ValueType = entities.ProductAttributeTypeInt
			attr.ValueInt = val.Int
			attr.Value = *val.Int
			attributes = append(attributes, attr)
			continue
		}

		if val.Bool != nil {
			attr.ValueType = entities.ProductAttributeTypeBool
			attr.ValueBool = val.Bool
			attr.Value = *val.Bool
			attributes = append(attributes, attr)
			continue
		}

		if val.Date != nil {
			dateValue, ok := parseAttributeDate(*val.Date)
			if !ok {
				s.ctx.Logger().Warnf("invalid date value for attribute %s: %q", key, *val.Date)
				continue
			}

			attr.ValueType = entities.ProductAttributeTypeDate
			attr.ValueDate = &dateValue
			attr.Value = dateValue
			attributes = append(attributes, attr)
			continue
		}

		if val.Range != nil {
			rangeValue := datatypes.NewJSONType(*val.Range)
			attr.ValueType = entities.ProductAttributeTypeJSON
			attr.ValueRange = &rangeValue
			attr.Value = rangeValue.Data()
			attributes = append(attributes, attr)
			continue
		}

		if val.JSON != nil {
			payload, ok := marshalAttributeJSON(val.JSON)
			if !ok {
				s.ctx.Logger().Warnf("failed to marshal json value for attribute %s", key)
				continue
			}

			attr.ValueType = entities.ProductAttributeTypeJSON
			attr.ValueJSON = payload
			attr.Value = val.JSON
			attributes = append(attributes, attr)
			continue
		}

		if val.Float != nil {
			attr.ValueType = entities.ProductAttributeTypeFloat
			attr.ValueFloat = val.Float
			attr.Value = *val.Float
			attributes = append(attributes, attr)
		}
	}

	return attributes
}

func (s *Service) deleteUploadedImages(images []string) {
	if err := s.providers.Files().DeleteFiles(images); err != nil {
		s.ctx.Logger().Warn(err)
	}
}

func (s *Service) processUploadedImages(paths []string) ([]entities.ProductVariantMedia, error) {
	paths, err := s.providers.Files().CopyUploadedFiles(paths, "www/images")
	if err != nil {
		return nil, err
	}
	var images []entities.ProductVariantMedia
	for _, val := range paths {
		images = append(images, entities.ProductVariantMedia{
			Path: val,
		})
	}
	return images, nil
}

func (s *Service) deleteProcessedImages(images []entities.ProductVariantMedia) {
	var paths []string
	for _, val := range images {
		paths = append(paths, val.Path)
	}
	if err := s.providers.Files().DeleteFiles(paths); err != nil {
		s.ctx.Logger().Warn(err)
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
	images, err := s.processUploadedImages(data.Images)
	if err != nil {
		return nil, err
	}
	var variants []entities.ProductVariant
	if len(images) > 0 {
		variants = append(variants, entities.ProductVariant{
			Media: images,
		})
	}
	product, err := query.Products().Create(entities.Product{
		LocationD:     data.LocationID,
		SubcategoryID: data.SubcategoryID,
		BrandID:       data.BrandID,
		Name:          data.Name,
		Status:        entities.ProductStatusAvailable,
		Attributes:    s.getAttributes(data.Attributes),
		Description:   data.Description,
		Variants:      variants,
	})
	if err != nil {
		s.deleteProcessedImages(images)
		return nil, err
	}
	s.deleteUploadedImages(data.Images)
	return &AddResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddResultPayload{
			Product: product,
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
	exists, err := query.Products().Exists(queries.ProductParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	product, exists, err := query.Products().FindOne(queries.ProductParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	var update bool
	updated := entities.Product{}
	if data.SubcategoryID != nil {
		update = true
		updated.SubcategoryID = *data.SubcategoryID
	}
	if data.BrandID != nil {
		update = true
		updated.BrandID = *data.BrandID
	}
	if data.Name != nil {
		update = true
		updated.Name = *data.Name
	}
	if data.Description != nil {
		update = true
		updated.Description = *data.Description
	}
	if data.Status != nil {
		update = true
		updated.Status = *data.Status
	}
	if update {
		product, _, err = query.Products().Update(queries.ProductParams{
			ID:         &data.ID,
			LocationID: &data.LocationID,
		}, updated)
		if err != nil {
			return nil, err
		}
	}
	return &UpdateResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateResultPayload{
			Product: product,
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
	result, err := query.Products().Delete(queries.ProductParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
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
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetOneResult{Code: types.ServiceResultCodeFailed}, nil
	}
	product, exists, err := query.Products().FindOne(queries.ProductParams{
		ID:         data.ID,
		LocationID: &data.LocationID,
		BinID:      data.BinID,
		Related:    true,
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
			Product: product,
		},
	}, nil
}

func (s *Service) GetMany(data GetManyData) (*GetManyResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetManyResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.Products().FindMany(
		queries.ProductParams{
			LocationID: &data.LocationID,
			BinID:      data.BinID,
			Related:    true,
		},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}))
	if err != nil {
		return nil, err
	}
	return &GetManyResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetManyResultPayload{
			Products: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddCategory(data AddCategoryData) (*AddCategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddCategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddCategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.ProductCategories().Create(entities.ProductCategory{
		BrandID: data.BrandID,
		Name:    data.Name,
	})
	if err != nil {
		return nil, err
	}
	return &AddCategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddCategoryResultPayload{
			Category: result,
		},
	}, nil
}

func (s *Service) UpdateCategory(data UpdateCategoryData) (*UpdateCategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateCategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateCategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductCategories().FindOne(queries.ProductCategoriesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateCategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	var update bool
	var updated entities.ProductCategory
	if data.Name != nil {
		update = true
		updated.Name = *data.Name
	}
	if data.BrandID != nil {
		update = true
		updated.BrandID = *data.BrandID
	}
	if update {
		result, exists, err = query.ProductCategories().Update(queries.ProductCategoriesParams{
			ID: &data.ID,
		}, updated)
		if err != nil {
			return nil, err
		}
		if !exists {
			return &UpdateCategoryResult{Code: types.ServiceResultCodeNotFound}, nil
		}
	}

	return &UpdateCategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateCategoryResultPayload{
			Category: result,
		},
	}, nil
}

func (s *Service) DeleteCategory(data DeleteCategoryData) (*DeleteCategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteCategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteCategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.ProductCategories().Delete(queries.ProductCategoriesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteCategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteCategoryResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetCategory(data GetCategoryData) (*GetCategoryResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetCategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductCategories().FindOne(queries.ProductCategoriesParams{
		ID: data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetCategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetCategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetCategoryResultPayload{
			Category: result,
		},
	}, nil
}

func (s *Service) GetCategories(data GetCategoriesData) (*GetCategoriesResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetCategoriesResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.ProductCategories().FindMany(
		queries.ProductCategoriesParams{BrandID: data.BrandID},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetCategoriesResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetCategoriesResultPayload{
			Categories: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddSubcategory(data AddSubcategoryData) (*AddSubcategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddSubcategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddSubcategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	categoryExists, err := query.ProductCategories().Exists(queries.ProductCategoriesParams{
		ID: &data.CategoryID,
	})
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return &AddSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.ProductSubcategories().Create(entities.ProductSubcategory{
		CategoryID: data.CategoryID,
		Name:       data.Name,
	})
	if err != nil {
		return nil, err
	}
	return &AddSubcategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddSubcategoryResultPayload{
			Subcategory: result,
		},
	}, nil
}

func (s *Service) UpdateSubcategory(data UpdateSubcategoryData) (*UpdateSubcategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateSubcategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateSubcategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	categoryExists, err := query.ProductCategories().Exists(queries.ProductCategoriesParams{
		ID: &data.CategoryID,
	})
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return &UpdateSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, exists, err := query.ProductSubcategories().Update(queries.ProductSubcategoriesParams{
		ID: &data.ID,
	}, entities.ProductSubcategory{
		CategoryID: data.CategoryID,
		Name:       data.Name,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateSubcategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateSubcategoryResultPayload{
			Subcategory: result,
		},
	}, nil
}

func (s *Service) DeleteSubcategory(data DeleteSubcategoryData) (*DeleteSubcategoryResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteSubcategoryResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteSubcategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	subcategory, exists, err := query.ProductSubcategories().FindOne(queries.ProductSubcategoriesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	categoryExists, err := query.ProductCategories().Exists(queries.ProductCategoriesParams{
		ID: &subcategory.CategoryID,
	})
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return &DeleteSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err := query.ProductSubcategories().Delete(queries.ProductSubcategoriesParams{
		ID: &data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteSubcategoryResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetSubcategory(data GetSubcategoryData) (*GetSubcategoryResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetSubcategoryResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductSubcategories().FindOne(queries.ProductSubcategoriesParams{
		ID: data.ID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	categoryExists, err := query.ProductCategories().Exists(queries.ProductCategoriesParams{
		ID: &result.CategoryID,
	})
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return &GetSubcategoryResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetSubcategoryResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetSubcategoryResultPayload{
			Subcategory: result,
		},
	}, nil
}

func (s *Service) GetSubcategories(data GetSubcategoriesData) (*GetSubcategoriesResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetSubcategoriesResult{Code: types.ServiceResultCodeFailed}, nil
	}
	params := queries.ProductSubcategoriesParams{
		CategoryID: data.CategoryID,
	}
	if data.CategoryID != nil {
		categoryExists, err := query.ProductCategories().Exists(queries.ProductCategoriesParams{
			ID: data.CategoryID,
		})
		if err != nil {
			return nil, err
		}
		if !categoryExists {
			return &GetSubcategoriesResult{Code: types.ServiceResultCodeNotFound}, nil
		}
	}
	results, paging, err := query.ProductSubcategories().FindMany(
		params,
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetSubcategoriesResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetSubcategoriesResultPayload{
			Subcategories: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddBrand(data AddBrandData) (*AddBrandResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddBrandResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddBrandResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.ProductBrands().Create(entities.ProductBrand{
		LocationD: data.LocationID,
		Name:      data.Name,
	})
	if err != nil {
		return nil, err
	}
	return &AddBrandResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddBrandResultPayload{
			Brand: result,
		},
	}, nil
}

func (s *Service) UpdateBrand(data UpdateBrandData) (*UpdateBrandResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateBrandResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateBrandResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductBrands().Update(queries.ProductBrandsParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	}, entities.ProductBrand{
		Name: data.Name,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateBrandResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateBrandResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateBrandResultPayload{
			Brand: result,
		},
	}, nil
}

func (s *Service) DeleteBrand(data DeleteBrandData) (*DeleteBrandResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteBrandResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteBrandResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.ProductBrands().Delete(queries.ProductBrandsParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteBrandResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteBrandResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetBrand(data GetBrandData) (*GetBrandResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBrandResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductBrands().FindOne(queries.ProductBrandsParams{
		ID:         data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetBrandResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetBrandResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBrandResultPayload{
			Brand: result,
		},
	}, nil
}

func (s *Service) GetBrands(data GetBrandsData) (*GetBrandsResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBrandsResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.ProductBrands().FindMany(
		queries.ProductBrandsParams{LocationID: &data.LocationID},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetBrandsResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBrandsResultPayload{
			Brands: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) AddAttribute(data AddAttributeData) (*AddAttributeResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddAttributeResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddAttributeResult{Code: types.ServiceResultCodeFailed}, nil
	}
	appliesTo := data.AppliesTo
	if appliesTo == "" {
		appliesTo = entities.ProductAttributeAppliesToProduct
	}
	dataType := data.DataType
	if dataType == "" {
		dataType = entities.ProductAttributeTypeString
	}
	fieldType := data.FieldType
	if fieldType == "" {
		fieldType = entities.DisplayText
	}
	result, err := query.ProductAttributeDefinitions().Create(entities.ProductAttributeDefinition{
		LocationD:       data.LocationID,
		IsSystem:        data.IsSystem,
		Label:           data.Label,
		AppliesTo:       appliesTo,
		DataType:        dataType,
		FieldType:       fieldType,
		ValidationRules: data.ValidationRules,
	})
	if err != nil {
		return nil, err
	}
	return &AddAttributeResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddAttributeResultPayload{
			Attribute: result,
		},
	}, nil
}

func (s *Service) UpdateAttribute(data UpdateAttributeData) (*UpdateAttributeResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateAttributeResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateAttributeResult{Code: types.ServiceResultCodeFailed}, nil
	}
	current, exists, err := query.ProductAttributeDefinitions().FindOne(queries.ProductAttributeDefinitionsParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateAttributeResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	updated := entities.ProductAttributeDefinition{
		IsSystem:        data.IsSystem,
		Label:           data.Label,
		AppliesTo:       data.AppliesTo,
		DataType:        data.DataType,
		FieldType:       data.FieldType,
		ValidationRules: data.ValidationRules,
	}
	if updated.Label == "" {
		updated.Label = current.Label
	}
	if updated.AppliesTo == "" {
		updated.AppliesTo = current.AppliesTo
	}
	if updated.DataType == "" {
		updated.DataType = current.DataType
	}
	if updated.FieldType == "" {
		updated.FieldType = current.FieldType
	}
	result, exists, err := query.ProductAttributeDefinitions().Update(queries.ProductAttributeDefinitionsParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	}, updated)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateAttributeResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateAttributeResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateAttributeResultPayload{
			Attribute: result,
		},
	}, nil
}

func (s *Service) DeleteAttribute(data DeleteAttributeData) (*DeleteAttributeResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteAttributeResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteAttributeResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, err := query.ProductAttributeDefinitions().Delete(queries.ProductAttributeDefinitionsParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !result {
		return &DeleteAttributeResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteAttributeResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetAttribute(data GetAttributeData) (*GetAttributeResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetAttributeResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.ProductAttributeDefinitions().FindOne(queries.ProductAttributeDefinitionsParams{
		ID:         data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetAttributeResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &GetAttributeResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetAttributeResultPayload{
			Attribute: result,
		},
	}, nil
}

func (s *Service) GetAttributes(data GetAttributesData) (*GetAttributesResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetAttributesResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.ProductAttributeDefinitions().FindMany(
		queries.ProductAttributeDefinitionsParams{LocationID: &data.LocationID},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetAttributesResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetAttributesResultPayload{
			Attributes: results,
		},
		Pagination: paging,
	}, nil
}

func (s *Service) GetBarcode(data GetBarcodeData) (*GetBarcodeResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetBarcodeResult{Code: types.ServiceResultCodeFailed}, nil
	}
	product, exists, err := query.Products().FindOne(queries.ProductParams{
		ID:         &data.ID,
		LocationID: &data.LocationID,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetBarcodeResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	dir := "tmp/barcodes"
	if !fileutil.Exists(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	file := cryptoutil.Sha1(product.ID)
	file = filepath.Join(dir, file+".png")
	if fileutil.Exists(file) {
		return &GetBarcodeResult{
			Code: types.ServiceResultCodeSuccess,
			Payload: GetBarcodeResultPayload{
				File: GetBarcodeResultPayloadFile{
					Path: file,
				},
			},
		}, nil
	}
	barcode, err := imageutil.GenerateBarcode(product.ID)
	if err != nil {
		return nil, err
	}
	if err = os.WriteFile(file, barcode, os.ModePerm); err != nil {
		return nil, err
	}
	return &GetBarcodeResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetBarcodeResultPayload{
			File: GetBarcodeResultPayloadFile{
				Path: file,
			},
		},
	}, nil
}
