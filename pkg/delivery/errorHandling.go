package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleOK(c *gin.Context, data interface{}) {

	resp := map[string]interface{}{
		"message": "OK",
		"data":    data,
	}

	c.JSON(http.StatusOK, resp)
}

func HandleError(c *gin.Context, statusCode int) {

	errMsg := ""

	switch statusCode {
	case http.StatusBadRequest:
		errMsg = "Bad Request"
	case http.StatusUnauthorized:
		errMsg = "Unauthorized"
	case http.StatusForbidden:
		errMsg = "Forbidden"
	case http.StatusNotFound:
		errMsg = "Not Found"
	case http.StatusInternalServerError:
		errMsg = "Internal Server Error"
	default:
		errMsg = "Unknown Error"
	}

	resp := map[string]interface{}{
		"message": errMsg,
	}

	c.JSON(statusCode, resp)
}
