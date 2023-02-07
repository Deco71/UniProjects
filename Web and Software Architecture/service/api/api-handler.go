package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.postLogin))
	rt.router.PUT("/setMyUsername", rt.wrap(rt.putUsername))
	rt.router.GET("/user/:name/follow", rt.wrap(rt.getFollowed))
	rt.router.GET("/user/:name/followers", rt.wrap(rt.getFollowers))
	rt.router.GET("/user/:name", rt.wrap(rt.getProfile))
	rt.router.POST("/uploadImage", rt.wrap(rt.postPhoto))
	rt.router.DELETE("/image/:img", rt.wrap(rt.deletePhoto))
	rt.router.GET("/image/:img", rt.wrap(rt.getPhoto))
	rt.router.PUT("/user/:me/follow/:followedName", rt.wrap(rt.followUser))
	rt.router.DELETE("/user/:me/follow/:followedName", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/user/:me/ban/:bannedUser", rt.wrap(rt.BanUser))
	rt.router.DELETE("/user/:me/ban/:bannedUser", rt.wrap(rt.unbanUser))
	rt.router.GET("/user/:name/ban", rt.wrap(rt.getBan))
	rt.router.GET("/feed", rt.wrap(rt.getFeed))
	rt.router.GET("/post/:img", rt.wrap(rt.getPost))
	rt.router.PUT("/post/:img/like/:me", rt.wrap(rt.likePost))
	rt.router.DELETE("/post/:img/like/:me", rt.wrap(rt.unlikePost))
	rt.router.POST("/post/:img/comment", rt.wrap(rt.commentPost))
	rt.router.DELETE("/post/:img/comment/:comment", rt.wrap(rt.uncommentPost))
	rt.router.GET("/post/:img/likes", rt.wrap(rt.getLikes))
	rt.router.GET("/post/:img/comments", rt.wrap(rt.getComments))

	// TODO: Autoremove Follow on ban and remove header on 204 AND get likes and get comment
	// ADD GROUP BY ASC/DESC ON GET FEED
	// If you want: Fix the /getComment that permits a banned user to see a comment of a post of which he's banned

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
