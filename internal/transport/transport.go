package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/maxik12233/task-junior/internal/service"
	"go.uber.org/zap"
)

type ITransport interface {
	RegisterRoutes(router *gin.Engine)
}

//

const (
	entityURL = "stats"
)

type Transport struct {
	svc    service.IService
	logger *zap.Logger
}

func NewTransport(svc service.IService, logger *zap.Logger) ITransport {
	return &Transport{
		svc:    svc,
		logger: logger,
	}
}

func (t *Transport) RegisterRoutes(router *gin.Engine) {

	public := router.Group("")

	statisticEntity := public.Group(entityURL)
	statisticEntity.GET("", t.AddNameStats)

}

func (t *Transport) AddNameStats(c *gin.Context) {
	// TODO
}
