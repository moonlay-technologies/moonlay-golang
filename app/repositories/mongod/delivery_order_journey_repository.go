package repositories

import (
	"context"
	"order-service/app/models"
)

type DeliveryOrderJourneysRepositoryInterface interface {
	Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan)
}
