package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/locations"

	"github.com/gin-gonic/gin"
)

type (
	getLocationsRequestBody  locations.GetLocationsData
	getLocationsResponseBody struct {
		apiResponseBody
		Data locations.GetLocationsResultPayload `json:"data"`
	}
	addLocationAisleRequestBody  locations.AddAisleData
	addLocationAisleResponseBody struct {
		apiResponseBody
		Data locations.AddAisleResultPayload `json:"data"`
	}
	updateLocationAisleRequestBody  locations.UpdateAisleData
	updateLocationAisleResponseBody struct {
		apiResponseBody
		Data locations.UpdateAisleResultPayload `json:"data"`
	}
	deleteLocationAisleResponseBody struct {
		apiResponseBody
	}
	getLocationAisleResponseBody struct {
		apiResponseBody
		Data locations.GetAisleResultPayload `json:"data"`
	}
	getLocationAislesRequestBody  locations.GetAislesData
	getLocationAislesResponseBody struct {
		apiResponseBody
		Data locations.GetAislesResultPayload `json:"data"`
	}
)

type (
	addLocationBayRequestBody  locations.AddBayData
	addLocationBayResponseBody struct {
		apiResponseBody
		Data locations.AddBayResultPayload `json:"data"`
	}
	updateLocationBayRequestBody  locations.UpdateBayData
	updateLocationBayResponseBody struct {
		apiResponseBody
		Data locations.UpdateBayResultPayload `json:"data"`
	}
	deleteLocationBayResponseBody struct {
		apiResponseBody
	}
	getLocationBayResponseBody struct {
		apiResponseBody
		Data locations.GetBayResultPayload `json:"data"`
	}
	getLocationBaysRequestBody  locations.GetBaysData
	getLocationBaysResponseBody struct {
		apiResponseBody
		Data locations.GetBaysResultPayload `json:"data"`
	}
)

type (
	addLocationShelfRequestBody  locations.AddShelfData
	addLocationShelfResponseBody struct {
		apiResponseBody
		Data locations.AddShelfResultPayload `json:"data"`
	}
	updateLocationShelfRequestBody  locations.UpdateShelfData
	updateLocationShelfResponseBody struct {
		apiResponseBody
		Data locations.UpdateShelfResultPayload `json:"data"`
	}
	deleteLocationShelfResponseBody struct {
		apiResponseBody
	}
	getLocationShelfResponseBody struct {
		apiResponseBody
		Data locations.GetShelfResultPayload `json:"data"`
	}
	getLocationShelvesRequestBody  locations.GetShelvesData
	getLocationShelvesResponseBody struct {
		apiResponseBody
		Data locations.GetShelvesResultPayload `json:"data"`
	}
)

type (
	addLocationShelfLevelRequestBody  locations.AddShelfLevelData
	addLocationShelfLevelResponseBody struct {
		apiResponseBody
		Data locations.AddShelfLevelResultPayload `json:"data"`
	}
	updateLocationShelfLevelRequestBody  locations.UpdateShelfLevelData
	updateLocationShelfLevelResponseBody struct {
		apiResponseBody
		Data locations.UpdateShelfLevelResultPayload `json:"data"`
	}
	deleteLocationShelfLevelResponseBody struct {
		apiResponseBody
	}
	getLocationShelfLevelResponseBody struct {
		apiResponseBody
		Data locations.GetShelfLevelResultPayload `json:"data"`
	}
	getLocationShelfLevelsRequestBody  locations.GetShelfLevelsData
	getLocationShelfLevelsResponseBody struct {
		apiResponseBody
		Data locations.GetShelfLevelsResultPayload `json:"data"`
	}
)

type (
	addLocationBinRequestBody  locations.AddBinData
	addLocationBinResponseBody struct {
		apiResponseBody
		Data locations.AddBinResultPayload `json:"data"`
	}
	updateLocationBinRequestBody  locations.UpdateBinData
	updateLocationBinResponseBody struct {
		apiResponseBody
		Data locations.UpdateBinResultPayload `json:"data"`
	}
	deleteLocationBinResponseBody struct {
		apiResponseBody
	}
	getLocationBinResponseBody struct {
		apiResponseBody
		Data locations.GetBinResultPayload `json:"data"`
	}
	getLocationBinsRequestBody  locations.GetBinsData
	getLocationBinsResponseBody struct {
		apiResponseBody
		Data locations.GetBinsResultPayload `json:"data"`
	}
)

func (h *handler) getLocations(context *gin.Context) {
	var (
		requestBody  getLocationsRequestBody
		responseBody getLocationsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	result, err := h.services.Locations().GetLocations(locations.GetLocationsData(requestBody))
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

func (h *handler) addLocationAisle(context *gin.Context) {
	var (
		requestBody  addLocationAisleRequestBody
		responseBody addLocationAisleResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().AddAisle(locations.AddAisleData(requestBody))
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

func (h *handler) updateLocationAisle(context *gin.Context) {
	var (
		requestBody  updateLocationAisleRequestBody
		responseBody updateLocationAisleResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().UpdateAisle(locations.UpdateAisleData(requestBody))
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

func (h *handler) deleteLocationAisle(context *gin.Context) {
	var responseBody deleteLocationAisleResponseBody
	result, err := h.services.Locations().DeleteAisle(locations.DeleteAisleData{
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

func (h *handler) getLocationAisle(context *gin.Context) {
	var responseBody getLocationAisleResponseBody
	id := context.Param("id")
	result, err := h.services.Locations().GetAisle(locations.GetAisleData{
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

func (h *handler) getLocationAisles(context *gin.Context) {
	var (
		requestBody  getLocationAislesRequestBody
		responseBody getLocationAislesResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().GetAisles(locations.GetAislesData(requestBody))
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

func (h *handler) addLocationBay(context *gin.Context) {
	var (
		requestBody  addLocationBayRequestBody
		responseBody addLocationBayResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().AddBay(locations.AddBayData(requestBody))
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

func (h *handler) updateLocationBay(context *gin.Context) {
	var (
		requestBody  updateLocationBayRequestBody
		responseBody updateLocationBayResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().UpdateBay(locations.UpdateBayData(requestBody))
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

func (h *handler) deleteLocationBay(context *gin.Context) {
	var responseBody deleteLocationBayResponseBody
	result, err := h.services.Locations().DeleteBay(locations.DeleteBayData{
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

func (h *handler) getLocationBay(context *gin.Context) {
	var responseBody getLocationBayResponseBody
	id := context.Param("id")
	result, err := h.services.Locations().GetBay(locations.GetBayData{
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

func (h *handler) getLocationBays(context *gin.Context) {
	var (
		requestBody  getLocationBaysRequestBody
		responseBody getLocationBaysResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().GetBays(locations.GetBaysData(requestBody))
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

func (h *handler) addLocationShelf(context *gin.Context) {
	var (
		requestBody  addLocationShelfRequestBody
		responseBody addLocationShelfResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().AddShelf(locations.AddShelfData(requestBody))
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

func (h *handler) updateLocationShelf(context *gin.Context) {
	var (
		requestBody  updateLocationShelfRequestBody
		responseBody updateLocationShelfResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().UpdateShelf(locations.UpdateShelfData(requestBody))
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

func (h *handler) deleteLocationShelf(context *gin.Context) {
	var responseBody deleteLocationShelfResponseBody
	result, err := h.services.Locations().DeleteShelf(locations.DeleteShelfData{
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

func (h *handler) getLocationShelf(context *gin.Context) {
	var responseBody getLocationShelfResponseBody
	id := context.Param("id")
	result, err := h.services.Locations().GetShelf(locations.GetShelfData{
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

func (h *handler) getLocationShelves(context *gin.Context) {
	var (
		requestBody  getLocationShelvesRequestBody
		responseBody getLocationShelvesResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().GetShelves(locations.GetShelvesData(requestBody))
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

func (h *handler) addLocationShelfLevel(context *gin.Context) {
	var (
		requestBody  addLocationShelfLevelRequestBody
		responseBody addLocationShelfLevelResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().AddShelfLevel(locations.AddShelfLevelData(requestBody))
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

func (h *handler) updateLocationShelfLevel(context *gin.Context) {
	var (
		requestBody  updateLocationShelfLevelRequestBody
		responseBody updateLocationShelfLevelResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().UpdateShelfLevel(locations.UpdateShelfLevelData(requestBody))
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

func (h *handler) deleteLocationShelfLevel(context *gin.Context) {
	var responseBody deleteLocationShelfLevelResponseBody
	result, err := h.services.Locations().DeleteShelfLevel(locations.DeleteShelfLevelData{
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

func (h *handler) getLocationShelfLevel(context *gin.Context) {
	var responseBody getLocationShelfLevelResponseBody
	id := context.Param("id")
	result, err := h.services.Locations().GetShelfLevel(locations.GetShelfLevelData{
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

func (h *handler) getLocationShelfLevels(context *gin.Context) {
	var (
		requestBody  getLocationShelfLevelsRequestBody
		responseBody getLocationShelfLevelsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().GetShelfLevels(locations.GetShelfLevelsData(requestBody))
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

func (h *handler) addLocationBin(context *gin.Context) {
	var (
		requestBody  addLocationBinRequestBody
		responseBody addLocationBinResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().AddBin(locations.AddBinData(requestBody))
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

func (h *handler) updateLocationBin(context *gin.Context) {
	var (
		requestBody  updateLocationBinRequestBody
		responseBody updateLocationBinResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().UpdateBin(locations.UpdateBinData(requestBody))
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

func (h *handler) deleteLocationBin(context *gin.Context) {
	var responseBody deleteLocationBinResponseBody
	result, err := h.services.Locations().DeleteBin(locations.DeleteBinData{
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

func (h *handler) getLocationBin(context *gin.Context) {
	var responseBody getLocationBinResponseBody
	id := context.Param("id")
	result, err := h.services.Locations().GetBin(locations.GetBinData{
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

func (h *handler) getLocationBins(context *gin.Context) {
	var (
		requestBody  getLocationBinsRequestBody
		responseBody getLocationBinsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Locations().GetBins(locations.GetBinsData(requestBody))
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
