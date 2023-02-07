package api

import (
	"encoding/json"
	"net/http"
	"wasaPhoto/service/api/reqcontext"
)

func (rt *_router) ErrLoggerAndSender(w http.ResponseWriter, ctx reqcontext.RequestContext, logString string, err error) {
	ctx.Logger.Error(logString, err)
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal Server Error"})
}

func (rt *_router) HttpErrCodeSender(w http.ResponseWriter, httpCode int, msg string) {
	w.WriteHeader(httpCode)
	_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: msg})
}
