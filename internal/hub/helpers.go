package hub

import "net/http"

func UnauthorizedError(responseMsg *WSResponseMessage) {
	responseMsg.Type = "error"
	responseMsg.StatusCode = http.StatusUnauthorized
	responseMsg.Message = http.StatusText(http.StatusUnauthorized)
}

func ServerError(responseMsg *WSResponseMessage) {
	responseMsg.Type = "error"
	responseMsg.StatusCode = http.StatusInternalServerError
	responseMsg.Message = http.StatusText(http.StatusInternalServerError)
}

func Success(responseMsg *WSResponseMessage) {
	responseMsg.Type = "success"
	responseMsg.StatusCode = http.StatusOK
}
