package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/inventory"

	"github.com/gin-gonic/gin"
)

type (
	addInventoryRequestBody  inventory.AddData
	addInventoryResponseBody struct {
		apiResponseBody
		Data inventory.AddResultPayload `json:"data"`
	}
	updateInventoryRequestBody  inventory.UpdateData
	updateInventoryResponseBody struct {
		apiResponseBody
		Data inventory.UpdateResultPayload `json:"data"`
	}
	deleteInventoryResponseBody struct {
		apiResponseBody
	}
	getInventoryItemResponseBody struct {
		apiResponseBody
		Data inventory.GetOneResultPayload `json:"data"`
	}
	getInventoryRequestBody  inventory.GetManyData
	getInventoryResponseBody struct {
		apiResponseBody
		Data inventory.GetManyResultPayload `json:"data"`
	}
	getInventorySummaryRequestBody  inventory.GetSummaryData
	getInventorySummaryResponseBody struct {
		apiResponseBody
		Data inventory.GetSummaryResultPayload `json:"data"`
	}
)

type (
	addInventoryMovementRequestBody  inventory.AddMovementData
	addInventoryMovementResponseBody struct {
		apiResponseBody
		Data inventory.AddMovementResultPayload `json:"data"`
	}
	updateInventoryMovementRequestBody  inventory.UpdateMovementData
	updateInventoryMovementResponseBody struct {
		apiResponseBody
		Data inventory.UpdateMovementResultPayload `json:"data"`
	}
	deleteInventoryMovementResponseBody struct {
		apiResponseBody
	}
	getInventoryMovementResponseBody struct {
		apiResponseBody
		Data inventory.GetMovementResultPayload `json:"data"`
	}
	getInventoryMovementsRequestBody  inventory.GetMovementsData
	getInventoryMovementsResponseBody struct {
		apiResponseBody
		Data inventory.GetMovementsResultPayload `json:"data"`
	}
)

func (h *handler) addInventory(context *gin.Context) {
	var (
		requestBody  addInventoryRequestBody
		responseBody addInventoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().Add(inventory.AddData(requestBody))
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

func (h *handler) updateInventory(context *gin.Context) {
	var (
		requestBody  updateInventoryRequestBody
		responseBody updateInventoryResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().Update(inventory.UpdateData(requestBody))
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

func (h *handler) deleteInventory(context *gin.Context) {
	var responseBody deleteInventoryResponseBody
	result, err := h.services.Inventory().Delete(inventory.DeleteData{
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

func (h *handler) getInventoryItem(context *gin.Context) {
	var responseBody getInventoryItemResponseBody
	id := context.Param("id")
	result, err := h.services.Inventory().GetOne(inventory.GetOneData{
		ID:         &id,
		LocationID: context.Param("lid"),
		Related:    true,
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

func (h *handler) getInventory(context *gin.Context) {
	var (
		requestBody  getInventoryRequestBody
		responseBody getInventoryResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.Related = true
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().GetMany(inventory.GetManyData(requestBody))
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

func (h *handler) getInventorySummary(context *gin.Context) {
	var (
		requestBody  getInventorySummaryRequestBody
		responseBody getInventorySummaryResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().GetSummary(inventory.GetSummaryData(requestBody))
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

func (h *handler) addInventoryMovement(context *gin.Context) {
	var (
		requestBody  addInventoryMovementRequestBody
		responseBody addInventoryMovementResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().AddMovement(inventory.AddMovementData(requestBody))
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

func (h *handler) updateInventoryMovement(context *gin.Context) {
	var (
		requestBody  updateInventoryMovementRequestBody
		responseBody updateInventoryMovementResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().UpdateMovement(inventory.UpdateMovementData(requestBody))
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

func (h *handler) deleteInventoryMovement(context *gin.Context) {
	var responseBody deleteInventoryMovementResponseBody
	result, err := h.services.Inventory().DeleteMovement(inventory.DeleteMovementData{
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

func (h *handler) getInventoryMovement(context *gin.Context) {
	var responseBody getInventoryMovementResponseBody
	id := context.Param("id")
	result, err := h.services.Inventory().GetMovement(inventory.GetMovementData{
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

func (h *handler) getInventoryMovements(context *gin.Context) {
	var (
		requestBody  getInventoryMovementsRequestBody
		responseBody getInventoryMovementsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Inventory().GetMovements(inventory.GetMovementsData(requestBody))
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
