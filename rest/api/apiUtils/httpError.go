package apiUtils

import (
	"net/http"

	"go.uber.org/zap"
)

// To be called in all http backend api handlers to log an error
func HttpError(response http.ResponseWriter, logger *zap.SugaredLogger, message string, httpStatusCode int) {
	logger.Error(message)
	http.Error(response, message, httpStatusCode)
}

// To be called in all http backend api handlers to log an error that varies for both the logger and the response
func HttpErrorDuoMessages(response http.ResponseWriter, logger *zap.SugaredLogger, loggerMessage string, responseMessage string, httpStatusCode int) {
	logger.Error(loggerMessage)
	http.Error(response, responseMessage, httpStatusCode)
}

// To be called in http backend api handlers to log an error, together with other information (treated as fields)
func HttpErrorWithFields(response http.ResponseWriter, logger *zap.SugaredLogger, message string, httpStatusCode int, fields []zap.Field) {
	logger.Error(message, fields)
	http.Error(response, message, httpStatusCode)
}

// To be called in http backend api handlers to log an error that varies for both the logger and the response, together with other information (treated as fields)
func HttpErrorDuoMessagesWithFields(response http.ResponseWriter, logger *zap.SugaredLogger, loggerMessage string, responseMessage string, httpStatusCode int, fields []zap.Field) {
	logger.Error(loggerMessage, fields)
	http.Error(response, responseMessage, httpStatusCode)
}
