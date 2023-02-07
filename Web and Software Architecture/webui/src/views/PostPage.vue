<script>
export default {
    emits: ['login'],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            posted: null,
			id: null,
            string: null,
            comments: true,
            CommentsList: null,
            liked: false,
            commentChanged: null,
            ItisWhatItis: 0,
            reMounted: true,
        };
    },
    methods: {
        async refresh() {
            this.string = this.$route.path.replace('/post/', '')
            this.errormsg = null
            if (this.$route.query.mode == "comments") {
                this.comments = true
            } else {
                this.comments = false
            }
            try{
                this.id = this.string.split("/")[0]
            } catch(err) {
                if (this.id == null) {
                    this.errormsg = err
                }
            }
        },
        changeToComments() {
            this.$router.replace({ path: `/post/${this.id}/`, query: { mode: 'comments' }})
        },
        changeToLikes() {
            this.$router.replace({ path: `/post/${this.id}/`, query: { mode: 'likes' }})
        },
        goBack() {
            if (this.reMounted) {
                this.$router.go(-2)
            }
            else{
                this.$router.go(-1)
            }
        },
        changeLiked(status) {
            this.liked = status
        },
        CommentChange(status) {
            this.commentChanged = [status, this.ItisWhatItis]
            this.ItisWhatItis = this.ItisWhatItis + 1
        },
    },
    watch: {
    '$route': {
        handler() {
            this.reMounted = !this.reMounted
            this.refresh();
        },
        }
	},
    mounted() {
        this.refresh();
    },
}
</script>

<template>
    <div v-if="this.id!==null" class="full-post">
        <Post :post="this.id" :removeActive="true" :commentChange="commentChanged" @removed="goBack" @likeStatus="changeLiked"></Post>
        <div class="interaction-container">
            <div class="interaction-header">
                <button class="small-text woomy" @click="changeToLikes">See Likes</button>
                <button class="small-text woomy" @click="changeToComments">See Comments</button>
            </div>
            <div v-if="this.comments" class="comment-section">
                <Comments :post="this.id" @commentsChanged="CommentChange"></Comments>
            </div>
            <div v-else class="like-section">
                <Likes :post="this.id" :liked="this.liked"></Likes>
            </div>
        </div>
    </div>
    
</template>

<style>
</style>
