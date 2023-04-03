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
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"

	"github.com/gin-gonic/gin"
)

type DeliveryOrderValidatorInterface interface {
	CreateDeliveryOrderValidator(*models.DeliveryOrderStoreRequest, *gin.Context) error
	GetDeliveryOrderValidator(*gin.Context) (*models.DeliveryOrderRequest, error)
	ExportDeliveryOrderValidator(*gin.Context) (*models.DeliveryOrderExportRequest, error)
	ExportDeliveryOrderDetailValidator(*gin.Context) (*models.DeliveryOrderDetailExportRequest, error)
	GetDeliveryOrderDetailValidator(*gin.Context) (*models.DeliveryOrderDetailOpenSearchRequest, error)
	GetDeliveryOrderDetailByDoIDValidator(*gin.Context) (*models.DeliveryOrderDetailRequest, error)
	GetDeliveryOrderBySalesmanIDValidator(*gin.Context) (*models.DeliveryOrderRequest, error)
	GetDeliveryOrderSyncToKafkaHistoriesValidator(*gin.Context) (*models.DeliveryOrderEventLogRequest, error)
	GetDeliveryOrderJourneysValidator(*gin.Context) (*models.DeliveryOrderJourneysRequest, error)
	GetDOUploadHistoriesValidator(*gin.Context) (*models.GetDoUploadHistoriesRequest, error)
	UpdateDeliveryOrderByIDValidator(int, *models.DeliveryOrderUpdateByIDRequest, *gin.Context) error
	UpdateDeliveryOrderDetailByDoIDValidator(int, []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, *gin.Context) error
	UpdateDeliveryOrderDetailByIDValidator(int, *models.DeliveryOrderDetailUpdateByIDRequest, *gin.Context) error
	DeleteDeliveryOrderByIDValidator(string, *gin.Context) error
	DeleteDeliveryOrderDetailByIDValidator(string, *gin.Context) error
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
	sDoDateHigherOrEqualSoDate := fmt.Sprintf("DATE(so_date) <= DATE('%s') AND DATE(so_date) <= DATE('%s')", insertRequest.DoRefDate, sDoDate)
	// sDoDateLowerOrEqualSoRefDate := fmt.Sprintf("DAY(so_ref_date) >= DAY('%s') AND DAY(so_ref_date) >= DAY(%s)", insertRequest.DoRefDate, sDoDate)
	sDoDateLowerOrEqualToday := fmt.Sprintf("DATE('%s') <= DATE('%s')", insertRequest.DoRefDate, sDoDate)
	sSoDateEqualDoDate := fmt.Sprintf("IF(DATE(so_date) = DATE('%[2]s'), IF(DATE('%[2]s') = DATE('%[1]s'), TRUE, FALSE), TRUE)", insertRequest.DoRefDate, sDoDate)
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
			Table:           "sales_orders s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("s.id = %d AND s.order_status_id NOT IN (5,7)", insertRequest.SalesOrderID),
			MessageFormat:   "Status Sales Order <result>",
		},
	}
	totalQty := 0
	for _, x := range insertRequest.DeliveryOrderDetails {
		if x.Qty < 0 {
			err = helper.NewError(fmt.Sprintf("qty delivery order detail %d must equal or higher than 0", x.SoDetailID))
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
			Table:           "sales_order_details d JOIN order_statuses s ON d.order_status_id = s.id",
			SelectedCollumn: "s.name",
			Clause:          fmt.Sprintf("d.id = %d AND d.order_status_id NOT IN (11, 13)", x.SoDetailID),
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

func (d *DeliveryOrderValidator) UpdateDeliveryOrderByIDValidator(id int, insertRequest *models.DeliveryOrderUpdateByIDRequest, ctx *gin.Context) error {
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

	mustActiveField417 := []*models.MustActiveRequest{
		{
			Table:    "delivery_orders",
			ReqField: "delivery_order_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
		},
		{
			Table:    "sales_orders s JOIN delivery_orders d ON d.sales_order_id = s.id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("d.id = %d AND s.deleted_at IS NULL", id),
		},
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "delivery_orders d JOIN order_statuses s ON d.order_status_id = s.id",
			SelectedCollumn: "s.name",
			Clause:          fmt.Sprintf("d.id = %d AND d.order_status_id NOT IN (17)", id),
			MessageFormat:   "Status Delivery Order <result>",
		},
		{
			Table:           "sales_orders s JOIN delivery_orders d JOIN order_statuses o ON d.sales_order_id = s.id AND o.id = s.order_status_id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("d.id = %d AND s.order_status_id NOT IN (5,7,8)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
	}

	for _, v := range insertRequest.DeliveryOrderDetails {
		if v.Qty < 0 {
			err := helper.NewError("qty must be higher or equal 0")
			result := &baseModel.Response{
				Error: helper.NewWriteLog(baseModel.ErrorLog{
					Message:       err.Error(),
					SystemMessage: err,
					StatusCode:    422,
				}),
			}
			ctx.JSON(result.StatusCode, result)
			return err
		}
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:         "delivery_order_details",
			ReqField:      "delivery_order_detail_id",
			Clause:        fmt.Sprintf("id = %d AND delivery_order_id = %d AND deleted_at IS NULL", v.ID, id),
			CustomMessage: fmt.Sprintf("delivery_order_detail_id %d not found in delivery_order id %d", v.ID, id),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "sales_order_details s JOIN delivery_order_details d ON d.so_detail_id = s.id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("d.id = %d AND s.deleted_at IS NULL", v.ID),
		})
		mustEmpties = append(mustEmpties, &models.MustEmptyValidationRequest{
			Table:           "sales_order_details s JOIN delivery_order_details d ON s.id = d.so_detail_id",
			SelectedCollumn: "s.residual_qty+d.qty",
			Clause:          fmt.Sprintf("d.id = %[1]d AND s.residual_qty+d.qty < %[2]d", v.ID, v.Qty),
			MessageFormat:   fmt.Sprintf("DO Detail %d must be lower than or equal residual qty (<result>)", v.ID),
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

	return nil
}

func (d *DeliveryOrderValidator) UpdateDeliveryOrderDetailByDoIDValidator(id int, insertRequest []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, ctx *gin.Context) error {
	mustActiveField417 := []*models.MustActiveRequest{
		{
			Table:    "delivery_orders",
			ReqField: "delivery_order_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
		},
		{
			Table:    "sales_orders s JOIN delivery_orders d ON d.sales_order_id = s.id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("d.id = %d AND s.deleted_at IS NULL", id),
		},
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "delivery_orders d JOIN order_statuses s ON d.order_status_id = s.id",
			SelectedCollumn: "s.name",
			Clause:          fmt.Sprintf("d.id = %d AND d.order_status_id NOT IN (17)", id),
			MessageFormat:   "Status Delivery Order <result>",
		},
		{
			Table:           "sales_orders s JOIN delivery_orders d JOIN order_statuses o ON d.sales_order_id = s.id AND o.id = s.order_status_id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("d.id = %d AND s.order_status_id NOT IN (5,7,8)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
	}

	for _, v := range insertRequest {
		if v.Qty < 0 {
			err := helper.NewError("qty must be higher or equal 0")
			result := &baseModel.Response{
				Error: helper.NewWriteLog(baseModel.ErrorLog{
					Message:       err.Error(),
					SystemMessage: err,
					StatusCode:    422,
				}),
			}
			ctx.JSON(result.StatusCode, result)
			return err
		}
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "delivery_order_details",
			ReqField: "delivery_order_detail_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", v.ID),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "sales_order_details s JOIN delivery_order_details d ON d.so_detail_id = s.id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("d.id = %d AND s.deleted_at IS NULL", v.ID),
		})
		mustEmpties = append(mustEmpties, &models.MustEmptyValidationRequest{
			Table:           "sales_order_details s JOIN delivery_order_details d ON s.id = d.so_detail_id",
			SelectedCollumn: "s.id",
			Clause:          fmt.Sprintf("d.id = %d AND s.qty < %d", v.ID, v.Qty),
			MessageFormat:   fmt.Sprintf("Qty SO Detail <result> FROM DO Detail %d must be higher than or equal delivery order qty", v.ID),
		})
	}

	err := d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField417)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeliveryOrderValidator) UpdateDeliveryOrderDetailByIDValidator(detailId int, insertRequest *models.DeliveryOrderDetailUpdateByIDRequest, ctx *gin.Context) error {
	mustActiveField417 := []*models.MustActiveRequest{
		{
			Table:    "delivery_orders d JOIN delivery_order_details dd ON d.id = dd.delivery_order_id",
			ReqField: "d.delivery_order_id",
			Clause:   fmt.Sprintf("dd.id = %d AND d.deleted_at IS NULL", detailId),
		},
		{
			Table:    "sales_orders s JOIN delivery_orders d ON d.sales_order_id = s.id JOIN delivery_order_details dd ON d.id = dd.delivery_order_id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("dd.id = %d AND s.deleted_at IS NULL", detailId),
		},
		{
			Table:    "delivery_order_details",
			ReqField: "delivery_order_detail_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", detailId),
		},
		{
			Table:    "sales_order_details s JOIN delivery_order_details d ON d.so_detail_id = s.id",
			ReqField: "sales-order_id",
			Clause:   fmt.Sprintf("d.id = %d AND s.deleted_at IS NULL", detailId),
		},
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "delivery_orders d JOIN order_statuses s ON d.order_status_id = s.id JOIN delivery_order_details dd ON d.id = dd.delivery_order_id",
			SelectedCollumn: "s.name",
			Clause:          fmt.Sprintf("dd.id = %d AND d.order_status_id NOT IN (17)", detailId),
			MessageFormat:   "Status Delivery Order <result>",
		},
		{
			Table:           "sales_orders s JOIN delivery_orders d JOIN order_statuses o ON d.sales_order_id = s.id AND o.id = s.order_status_id JOIN delivery_order_details dd ON d.id = dd.delivery_order_id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("dd.id = %d AND s.order_status_id NOT IN (5,7,8)", detailId),
			MessageFormat:   "Status Sales Order <result>",
		},
		{
			Table:           "sales_order_details s JOIN delivery_order_details d ON s.id = d.so_detail_id",
			SelectedCollumn: "s.id",
			Clause:          fmt.Sprintf("d.id = %d AND s.qty < %d", insertRequest.ID, insertRequest.Qty),
			MessageFormat:   fmt.Sprintf("Qty SO Detail <result> FROM DO Detail %d must be higher than or equal delivery order qty", insertRequest.ID),
		},
	}
	if insertRequest.Qty < 0 {
		err := helper.NewError("qty must be higher or equal 0")
		result := &baseModel.Response{
			Error: helper.NewWriteLog(baseModel.ErrorLog{
				Message:       err.Error(),
				SystemMessage: err,
				StatusCode:    422,
			}),
		}
		ctx.JSON(result.StatusCode, result)
		return err
	}

	err := d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField417)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeliveryOrderValidator) DeleteDeliveryOrderByIDValidator(sId string, ctx *gin.Context) error {
	var result baseModel.Response
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, err.Error())
		ctx.JSON(result.StatusCode, result)
		return err
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "delivery_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d", id),
		},
	}
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:         "delivery_orders",
			ReqField:      "id",
			Clause:        fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
			CustomMessage: fmt.Sprintf("DO id = %d was deleted", id),
		},
	}

	err = d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustActiveValidation422(ctx, mustActiveField422)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeliveryOrderValidator) DeleteDeliveryOrderDetailByIDValidator(sId string, ctx *gin.Context) error {
	var result baseModel.Response
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, err.Error())
		ctx.JSON(result.StatusCode, result)
		return err
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "delivery_order_details",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d", id),
		},
	}
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:         "delivery_order_details",
			ReqField:      "id",
			Clause:        fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
			CustomMessage: fmt.Sprintf("DO detail id = %d was deleted", id),
		},
	}

	err = d.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return err
	}

	err = d.requestValidationMiddleware.MustActiveValidation422(ctx, mustActiveField422)
	if err != nil {
		return err
	}

	return nil
}

func (c *DeliveryOrderValidator) GetDeliveryOrderValidator(ctx *gin.Context) (*models.DeliveryOrderRequest, error) {
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

	dateFields := []*models.DateInputRequest{}

	startDoDate, dateFields := c.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endDoDate, dateFields := c.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	updatedAt, dateFields := c.getQueryWithDateValidation("updated_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
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
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         updatedAt,
	}
	return deliveryOrderReqeuest, nil
}

func (c *DeliveryOrderValidator) ExportDeliveryOrderValidator(ctx *gin.Context) (*models.DeliveryOrderExportRequest, error) {
	var result baseModel.Response

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "store_id" && sortField != "created_at" && sortField != "updated_at" {
		err := helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
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

	dateFields := []*models.DateInputRequest{}

	startDoDate, dateFields := c.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endDoDate, dateFields := c.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	deliveryOrderReqeuest := &models.DeliveryOrderExportRequest{
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		FileType:          c.getQueryWithDefault("file_type", "xlsx", ctx),
		ID:                intID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         c.getQueryWithDefault("updated_at", "", ctx),
	}
	return deliveryOrderReqeuest, nil
}

func (c *DeliveryOrderValidator) ExportDeliveryOrderDetailValidator(ctx *gin.Context) (*models.DeliveryOrderDetailExportRequest, error) {
	var result baseModel.Response

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "store_id" && sortField != "created_at" && sortField != "updated_at" {
		err := helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDeliveryOrderID, err := c.getIntQueryWithDefault("do_id", "0", false, ctx)
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

	dateFields := []*models.DateInputRequest{}

	startDoDate, dateFields := c.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endDoDate, dateFields := c.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	deliveryOrderDetailExportRequest := &models.DeliveryOrderDetailExportRequest{
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		FileType:          c.getQueryWithDefault("file_type", "xlsx", ctx),
		ID:                intID,
		DeliveryOrderID:   intDeliveryOrderID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         c.getQueryWithDefault("updated_at", "", ctx),
	}
	return deliveryOrderDetailExportRequest, nil
}

func (c *DeliveryOrderValidator) GetDeliveryOrderDetailValidator(ctx *gin.Context) (*models.DeliveryOrderDetailOpenSearchRequest, error) {
	var result baseModel.Response

	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "order_status_id", ctx)

	if sortField != "order_status_id" && sortField != "do_code" && sortField != "so_code" && sortField != "agent_id" && sortField != "store_id" && sortField != "product_id" && sortField != "qty" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_code' or 'so_code' or 'agent_id' or 'store_id' or 'product_id' or 'qty'")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDeliveryOrderID, err := c.getIntQueryWithDefault("do_id", "0", false, ctx)
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

	intQty, err := c.getIntQueryWithDefault("qty", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startDoDate, dateFields := c.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endDoDate, dateFields := c.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	deliveryOrderReqeuest := &models.DeliveryOrderDetailOpenSearchRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		ID:                intID,
		DeliveryOrderID:   intDeliveryOrderID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		SoCode:            c.getQueryWithDefault("so_code", "", ctx),
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		Qty:               intQty,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         c.getQueryWithDefault("updated_at", "", ctx),
	}
	return deliveryOrderReqeuest, nil
}

func (c *DeliveryOrderValidator) GetDeliveryOrderDetailByDoIDValidator(ctx *gin.Context) (*models.DeliveryOrderDetailRequest, error) {
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

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intDoDetailID, err := c.getIntQueryWithDefault("do_detail_id", "0", false, ctx)
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

	intQty, err := c.getIntQueryWithDefault("qty", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startDoDate, dateFields := c.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endDoDate, dateFields := c.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	deliveryOrderReqeuest := &models.DeliveryOrderDetailRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		DoDetailID:        intDoDetailID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            c.getQueryWithDefault("do_code", "", ctx),
		SoCode:            c.getQueryWithDefault("so_code", "", ctx),
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         c.getQueryWithDefault("do_ref_code", "", ctx),
		DoRefDate:         c.getQueryWithDefault("do_ref_date", "", ctx),
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		Qty:               intQty,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
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

func (c *DeliveryOrderValidator) GetDeliveryOrderSyncToKafkaHistoriesValidator(ctx *gin.Context) (*models.DeliveryOrderEventLogRequest, error) {
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

	if sortField != "do_code" && sortField != "status" && sortField != "agent_name" && sortField != "created_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'do_code' or 'status' or 'agent_name' or 'created_at' ")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	salesOrderRequest := &models.DeliveryOrderEventLogRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue: strings.ToLower(c.getQueryWithDefault("global_search_value", "", ctx)),
		ID:                c.getQueryWithDefault("id", "", ctx),
		RequestID:         c.getQueryWithDefault("request_id", "", ctx),
		AgentID:           intAgentID,
		Status:            strings.ToLower(c.getQueryWithDefault("status", "", ctx)),
	}
	return salesOrderRequest, nil
}

func (c *DeliveryOrderValidator) GetDeliveryOrderJourneysValidator(ctx *gin.Context) (*models.DeliveryOrderJourneysRequest, error) {

	intDoID, err := c.getIntQueryWithDefault("do_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	deliveryOrderJourneysRequest := models.DeliveryOrderJourneysRequest{
		DoId:      intDoID,
		DoDate:    c.getQueryWithDefault("do_date", "", ctx),
		Status:    c.getQueryWithDefault("status", "", ctx),
		Remark:    c.getQueryWithDefault("remark", "", ctx),
		Reason:    c.getQueryWithDefault("reason", "", ctx),
		CreatedAt: c.getQueryWithDefault("created_at", "", ctx),
	}

	return &deliveryOrderJourneysRequest, nil
}

func (c *DeliveryOrderValidator) GetDOUploadHistoriesValidator(ctx *gin.Context) (*models.GetDoUploadHistoriesRequest, error) {
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

	if sortField != "agent_name" && sortField != "file_name" && sortField != "status" && sortField != "created_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'agent_name' or 'file_name' or 'status' or 'created_at' ")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intUploadedBy, err := c.getIntQueryWithDefault("uploaded_by", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startUploadAt, dateFields := c.getQueryWithDateValidation("start_upload_at", "", dateFields, ctx)

	endUploadAt, dateFields := c.getQueryWithDateValidation("end_upload_at", "", dateFields, ctx)

	finishProcessDateStart, dateFields := c.getQueryWithDateValidation("finish_process_date_start", "", dateFields, ctx)

	finishProcessDateEnd, dateFields := c.getQueryWithDateValidation("finish_process_date_end", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	getDoUploadHistoriesRequest := &models.GetDoUploadHistoriesRequest{
		ID:                     c.getQueryWithDefault("id", "", ctx),
		Page:                   pageInt,
		PerPage:                perPageInt,
		SortField:              sortField,
		SortValue:              strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue:      strings.ToLower(c.getQueryWithDefault("global_search_value", "", ctx)),
		RequestID:              c.getQueryWithDefault("request_id", "", ctx),
		FileName:               c.getQueryWithDefault("file_name", "", ctx),
		BulkCode:               c.getQueryWithDefault("bulk_code", "", ctx),
		AgentID:                intAgentID,
		Status:                 strings.ToLower(c.getQueryWithDefault("status", "", ctx)),
		UploadedBy:             intUploadedBy,
		StartUploadAt:          startUploadAt,
		EndUploadAt:            endUploadAt,
		FinishProcessDateStart: finishProcessDateStart,
		FinishProcessDateEnd:   finishProcessDateEnd,
	}

	return getDoUploadHistoriesRequest, nil
}

func (d *DeliveryOrderValidator) getQueryWithDefault(param string, empty string, ctx *gin.Context) string {
	result, isStartCreatedAt := ctx.GetQuery(param)
	if isStartCreatedAt == false {
		result = empty
	}
	return result
}

func (d *DeliveryOrderValidator) getQueryWithDateValidation(param string, empty string, dateFields []*models.DateInputRequest, ctx *gin.Context) (string, []*models.DateInputRequest) {
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

func (d *DeliveryOrderValidator) getIntQueryWithDefault(param string, empty string, isNotZero bool, ctx *gin.Context) (int, error) {
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
