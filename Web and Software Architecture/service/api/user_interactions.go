package api

import (
	"encoding/json"
	"net/http"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	ownerUser := ps.ByName("followedName")
	if myself == ownerUser {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "You can't execute this action on yourself!"})
		return
	}
	_, err := rt.AutorizeAndCheckFull(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}
	err = rt.db.Follow(myself, ownerUser)

	if err != nil && err.Error() != "UNIQUE constraint failed: Follow.user, Follow.followed" {
		rt.baseLogger.WithError(err).Error("Error in Follow")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	ownerUser := ps.ByName("followedName")
	_, err := rt.AutorizeAndCheckNoBan(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}

	err = rt.db.Unfollow(myself, ownerUser)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error in Unfollow")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	ownerUser := ps.ByName("bannedUser")
	if myself == ownerUser {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "You can't execute this action on yourself!"})
		return
	}
	_, err := rt.AutorizeAndCheckNoBan(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}

	err = rt.db.Ban(myself, ownerUser)
	if err != nil && err.Error() != "UNIQUE constraint failed: Ban.user, Ban.banned" {
		rt.baseLogger.WithError(err).Error("Error in Follow")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TODO: Automatic follow removal
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	ownerUser := ps.ByName("bannedUser")
	_, err := rt.AutorizeAndCheckNoBan(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}

	err = rt.db.Unban(myself, ownerUser)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error in Unfollow")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
