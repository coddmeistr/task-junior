package sort

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxik12233/task-junior/pkg/api"
)

const (
	ASC               = "asc"
	DESC              = "desc"
	OptionsContextKey = "sort_options"
)

type Options struct {
	Field string
	Order string
}

func Middleware(defaultSortField, defaultSortOrder string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sortBy := c.Request.URL.Query().Get("sort_by")
		sortOrder := c.Request.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = defaultSortField
		}

		if sortOrder == "" {
			sortOrder = defaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != ASC && upperSortOrder != DESC {
				c.Writer.WriteHeader(http.StatusBadRequest)
				err := api.ErrorResponse{
					Message: "collation must be asc or desc",
					Details: nil,
				}
				c.Writer.Write(err.Marshal())
				return
			}
		}

		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}
		ctx := context.WithValue(c.Request.Context(), OptionsContextKey, options)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
