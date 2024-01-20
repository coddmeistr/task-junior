package transport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	app "github.com/maxik12233/task-junior"
	"github.com/maxik12233/task-junior/internal/service"
	"github.com/maxik12233/task-junior/pkg/api/paginate"
	"github.com/maxik12233/task-junior/pkg/api/sort"
	"go.uber.org/zap"
)

type ITransport interface {
	RegisterRoutes(router *gin.Engine)
}

//

const (
	entityURL = "person"
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
		t.logger.Error("Error given bad json body", zap.Error(err))
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		t.logger.Error("Failed struct validation", zap.Error(err))
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body. Failed validation").Error())
		return
	}

	info, err := t.svc.CreatePersonInfo(c.Request.Context(), req.ToDomain())
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
		t.logger.Error("Error given bad json body", zap.Error(err))
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		t.logger.Error("Failed struct validation", zap.Error(err))
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
	var sortOptions *sort.Options
	if options, ok := c.Request.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = &options
	}

	var paginateOptions *paginate.Options
	if options, ok := c.Request.Context().Value(paginate.OptionsContextKey).(paginate.Options); ok {
		paginateOptions = &options
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err == nil {
		person, err := t.svc.GetPersonInfo(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(app.GetHTTPCodeFromError(err), err.Error())
			return
		}

		c.JSON(http.StatusOK, GetPersonInfoResponse{
			PersonResponse: PersonResponse{
				Id:          person.ID,
				Name:        person.Name,
				Surname:     person.Surname,
				Patronymic:  person.Patronymic,
				Age:         person.Characteristic.Age,
				Gender:      person.Characteristic.Gender,
				Nationality: person.Characteristic.Nationality,
			},
		})
	} else {

		persons, err := t.svc.GetAllPersonInfo(c.Request.Context(), sortOptions, paginateOptions)
		if err != nil {
			c.JSON(app.GetHTTPCodeFromError(err), err.Error())
			return
		}

		count, err := t.svc.GetPersonCount(c.Request.Context())
		if err != nil {
			c.JSON(app.GetHTTPCodeFromError(err), err.Error())
			return
		}

		personResponses := make([]PersonResponse, len(persons))
		for i, v := range persons {
			personResponses[i] = PersonResponse{
				Id:          v.ID,
				Name:        v.Name,
				Surname:     v.Surname,
				Patronymic:  v.Patronymic,
				Age:         v.Characteristic.Age,
				Gender:      v.Characteristic.Gender,
				Nationality: v.Characteristic.Nationality,
			}
		}

		c.JSON(http.StatusOK, GetPersonInfoResponse{
			TotalCount: &count,
			Persons:    personResponses,
		})
	}
}

func (t *Transport) UpdatePersonInfo(c *gin.Context) {
	var req UpdatePersonInfoRequest
	if err := c.BindJSON(&req); err != nil {
		t.logger.Error("Error given bad json body", zap.Error(err))
		c.JSON(app.GetHTTPCodeFromError(app.ErrBadRequest), app.WrapE(app.ErrBadRequest, "Bad JSON body").Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		t.logger.Error("Failed struct validation", zap.Error(err))
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
