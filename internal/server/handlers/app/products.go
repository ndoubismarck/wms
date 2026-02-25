package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/products"

	"github.com/gin-gonic/gin"
)

type (
	addProductRequestBody  products.AddData
	addProductResponseBody struct {
		apiResponseBody
		Data products.AddResultPayload `json:"data"`
	}
	updateProductRequestBody  products.UpdateData
	updateProductResponseBody struct {
		apiResponseBody
		Data products.UpdateResultPayload `json:"data"`
	}
	deleteProductResponseBody struct {
		apiResponseBody
	}
	getProductResponseBody struct {
		apiResponseBody
		Data products.GetOneResultPayload `json:"data"`
	}
	getProductsRequestBody  products.GetManyData
	getProductsResponseBody struct {
		apiResponseBody
		Data products.GetManyResultPayload `json:"data"`
	}
)

type (
	addProductCategoryRequestBody  products.AddCategoryData
	addProductCategoryResponseBody struct {
		apiResponseBody
		Data products.AddCategoryResultPayload `json:"data"`
	}
	updateProductCategoryRequestBody  products.UpdateCategoryData
	updateProductCategoryResponseBody struct {
		apiResponseBody
		Data products.UpdateCategoryResultPayload `json:"data"`
	}
	deleteProductCategoryResponseBody struct {
		apiResponseBody
	}
	getProductCategoryResponseBody struct {
		apiResponseBody
		Data products.GetCategoryResultPayload `json:"data"`
	}
	getProductCategoriesRequestBody  products.GetCategoriesData
	getProductCategoriesResponseBody struct {
		apiResponseBody
		Data products.GetCategoriesResultPayload `json:"data"`
	}
)

type (
	addProductSubcategoryRequestBody  products.AddSubcategoryData
	addProductSubcategoryResponseBody struct {
		apiResponseBody
		Data products.AddSubcategoryResultPayload `json:"data"`
	}
	updateProductSubcategoryRequestBody  products.UpdateSubcategoryData
	updateProductSubcategoryResponseBody struct {
		apiResponseBody
		Data products.UpdateSubcategoryResultPayload `json:"data"`
	}
	deleteProductSubcategoryResponseBody struct {
		apiResponseBody
	}
	getProductSubcategoryResponseBody struct {
		apiResponseBody
		Data products.GetSubcategoryResultPayload `json:"data"`
	}
	getProductSubcategoriesRequestBody  products.GetSubcategoriesData
	getProductSubcategoriesResponseBody struct {
		apiResponseBody
		Data products.GetSubcategoriesResultPayload `json:"data"`
	}
)

type (
	addProductBrandRequestBody  products.AddBrandData
	addProductBrandResponseBody struct {
		apiResponseBody
		Data products.AddBrandResultPayload `json:"data"`
	}
	updateProductBrandRequestBody  products.UpdateBrandData
	updateProductBrandResponseBody struct {
		apiResponseBody
		Data products.UpdateBrandResultPayload `json:"data"`
	}
	deleteProductBrandResponseBody struct {
		apiResponseBody
	}
	getProductBrandResponseBody struct {
		apiResponseBody
		Data products.GetBrandResultPayload `json:"data"`
	}
	getProductBrandsRequestBody  products.GetBrandsData
	getProductBrandsResponseBody struct {
		apiResponseBody
		Data products.GetBrandsResultPayload `json:"data"`
	}
)

type (
	addProductAttributeRequestBody  products.AddAttributeData
	addProductAttributeResponseBody struct {
		apiResponseBody
		Data products.AddAttributeResultPayload `json:"data"`
	}
	updateProductAttributeRequestBody  products.UpdateAttributeData
	updateProductAttributeResponseBody struct {
		apiResponseBody
		Data products.UpdateAttributeResultPayload `json:"data"`
	}
	deleteProductAttributeResponseBody struct {
		apiResponseBody
	}
	getProductAttributeResponseBody struct {
		apiResponseBody
		Data products.GetAttributeResultPayload `json:"data"`
	}
	getProductAttributesRequestBody  products.GetAttributesData
	getProductAttributesResponseBody struct {
		apiResponseBody
		Data products.GetAttributesResultPayload `json:"data"`
	}
)

func (h *handler) addProduct(context *gin.Context) {
	var (
		requestBody  addProductRequestBody
		responseBody addProductResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().Add(products.AddData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateProduct(context *gin.Context) {
	var (
		requestBody  updateProductRequestBody
		responseBody updateProductResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ID = context.Param("id")
	result, err := h.services.Products().Update(products.UpdateData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) deleteProduct(context *gin.Context) {
	var responseBody deleteProductResponseBody
	result, err := h.services.Products().Delete(products.DeleteData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProduct(context *gin.Context) {
	var responseBody getProductResponseBody
	id := context.Param("id")
	result, err := h.services.Products().GetOne(products.GetOneData{
		LocationID: context.Param("lid"),
		ID:         &id,
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProducts(context *gin.Context) {
	var (
		requestBody  getProductsRequestBody
		responseBody getProductsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().GetMany(products.GetManyData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) addProductCategory(context *gin.Context) {
	var (
		requestBody  addProductCategoryRequestBody
		responseBody addProductCategoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().AddCategory(products.AddCategoryData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateProductCategory(context *gin.Context) {
	var (
		requestBody  updateProductCategoryRequestBody
		responseBody updateProductCategoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ID = context.Param("id")
	result, err := h.services.Products().UpdateCategory(products.UpdateCategoryData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) deleteProductCategory(context *gin.Context) {
	var responseBody deleteProductCategoryResponseBody
	result, err := h.services.Products().DeleteCategory(products.DeleteCategoryData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductCategory(context *gin.Context) {
	var responseBody getProductCategoryResponseBody
	id := context.Param("id")
	result, err := h.services.Products().GetCategory(products.GetCategoryData{
		LocationID: context.Param("lid"),
		ID:         &id,
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductCategories(context *gin.Context) {
	var (
		requestBody  getProductCategoriesRequestBody
		responseBody getProductCategoriesResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().GetCategories(products.GetCategoriesData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) addProductSubcategory(context *gin.Context) {
	var (
		requestBody  addProductSubcategoryRequestBody
		responseBody addProductSubcategoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().AddSubcategory(products.AddSubcategoryData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateProductSubcategory(context *gin.Context) {
	var (
		requestBody  updateProductSubcategoryRequestBody
		responseBody updateProductSubcategoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ID = context.Param("id")
	result, err := h.services.Products().UpdateSubcategory(products.UpdateSubcategoryData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) deleteProductSubcategory(context *gin.Context) {
	var responseBody deleteProductSubcategoryResponseBody
	result, err := h.services.Products().DeleteSubcategory(products.DeleteSubcategoryData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductSubcategory(context *gin.Context) {
	var responseBody getProductSubcategoryResponseBody
	id := context.Param("id")
	result, err := h.services.Products().GetSubcategory(products.GetSubcategoryData{
		LocationID: context.Param("lid"),
		ID:         &id,
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductSubcategories(context *gin.Context) {
	var (
		requestBody  getProductSubcategoriesRequestBody
		responseBody getProductSubcategoriesResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().GetSubcategories(products.GetSubcategoriesData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) addProductBrand(context *gin.Context) {
	var (
		requestBody  addProductBrandRequestBody
		responseBody addProductBrandResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().AddBrand(products.AddBrandData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateProductBrand(context *gin.Context) {
	var (
		requestBody  updateProductBrandRequestBody
		responseBody updateProductBrandResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ID = context.Param("id")
	result, err := h.services.Products().UpdateBrand(products.UpdateBrandData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) deleteProductBrand(context *gin.Context) {
	var responseBody deleteProductBrandResponseBody
	result, err := h.services.Products().DeleteBrand(products.DeleteBrandData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductBrand(context *gin.Context) {
	var responseBody getProductBrandResponseBody
	id := context.Param("id")
	result, err := h.services.Products().GetBrand(products.GetBrandData{
		LocationID: context.Param("lid"),
		ID:         &id,
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductBrands(context *gin.Context) {
	var (
		requestBody  getProductBrandsRequestBody
		responseBody getProductBrandsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().GetBrands(products.GetBrandsData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) addProductAttribute(context *gin.Context) {
	var (
		requestBody  addProductAttributeRequestBody
		responseBody addProductAttributeResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().AddAttribute(products.AddAttributeData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateProductAttribute(context *gin.Context) {
	var (
		requestBody  updateProductAttributeRequestBody
		responseBody updateProductAttributeResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ID = context.Param("id")
	result, err := h.services.Products().UpdateAttribute(products.UpdateAttributeData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) deleteProductAttribute(context *gin.Context) {
	var responseBody deleteProductAttributeResponseBody
	result, err := h.services.Products().DeleteAttribute(products.DeleteAttributeData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductAttribute(context *gin.Context) {
	var responseBody getProductAttributeResponseBody
	id := context.Param("id")
	result, err := h.services.Products().GetAttribute(products.GetAttributeData{
		LocationID: context.Param("lid"),
		ID:         &id,
	})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductAttributes(context *gin.Context) {
	var (
		requestBody  getProductAttributesRequestBody
		responseBody getProductAttributesResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Products().GetAttributes(products.GetAttributesData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) getProductBarcode(context *gin.Context) {
	id := context.Param("id")
	result, err := h.services.Products().GetBarcode(products.GetBarcodeData{
		LocationID: context.Param("lid"),
		ID:         id,
	})
	if err != nil {
		context.Status(http.StatusInternalServerError)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		if result.Code == types.ServiceResultCodeNotFound {
			context.Status(http.StatusNotFound)
			return
		}
		context.Status(http.StatusForbidden)
		return
	}
	context.File(result.Payload.File.Path)
}
