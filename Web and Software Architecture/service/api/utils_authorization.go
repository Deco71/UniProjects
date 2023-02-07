package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"wasaPhoto/service/api/reqcontext"
)

// Only authorization
func (rt *_router) AutorizeToken(r *http.Request, w http.ResponseWriter) (string, error) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 1 {
		rt.HttpErrCodeSender(w, http.StatusUnauthorized, "The provided token was not valid. Please provide a valid token")
		return "", errors.New("token not found")
	} else {
		reqToken = splitToken[1]
		name, err := rt.db.CheckToken(reqToken)
		if err != nil {
			rt.HttpErrCodeSender(w, http.StatusUnauthorized, "The provided token was not valid. Please provide a valid token")
			return "", errors.New("token not valid")
		}
		return name, err
	}

}

// Authorization and Equality, Existence and Ban checks
func (rt *_router) AutorizeAndCheckFull(r *http.Request, w http.ResponseWriter, ownerUser string, myself string, ctx reqcontext.RequestContext) (string, error) {

	// Authorization
	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return "", err
	}

	// Equality check
	if myself != "" {
		if myself != name {
			rt.HttpErrCodeSender(w, http.StatusUnauthorized, "The provided token was not valid. Please provide a valid token")
			return "", errors.New("unauthorized")
		}
	}

	// Existence check
	err = rt.db.CheckUser(ownerUser)
	if errors.Is(err, sql.ErrNoRows) {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "User not found")
		return "", err
	} else if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in CheckUser function", err)
		return "", err
	}

	// Ban Check
	value, err := rt.db.CheckBan(ownerUser, name)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in CheckBan function", err)
		return "", err
	} else if value == 1 {
		rt.HttpErrCodeSender(w, http.StatusForbidden, "You are unable to access this resource, probably you are banned from this user. Try again later")
		return "", errors.New("banned")
	}
	return name, err
}

// Authorization, Existence and Ban checks
func (rt *_router) AutorizeAndCheck(r *http.Request, w http.ResponseWriter, ownerUser string, ctx reqcontext.RequestContext) (string, error) {
	return rt.AutorizeAndCheckFull(r, w, ownerUser, "", ctx)

}

// Authorization and Equality, Existence and Ban checks
func (rt *_router) AutorizeAndCheckNoBan(r *http.Request, w http.ResponseWriter, ownerUser string, myself string, ctx reqcontext.RequestContext) (string, error) {

	// Authorization
	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return "", err
	}

	// Equality check
	if myself != "" {
		if myself != name {
			rt.HttpErrCodeSender(w, http.StatusUnauthorized, "The provided token was not valid. Please provide a valid token")
			return "", errors.New("unauthorized")
		}
	}

	// Existence check
	err = rt.db.CheckUser(ownerUser)
	if errors.Is(err, sql.ErrNoRows) {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "User not found")
		return "", err
	} else if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in CheckUser function", err)
		return "", err
	}

	return name, err
}
