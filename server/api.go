package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/jholdstock/dcrwages/model"
)

type errmsg struct {
	Err string `json:"error"`
}

// apiReady is a HTTP handler for the API routes. It returns a 503
// HTTP status if the data model is not yet loaded, and then prevents
// further HTTP handlers from executing.
func apiReady() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !model.Initialised {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable,
				errmsg{
					"dcrwages is initialising",
				})
		}
	}
}

func getIntParam(key string, c *gin.Context) (int, error) {
	value, err := strconv.Atoi(c.Param(key))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			errmsg{
				fmt.Sprintf("Could not decode %s param. Should be an integer.", key),
			},
		)
		return -1, err
	}
	return value, nil
}

func jsonError(httpStatus int, err string, c *gin.Context) {
	c.JSON(httpStatus,
		errmsg{
			err,
		})
}

// Return all available data, json encoded
func getAllData(c *gin.Context) {
	c.JSON(http.StatusOK,
		model.FullHistory)
}

// Return a single year, json encoded
func getYear(c *gin.Context) {
	yearParam, err := getIntParam("year", c)
	if err != nil {
		return
	}

	year, err := model.FullHistory.FindYear(yearParam)
	if err != nil {
		jsonError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK,
		year)
}

// Return a single month, json encoded
func getMonth(c *gin.Context) {
	yearParam, err := getIntParam("year", c)
	if err != nil {
		return
	}

	monthParam, err := getIntParam("month", c)
	if err != nil {
		return
	}

	month, err := model.FullHistory.FindMonth(yearParam, monthParam)
	if err != nil {
		jsonError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, month)
}
