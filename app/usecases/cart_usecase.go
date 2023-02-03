package usecases

import (
	"context"
	"database/sql"
	"poc-order-service/app/models"
	"poc-order-service/app/repositories"
	"poc-order-service/global/utils/model"
	"time"

	"github.com/bxcodec/dbresolver"
)

type CartUseCaseInterface interface {
	Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.Cart, *model.ErrorLog)
}

type cartUseCase struct {
	cartRepository        repositories.CartRepositoryInterface
	cartDetailRepository  repositories.CartDetailRepositoryInterface
	orderStatusRepository repositories.OrderStatusRepositoryInterface
	db                    dbresolver.DB
	ctx                   context.Context
}

func InitCartUseCaseInterface(cartRepository repositories.CartRepositoryInterface, cartDetailRepository repositories.CartDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, db dbresolver.DB, ctx context.Context) CartUseCaseInterface {
	return &cartUseCase{
		cartRepository:        cartRepository,
		cartDetailRepository:  cartDetailRepository,
		orderStatusRepository: orderStatusRepository,
		db:                    db,
		ctx:                   ctx,
	}
}

func (u *cartUseCase) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.Cart, *model.ErrorLog) {
	now := time.Now()
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "cart", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.Cart{}, getOrderStatusResult.ErrorLog
	}

	cart := &models.Cart{
		AgentID:         request.AgentID,
		BrandID:         request.BrandID,
		VisitationID:    request.VisitationID,
		UserID:          request.UserID,
		StoreID:         request.StoreID,
		OrderStatusID:   getOrderStatusResult.OrderStatus.ID,
		OrderSourceID:   request.OrderSourceID,
		TotalAmount:     request.TotalAmount,
		TotalTonase:     request.TotalTonase,
		Note:            request.Note,
		CreatedBy:       request.UserID,
		LatestUpdatedBy: request.UserID,
		CreatedAt:       &now,
	}

	insertCartResultChan := make(chan *models.CartChan)
	go u.cartRepository.Insert(cart, sqlTransaction, ctx, insertCartResultChan)
	insertCartResult := <-insertCartResultChan

	if insertCartResult.Error != nil {
		return &models.Cart{}, insertCartResult.ErrorLog
	}

	getCartDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "cart_detail", false, ctx, getCartDetailStatusResultChan)
	getCartDetailStatusResult := <-getCartDetailStatusResultChan

	if getCartDetailStatusResult.Error != nil {
		return &models.Cart{}, getCartDetailStatusResult.ErrorLog
	}

	for _, v := range request.SalesOrderDetails {
		cartDetail := &models.CartDetail{
			CartID:        int(insertCartResult.ID),
			ProductID:     v.ProductID,
			UomID:         v.UomID,
			Qty:           v.Qty,
			Price:         v.Price,
			OrderStatusID: getCartDetailStatusResult.OrderStatus.ID,
			CreatedAt:     &now,
		}

		insertCartDetailResultChan := make(chan *models.CartDetailChan)
		go u.cartDetailRepository.Insert(cartDetail, sqlTransaction, ctx, insertCartDetailResultChan)
		insertCartDetailResult := <-insertCartDetailResultChan

		if insertCartDetailResult.Error != nil {
			return &models.Cart{}, insertCartDetailResult.ErrorLog
		}
	}

	return &models.Cart{}, nil
}
