package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// get user Profile
func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	ownerUser := ps.ByName("name")

	offset := getOffset(r)
	_, err := rt.AutorizeAndCheck(r, w, ownerUser, ctx)
	if err != nil {
		return
	}

	// Api Call

	images, err := rt.db.GetImages(ownerUser, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetImages function", err)
	}

	profile := UserProfile{Username: ownerUser, ProfileImages: images}

	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in JSON encoding in getProfile", err)
	}

}

func (rt *_router) getFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	ownerUser := ps.ByName("name")

	_, err := rt.AutorizeAndCheck(r, w, ownerUser, ctx)
	if err != nil {
		return
	}

	offset := getOffset(r)

	followers, err := rt.db.GetFollowed(ownerUser, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetFollowed function in getFollowed", err)
	}

	list := FollowerList{Followers: followers}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in JSON encoding in getFollowed", err)
	}

}
func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	ownerUser := ps.ByName("name")

	_, err := rt.AutorizeAndCheck(r, w, ownerUser, ctx)
	if err != nil {
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	followers, err := rt.db.GetFollowers(ownerUser, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetFollowers function in getFollowers", err)
	}

	list := FollowerList{Followers: followers}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in JSON encoding in getFollowers", err)
	}
}

func (rt *_router) getBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	yourself := ps.ByName("name")

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}
	if name != yourself {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "The provided token was not valid. Please provide a valid token"})
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	followers, err := rt.db.GetBan(yourself, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetBan function in getBan", err)
	}

	list := BanList{Followers: followers}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in JSON encoding in getBan", err)
	}
}
