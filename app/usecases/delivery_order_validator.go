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
	"time"

	"github.com/bxcodec/dbresolver"

	"github.com/gin-gonic/gin"
)

type DeliveryOrderValidatorInterface interface {
	CreateDeliveryOrderValidator(*models.DeliveryOrderStoreRequest, *gin.Context) error
	GetDeliveryOrderValidator(*gin.Context) (*models.DeliveryOrderRequest, error)
	GetDeliveryOrderBySalesmanIDValidator(*gin.Context) (*models.DeliveryOrderRequest, error)
	UpdateDeliveryOrderByIDValidator(*models.DeliveryOrderUpdateByIDRequest, *gin.Context) error
	DeleteDeliveryOrderByIDValidator(string, *gin.Context) error
}

type DeliveryOrderValidator struct {
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitDeliveryOrderValidator(requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderValidatorInterface {
	return &DeliveryOrderValidator{
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (d *DeliveryOrderValidator) CreateDeliveryOrderValidator(insertRequest *models.DeliveryOrderStoreRequest, ctx *gin.Context) error {
	var result baseModel.Response
	dateField := []*models.DateInputRequest{
		{
			Field: "so_date",
			Value: insertRequest.DoRefDate,
		},
	}
	err := d.requestValidationMiddleware.DateInputValidation(ctx, dateField, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return err
	}

	uniqueField := []*models.UniqueRequest{
		{
			Table: constants.DELIVERY_ORDERS_TABLE,
			Field: "do_ref_code",
			Value: insertRequest.DoRefCode,
		},
	}

	err = d.requestValidationMiddleware.UniqueValidation(ctx, uniqueField)
	if err != nil {
		fmt.Println("Error unique validation", err)
		return err
	}
	mustActiveField417 := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "sales_order_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", insertRequest.SalesOrderID),
		},
		{
			Table:    "warehouses a JOIN agents b JOIN sales_orders c ON a.owner_id = b.id AND c.agent_id = b.id",
			ReqField: "warehouse_id",
			Clause:   fmt.Sprintf("c.id = %d AND a.id = %d AND b.deleted_at IS NULL AND a.`status` = 1", insertRequest.SalesOrderID, insertRequest.WarehouseID),
		},
		{
			Table:    "stores a JOIN sales_orders b ON b.store_id = a.id",
			ReqField: "stores_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", insertRequest.SalesOrderID),
		},
		{
			Table:    "brands a JOIN sales_orders b ON b.brand_id = a.id",
			ReqField: "brands_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.status_active = 1 AND b.deleted_at IS NULL", insertRequest.SalesOrderID),
		},
	}
	now := time.Now().UTC().Add(7 * time.Hour)
	sDoDate := now.Format("2006-01-02")
	sDoDateEqualMonth := fmt.Sprintf("MONTH(so_date) = MONTH('%s') AND MONTH(so_date) = MONTH('%s')", insertRequest.DoRefDate, sDoDate)
	sDoDateHigherOrEqualSoDate := fmt.Sprintf("DAY(so_date) <= DAY('%s') AND DAY(so_date) <= DAY('%s')", insertRequest.DoRefDate, sDoDate)
	// sDoDateLowerOrEqualSoRefDate := fmt.Sprintf("DAY(so_ref_date) >= DAY('%s') AND DAY(so_ref_date) >= DAY(%s)", insertRequest.DoRefDate, sDoDate)
	sDoDateLowerOrEqualToday := fmt.Sprintf("DAY('%s') <= DAY('%s')", insertRequest.DoRefDate, sDoDate)
	sSoDateEqualDoDate := fmt.Sprintf("IF(DAY(so_date) = DAY('%[2]s'), IF(DAY('%[2]s') = DAY('%[1]s'), TRUE, FALSE), TRUE)", insertRequest.DoRefDate, sDoDate)
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:         "sales_orders",
			ReqField:      "so_date",
			Clause:        fmt.Sprintf("id = %d AND %s AND %s AND %s AND %s ", insertRequest.SalesOrderID, sDoDateEqualMonth, sDoDateHigherOrEqualSoDate, sDoDateLowerOrEqualToday, sSoDateEqualDoDate),
			CustomMessage: "do_date and do_ref_date must be equal less than today, must be equal more than so_date and must be in the current month",
		},
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_orders",
			TableJoin:       "order_statuses",
			ForeignKey:      "order_status_id",
			SelectedCollumn: "order_statuses.name",
			Clause:          fmt.Sprintf("sales_orders.id = %d AND sales_orders.order_status_id NOT IN (5,7)", insertRequest.SalesOrderID),
			MessageFormat:   "Status Sales Order <result>",
		},
	}
	totalQty := 0
	for _, x := range insertRequest.DeliveryOrderDetails {
		if x.Qty < 0 {
			err = helper.NewError(fmt.Sprintf("qty sales order detail %d must equal or higher than 0", x.SoDetailID))
			result.StatusCode = http.StatusUnprocessableEntity
			result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
			ctx.JSON(result.StatusCode, result)
			return err
		}

		totalQty += x.Qty
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:         "sales_orders a JOIN sales_order_details b ON b.sales_order_id = a.id",
			ReqField:      "sales_order_id",
			Clause:        fmt.Sprintf("b.id = %d AND a.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID, insertRequest.SalesOrderID),
			CustomMessage: fmt.Sprintf("sales order detail = %d tidak ditemukan", x.SoDetailID),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "products a JOIN sales_order_details b ON b.product_id = a.id",
			ReqField: "product_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "uoms a JOIN sales_order_details b ON b.uom_id = a.id",
			ReqField: "uoms_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID),
		})
		mustActiveField422 = append(mustActiveField422, &models.MustActiveRequest{
			Table:         "sales_order_details",
			ReqField:      "residual_qty",
			Clause:        fmt.Sprintf("id = %d AND deleted_at IS NULL AND residual_qty >= %d", x.SoDetailID, x.Qty),
			CustomMessage: fmt.Sprintf("Residual Qty SO Detail %d must be higher than or equal delivery order qty", x.SoDetailID),
		})
		mustEmpties = append(mustEmpties, &models.MustEmptyValidationRequest{
			Table:           "sales_order_details",
			TableJoin:       "order_statuses",
			ForeignKey:      "order_status_id",
			SelectedCollumn: "order_statuses.name",
			Clause:          fmt.Sprintf("sales_order_details.id = %d AND sales_order_details.order_status_id NOT IN (11, 13)", x.SoDetailID),
			MessageFormat:   fmt.Sprintf("Status Sales Order Detail %d <result>", x.SoDetailID),
		})
	}

	err = d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField417)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustActiveValidation422(ctx, mustActiveField422)
	if err != nil {
		return err
	}

	if totalQty <= 0 {
		err = helper.NewError("total qty must higher than 0")
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return err
	}
	return nil
}
func (d *DeliveryOrderValidator) UpdateDeliveryOrderByIDValidator(insertRequest *models.DeliveryOrderUpdateByIDRequest, ctx *gin.Context) error {
	uniqueField := []*models.UniqueRequest{
		{
			Table: constants.DELIVERY_ORDERS_TABLE,
			Field: "do_ref_code",
			Value: insertRequest.DoRefCode,
		},
	}

	err := d.requestValidationMiddleware.UniqueValidation(ctx, uniqueField)
	if err != nil {
		fmt.Println("Error unique validation", err)
		return err
	}
	return nil
}
func (d *DeliveryOrderValidator) DeleteDeliveryOrderByIDValidator(sId string, ctx *gin.Context) error {
	var result baseModel.Response
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, err.Error())
		ctx.JSON(result.StatusCode, result)
		return err
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "delivery_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
		},
	}

	err = d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return err
	}

	return nil
}
func (c *DeliveryOrderValidator) GetDeliveryOrderValidator(ctx *gin.Context) (*models.DeliveryOrderRequest, error) {
	var result baseModel.Response

	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intSalesOrderID, err := c.getIntQueryWithDefault("sales_order_id", "0", false, ctx)
	if err != nil {
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

	intOrderStatusID, err := c.getIntQueryWithDefault("order_status_id", "0", false, ctx)
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

	deliveryOrderReqeuest := &models.DeliveryOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		ID:                intID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		SoCode:            c.getQueryWithDefault("so_code", "", ctx),
		StartDoDate:       c.getQueryWithDefault("start_do_date", "", ctx),
		EndDoDate:         c.getQueryWithDefault("end_do_date", "", ctx),
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    c.getQueryWithDefault("start_created_at", "", ctx),
		EndCreatedAt:      c.getQueryWithDefault("end_created_at", "", ctx),
		UpdatedAt:         c.getQueryWithDefault("updated_at", "", ctx),
	}
	return deliveryOrderReqeuest, nil
}
func (c *DeliveryOrderValidator) GetDeliveryOrderBySalesmanIDValidator(ctx *gin.Context) (*models.DeliveryOrderRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)

	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intSalesOrderID, err := c.getIntQueryWithDefault("sales_order_id", "0", false, ctx)
	if err != nil {
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

	intOrderStatusID, err := c.getIntQueryWithDefault("order_status_id", "0", false, ctx)
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

	deliveryOrderRequest := &models.DeliveryOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		ID:                intID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		StartDoDate:       c.getQueryWithDefault("start_do_date", "", ctx),
		EndDoDate:         c.getQueryWithDefault("end_do_date", "", ctx),
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    c.getQueryWithDefault("start_created_at", "", ctx),
		EndCreatedAt:      c.getQueryWithDefault("end_created_at", "", ctx),
		UpdatedAt:         c.getQueryWithDefault("updated_at", "", ctx),
	}
	return deliveryOrderRequest, nil
}

func (d *DeliveryOrderValidator) getQueryWithDefault(param string, empty string, ctx *gin.Context) string {
	result, isStartCreatedAt := ctx.GetQuery(param)
	if isStartCreatedAt == false {
		result = empty
	}
	return result
}

func (d *DeliveryOrderValidator) getIntQueryWithDefault(param string, empty string, isNotZero bool, ctx *gin.Context) (int, error) {
	var response baseModel.Response
	sResult := d.getQueryWithDefault(param, empty, ctx)
	result, err := strconv.Atoi(sResult)
	if err != nil {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, result)
		return 0, err
	}
	if result == 0 && isNotZero {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer > 0", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, result)
		return 0, err
	}
	return result, nil
}
