package utils

import "github.com/gin-gonic/gin"

func SendErrResponse(gContext *gin.Context, httpCode int, errorText string) {
	var errorJson = map[string]string{
		"result":    "fail",
		"errorText": errorText,
	}
	gContext.JSON(httpCode, errorJson)
}

func SendSuccessResponse(gContext *gin.Context, httpCode int) {
	var successJson = map[string]string{
		"result": "success",
	}
	gContext.JSON(httpCode, successJson)
}
