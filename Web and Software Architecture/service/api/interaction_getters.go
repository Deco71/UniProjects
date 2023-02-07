package api

import (
	"encoding/json"
	"net/http"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	post := ps.ByName("img")
	offset := getOffset(r)

	postOwner, err := rt.db.CheckImage(post)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Post not found")
		return
	}

	name, err := rt.AutorizeAndCheck(r, w, postOwner, ctx)
	if err != nil {
		return
	}

	comments, err := rt.getCommentsArray(post, name, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetComments", err)
		return
	}

	_ = json.NewEncoder(w).Encode(comments)

}

func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	post := ps.ByName("img")
	offset := getOffset(r)

	postOwner, err := rt.db.CheckImage(post)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Post not found")
		return
	}

	name, err := rt.AutorizeAndCheck(r, w, postOwner, ctx)
	if err != nil {
		return
	}

	users, err := rt.db.GetLikes(post, name, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetLikes", err)
		return
	}

	_ = json.NewEncoder(w).Encode(Likers{Likers: users})
}
