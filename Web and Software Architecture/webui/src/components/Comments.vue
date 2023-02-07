<script>
export default {
    emits: ['commentsChanged'],
    props: ["post"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            id: null,
            commentsList: new Array(),
            userlink: "/user/",
            postInfo: null,
            me: localStorage.getItem("username"),
            newComment: '',
            emptyString: '',
            offset: 0,
            lastBatch: 0,
        };
    },
    methods: {
        async refresh() {
            this.id = this.post
            this.errormsg = null
            await this.$axios.get(`/post/${this.id}/comments`).then(response => {
                this.postInfo = response.data;
                this.lastBatch = response.data.comments.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async postComment() {
            this.errormsg = null
            await this.$axios.post(`/post/${this.id}/comment`, {
                comment: this.newComment,
            }).then(response => {
                this.newComment = '';
                this.$emit('commentsChanged', true);
                this.refresh();
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async removeComment(commentId) {
            this.errormsg = null
            await this.$axios.delete(`/post/${this.id}/comment/${commentId}`).then(response => {
                this.$emit('commentsChanged', false);
                this.refresh();
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
            this.refresh();
        },
        async MORE() {
            this.errormsg = null
            this.offset = this.offset + 1;
            await this.$axios.get(`/post/${this.id}/comments?offset=${this.offset}`).then(response => {
                this.postInfo.comments = this.postInfo.comments.concat(response.data.comments);
                this.lastBatch = response.data.comments.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
    },
    watch: {
    post: {
        handler() {
            this.refresh();
        },
        },
    },
    mounted() {
        this.refresh();
    },
}
</script>
<template>
    <textarea v-model="this.newComment" class="new-comment" placeholder="Write a comment..."></textarea>
    <button class="woomy" @click="postComment" :disabled="newComment===emptyString">Post</button>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <div v-if=this.postInfo class="something-list">
        <div v-for="comment in this.postInfo.comments" v-bind:key="comment" class="comment">
            <div class="comment-header">
                <div class="comment-user-name small-text">
                    <RouterLink :to="this.userlink+comment.username">
                        {{ comment.username }}
                    </RouterLink>
                </div>
                <button v-if="me==comment.username" class="small-text woomy-comment delete" @click="removeComment(comment.commentId)">Remove</button>
                <div class="small-text hour-nopadding"> {{ comment.date }}</div>
            </div>
            <div class="comment-text">
                <div class="small-text"> {{ comment.comment }}</div>
            </div>
        </div>
        <button v-if="lastBatch==30" class="woomy" @click="MORE">More</button>
    </div>
</template>
