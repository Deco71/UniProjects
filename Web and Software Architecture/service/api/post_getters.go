package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	username, likeCount, commentCount, created_at, liked, err := rt.db.GetPost(post, name)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetPost", err)
		return
	}

	photoInt, _ := strconv.Atoi(post)
	likeInt, _ := strconv.Atoi(likeCount)
	commentInt, _ := strconv.Atoi(commentCount)
	likedBool, _ := strconv.ParseBool(liked)

	_ = json.NewEncoder(w).Encode(Post{Username: username, ImageId: photoInt, LikesValue: likeInt, Liked: likedBool, CommentsValue: commentInt, Date: created_at})
}

func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	offset := getOffset(r)

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}

	list, err := rt.db.GetFeed(name, offset)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in GetFeed", err)
	}

	var posts PostArray

	for id := range list {
		username, likeCount, commentCount, created_at, liked, err := rt.db.GetPost(list[id], name)
		if err != nil {
			rt.ErrLoggerAndSender(w, ctx, "Error in GetPost", err)
			return
		}

		photoInt, _ := strconv.Atoi(list[id])
		likeInt, _ := strconv.Atoi(likeCount)
		commentInt, _ := strconv.Atoi(commentCount)
		likedBool, _ := strconv.ParseBool(liked)

		postSlice := Post{Username: username, ImageId: photoInt, LikesValue: likeInt, Liked: likedBool, CommentsValue: commentInt, Date: created_at}
		posts.Posts = append(posts.Posts, postSlice)
	}

	_ = json.NewEncoder(w).Encode(posts)
}
