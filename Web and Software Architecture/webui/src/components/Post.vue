<script>
export default {
    props: ["post", "removeActive", "commentChange"],
    emits: ["removed", "likeStatus"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            id: null,
            postInfo: null,
            postLink: "/post/",
			like: "?mode=like",
			comment: "?mode=comments",
            userlink: "/user/",
            postAuthority: null,
        };
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.id = this.post
            this.errormsg = null
            await this.$axios.get(`/post/${this.id}`, {
				headers: { Authorization: `Bearer ${localStorage.getItem("id")}`}
			    }).then(response => {
                    this.postInfo = response.data;
                    this.loading = false;
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                    this.loading = false;
                });
            if (this.postInfo.username === localStorage.getItem("username")) {
                this.postAuthority = true
            }
            else {
                this.postAuthority = false
            }
            if (this.postInfo.liked)
            {
                this.$emit("likeStatus", true)
            }
            else {
                this.$emit("likeStatus", false)
            }

        },
        async remove() {
            if (confirm('Are you sure you want to remove this post?\nThis action cannot be undone.')) {
                await this.$axios.delete(`/image/${this.postInfo.imageId}`).then(response => {
                    this.$emit("removed", this.postInfo.imageId)
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                });
            }
            
        },
        async addLike() {
            await this.$axios.put(`/post/${this.postInfo.imageId}/like/${localStorage.getItem("username")}`).then(response => {
                this.postInfo.liked = true
                this.postInfo.likesValue += 1
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
            this.$emit("likeStatus", true)
        },
        async removeLike() {
            await this.$axios.delete(`/post/${this.postInfo.imageId}/like/${localStorage.getItem("username")}`).then(response => {
                this.postInfo.liked = false
                this.postInfo.likesValue -= 1
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
            this.$emit("likeStatus", false)
        }
    },
    watch: {
        post: {
            handler() {
                this.refresh();
            },
            },
        commentChange: {
            handler() {
                if (this.commentChange === undefined || this.commentChange === null) {
                    return;
                }
                if (this.commentChange[0]) {
                    this.postInfo.commentsValue += 1;
                }
                else {
                    this.postInfo.commentsValue -= 1;
                }
            },

        }
    },
    mounted() {
        this.refresh();
    }
}
</script>
<template>

    <LoadingSpinner :load="this.loading"></LoadingSpinner>

    <div v-if="postInfo" class="post no-border no-margin">
        <div class="post-header">
            <div class="post-user-name medium-text">
                <RouterLink :to="this.userlink+this.postInfo.username">
                    {{ this.postInfo.username }}
                </RouterLink>
            </div>
            <div class="small-text hour">{{ this.postInfo.date }}</div>
            <button v-if="postAuthority && removeActive" class="small-text woomy delete" @click="remove">Remove Post</button>
        </div>
        <div class="post-photo">
            <div class="post-user-avatar">
                <Immagine v-bind:photoId="postInfo.imageId"></Immagine>
            </div>
        </div>
        <div class="post-options">
			<div class="post-buttons me-2">
				<svg v-if="!this.postInfo.liked" class="feather icon" @click="addLike(this.postInfo.imageId)"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
				<svg v-if="this.postInfo.liked" class="feather icon liked" @click="removeLike(this.postInfo.imageId)"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
                <RouterLink :to=postLink+postInfo.imageId+this.like class="small-text">{{ postInfo.likesValue }}</RouterLink>
			</div>
			<div class="post-buttons">
				<RouterLink :to=postLink+postInfo.imageId+this.comment class="commentViewer">
                    <svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
                    <div class="small-text">
                        {{ postInfo.commentsValue }}
                    </div>
                </RouterLink>
			</div>
		</div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>
