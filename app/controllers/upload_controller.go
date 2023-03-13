package controllers

import (
	"context"
	"order-service/app/usecases"

	"github.com/gin-gonic/gin"
)

type UploadControllerInterface interface {
	UploadSOSJ(ctx *gin.Context)
	UploadDO(ctx *gin.Context)
	UploadSO(ctx *gin.Context)
}

type uploadController struct {
	uploadUseCase usecases.UploadUseCaseInterface
	ctx           context.Context
}

func InitUploadController(uploadUseCase usecases.UploadUseCaseInterface, ctx context.Context) UploadControllerInterface {
	return &uploadController{
		uploadUseCase: uploadUseCase,
		ctx:           ctx,
	}
}

func (c *uploadController) UploadSOSJ(ctx *gin.Context) {

	c.uploadUseCase.UploadSOSJ(ctx)

}

func (c *uploadController) UploadDO(ctx *gin.Context) {

	c.uploadUseCase.UploadDO(ctx)

}

func (c *uploadController) UploadSO(ctx *gin.Context) {

	c.uploadUseCase.UploadSO(ctx)

}
