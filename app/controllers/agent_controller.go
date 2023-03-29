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

type AgentControllerInterface interface {
	GetSalesOrders(ctx *gin.Context)
	GetDeliveryOrders(ctx *gin.Context)
}

type agentController struct {
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	salesOrderValidator         usecases.SalesOrderValidatorInterface
	deliveryOrderValidator      usecases.DeliveryOrderValidatorInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitAgentController(salesOrderUseCase usecases.SalesOrderUseCaseInterface, deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, salesOrderValidator usecases.SalesOrderValidatorInterface, deliveryOrderValidator usecases.DeliveryOrderValidatorInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) AgentControllerInterface {
	return &agentController{
		salesOrderUseCase:           salesOrderUseCase,
		deliveryOrderUseCase:        deliveryOrderUseCase,
		salesOrderValidator:         salesOrderValidator,
		deliveryOrderValidator:      deliveryOrderValidator,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *agentController) GetSalesOrders(ctx *gin.Context) {
	agentIds := ctx.Param("id")
	agentId, err := strconv.Atoi(agentIds)

	if err != nil {
		err = helper.NewError("Parameter 'agent id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	salesOrderRequest, err := c.salesOrderValidator.GetSalesOrderValidator(ctx)
	if err != nil {
		return
	}

	salesOrderRequest.AgentID = agentId

	salesOrders, errorLog := c.salesOrderUseCase.Get(salesOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: salesOrders.SalesOrders, Total: salesOrders.Total, StatusCode: http.StatusOK})
}

func (c *agentController) GetDeliveryOrders(ctx *gin.Context) {
	agentIds := ctx.Param("id")
	agentId, err := strconv.Atoi(agentIds)

	if err != nil {
		err = helper.NewError("Parameter 'agent id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrderRequest.AgentID = agentId

	deliveryOrders, errorLog := c.deliveryOrderUseCase.Get(deliveryOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders.DeliveryOrders, Total: deliveryOrders.Total, StatusCode: http.StatusOK})
}
