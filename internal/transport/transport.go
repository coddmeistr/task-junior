package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	app "github.com/maxik12233/task-junior"
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
	statisticEntity.GET("", t.GetPersonInfo)
	statisticEntity.POST("", t.AddPersonInfo)
	statisticEntity.DELETE("", t.DeletePersonInfo)
	statisticEntity.PUT("", t.UpdatePersonInfo)
}

func (t *Transport) AddPersonInfo(c *gin.Context) {
	var req AddPersonInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body. Failed validation").Error())
		return
	}

	info, err := t.svc.AddNewNameInfo(c.Request.Context(), req.ToDomain())
	if err != nil {
		c.JSON(app.GetHTTPCodeFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, AddPersonInfoResponse{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         info.Age,
		Gender:      info.Gender,
		Nationality: info.Nationality,
	})
}

func (t *Transport) DeletePersonInfo(c *gin.Context) {
	var req DeletePersonInfoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body. Failed validation").Error())
		return
	}

	err := t.svc.DeletePersonInfo(c.Request.Context(), int(req.Id))
	if err != nil {
		c.JSON(app.GetHTTPCodeFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, "Entity was deleted")
}

func (t *Transport) GetPersonInfo(c *gin.Context) {
	// TODO
}

func (t *Transport) UpdatePersonInfo(c *gin.Context) {
	var req UpdatePersonInfoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body. Failed validation").Error())
		return
	}

	person, char := req.ToDomain()
	err := t.svc.UpdatePersonInfo(c.Request.Context(), person, char)
	if err != nil {
		c.JSON(app.GetHTTPCodeFromError(err), err.Error())
		return
	}

	c.JSON(http.StatusOK, "Entity was updated")
}
