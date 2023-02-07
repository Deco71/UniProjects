package api

import "strconv"

type UserProfile struct {
	Username      string   `json:"username"`
	ProfileImages []string `json:"images"`
}

type Image struct {
	Image int `json:"imageId"`
}

type FollowerList struct {
	Followers []string `json:"followers"`
}

type BanList struct {
	Followers []string `json:"ban"`
}

type BannedList struct {
	Followers []string `json:"ban"`
}

type UserLogin struct {
	Username string `json:"username"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
}

type CommentText struct {
	Comment string `json:"comment"`
}

type Comment struct {
	User      string `json:"username"`
	Comment   string `json:"comment"`
	CommentId int    `json:"commentId"`
	Date      string `json:"date"`
}

type CommentArray struct {
	Comments []Comment `json:"comments"`
}

type JSONErrorMsg struct {
	Message string `json:"message"`
}

type Likers struct {
	Likers []string `json:"likers"`
}

type PostArray struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Username      string `json:"username"`
	ImageId       int    `json:"imageId"`
	CommentsValue int    `json:"commentsValue"`
	LikesValue    int    `json:"likesValue"`
	Liked         bool   `json:"liked"`
	Date          string `json:"date"`
}

func (rt *_router) getCommentsArray(post string, name string, offset int) (CommentArray, error) {

	commentArray := make([]Comment, 0)

	ids, users, comments, dates, err := rt.db.GetComments(post, name, offset)

	for (err == nil) && (len(ids) > 0) {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			return CommentArray{}, err
		}
		commentArray = append(commentArray, Comment{
			CommentId: id,
			User:      users[0],
			Comment:   comments[0],
			Date:      dates[0],
		})
		ids = ids[1:]
		users = users[1:]
		comments = comments[1:]
		dates = dates[1:]
	}

	return CommentArray{Comments: commentArray}, err

}
