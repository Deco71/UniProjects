package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")

	var user UserLogin

	var token int

	regex, _ := regexp.Compile("^[A-Za-z0-9]*$")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Received an invalid json")
		return
	}
	if CheckStringLenght(user.Username) != nil {

		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Username must be between 3 and 16 characters")
		return
	}
	if !regex.MatchString(user.Username) {
		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Username must contain only letters and numbers")
		return
	}

	token, err = rt.db.GetToken(user.Username)
	if err != nil {

		err = rt.db.SetName(user.Username)
		if err != nil {
			rt.ErrLoggerAndSender(w, ctx, "Error in SetName function in postLogin", err)
			return
		}
		token, _ = rt.db.GetToken(user.Username)
	}

	tokenObj := Identifier{Identifier: strconv.Itoa(token)}

	err = json.NewEncoder(w).Encode(tokenObj)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in JSON encoding in postLogin", err)
	}

}

func (rt *_router) putUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}

	var user UserLogin

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Received an invalid json")
		return
	}
	if CheckStringLenght(user.Username) != nil {
		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Username must be between 3 and 16 characters")
		return
	}

	err = rt.db.ChangeName(name, user.Username)
	if err != nil && err.Error() == "UNIQUE constraint failed: User.username" {
		rt.HttpErrCodeSender(w, http.StatusConflict, "Username already taken")
		return
	} else if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in ChangeName function", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
