package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	post := ps.ByName("img")

	ownerUser, _ := rt.db.CheckImage(post)
	if ownerUser == "" {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Image Not Found")
		return
	}

	_, err := rt.AutorizeAndCheckFull(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}

	err = rt.db.Like(post, myself)

	if err != nil && err.Error() != "UNIQUE constraint failed: Like.user, Like.photo" {
		rt.ErrLoggerAndSender(w, ctx, "Error in Like", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unlikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")
	myself := ps.ByName("me")
	post := ps.ByName("img")

	ownerUser, _ := rt.db.CheckImage(post)
	if ownerUser == "" {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Image Not Found")
		return
	}
	_, err := rt.AutorizeAndCheckNoBan(r, w, ownerUser, myself, ctx)
	if err != nil {
		return
	}

	err = rt.db.Unlike(post, myself)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in Unlike", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) commentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	post := ps.ByName("img")

	postOwner, err := rt.db.CheckImage(post)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Post not found")
		return
	}

	name, err := rt.AutorizeAndCheck(r, w, postOwner, ctx)
	if err != nil {
		return
	}
	var comment CommentText

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusBadRequest, "Received an invalid json")
		return
	}

	id, createdAt, err := rt.db.Comment(post, name, comment.Comment)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in Comment", err)
		return
	}

	idInt, _ := strconv.Atoi(id)

	_ = json.NewEncoder(w).Encode(Comment{User: name, Comment: comment.Comment, CommentId: idInt, Date: createdAt})

}

func (rt *_router) uncommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	post := ps.ByName("img")
	commentId := ps.ByName("comment")

	_, post2, owner, _, _, err := rt.db.GetComment(commentId)
	if err != nil || post2 != post {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Comment or post not found")
		return
	}

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}

	if name != owner {
		rt.HttpErrCodeSender(w, http.StatusUnauthorized, "The provided token was not valid. Please provide a valid token")
		return
	}

	err = rt.db.Uncomment(commentId)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in Uncomment", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
