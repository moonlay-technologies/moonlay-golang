package usecases

import (
	"context"
	"fmt"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	baseModel "order-service/global/utils/model"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type SalesOrderValidatorInterface interface {
	GetSalesOrderValidator(*gin.Context) (*models.SalesOrderRequest, error)
}

type SalesOrderValidator struct {
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitSalesOrderValidator(requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) SalesOrderValidatorInterface {
	return &SalesOrderValidator{
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *SalesOrderValidator) GetSalesOrderValidator(ctx *gin.Context) (*models.SalesOrderRequest, error) {
	var result baseModel.Response

	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intStoreID, err := c.getIntQueryWithDefault("store_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intBrandID, err := c.getIntQueryWithDefault("brand_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderSourceID, err := c.getIntQueryWithDefault("order_source_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderStatusID, err := c.getIntQueryWithDefault("order_status_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProductID, err := c.getIntQueryWithDefault("product_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCategoryID, err := c.getIntQueryWithDefault("category_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intSalesmanID, err := c.getIntQueryWithDefault("salesman_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProvinceID, err := c.getIntQueryWithDefault("province_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCityID, err := c.getIntQueryWithDefault("city_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDistrictID, err := c.getIntQueryWithDefault("district_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intVillageID, err := c.getIntQueryWithDefault("village_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startSoDate, dateFields := c.getQueryWithDateValidation("start_so_date", "", dateFields, ctx)

	endSoDate, dateFields := c.getQueryWithDateValidation("end_so_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	salesOrderRequest := &models.SalesOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                intID,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}
	return salesOrderRequest, nil
}

func (d *SalesOrderValidator) getQueryWithDefault(param string, empty string, ctx *gin.Context) string {
	result, isStartCreatedAt := ctx.GetQuery(param)
	if isStartCreatedAt == false {
		result = empty
	}
	return result
}

func (d *SalesOrderValidator) getQueryWithDateValidation(param string, empty string, dateFields []*models.DateInputRequest, ctx *gin.Context) (string, []*models.DateInputRequest) {
	result, isResultExist := ctx.GetQuery(param)
	if isResultExist == false {
		result = empty
	} else {
		dateFields = append(dateFields, &models.DateInputRequest{
			Field: param,
			Value: result,
		})
	}
	return result, dateFields
}

func (d *SalesOrderValidator) getIntQueryWithDefault(param string, empty string, isNotZero bool, ctx *gin.Context) (int, error) {
	var response baseModel.Response
	sResult := d.getQueryWithDefault(param, empty, ctx)
	result, err := strconv.Atoi(sResult)
	if err != nil {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, response)
		return 0, err
	}
	if result == 0 && isNotZero {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer > 0", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, response)
		return 0, err
	}
	return result, nil
}
