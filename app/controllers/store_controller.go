package controllers

import (
	"context"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type StoreControllerInterface interface {
	GetSalesOrders(ctx *gin.Context)
	GetDeliveryOrders(ctx *gin.Context)
}

type storeController struct {
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	salesOrderValidator         usecases.SalesOrderValidatorInterface
	deliveryOrderValidator      usecases.DeliveryOrderValidatorInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitStoreController(salesOrderUseCase usecases.SalesOrderUseCaseInterface, deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, salesOrderValidator usecases.SalesOrderValidatorInterface, deliveryOrderValidator usecases.DeliveryOrderValidatorInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) StoreControllerInterface {
	return &storeController{
		salesOrderUseCase:           salesOrderUseCase,
		deliveryOrderUseCase:        deliveryOrderUseCase,
		salesOrderValidator:         salesOrderValidator,
		deliveryOrderValidator:      deliveryOrderValidator,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *storeController) GetSalesOrders(ctx *gin.Context) {
	storeIds := ctx.Param("id")
	storeId, err := strconv.Atoi(storeIds)

	if err != nil {
		err = helper.NewError("Parameter 'store id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	salesOrderRequest, err := c.salesOrderValidator.GetSalesOrderValidator(ctx)
	if err != nil {
		return
	}

	salesOrderRequest.StoreID = storeId

	salesOrders, errorLog := c.salesOrderUseCase.Get(salesOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: salesOrders.SalesOrders, Total: salesOrders.Total, StatusCode: http.StatusOK})
	return
}

func (c *storeController) GetDeliveryOrders(ctx *gin.Context) {
	storeIds := ctx.Param("id")
	storeId, err := strconv.Atoi(storeIds)

	if err != nil {
		err = helper.NewError("Parameter 'store id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrderRequest.StoreID = storeId

	deliveryOrders, errorLog := c.deliveryOrderUseCase.Get(deliveryOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders.DeliveryOrders, Total: deliveryOrders.Total, StatusCode: http.StatusOK})
	return
}
