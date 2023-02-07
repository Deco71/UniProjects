<script>
export default {
    props: ["post", "liked"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            id: null,
            commentsList: new Array(),
            userlink: "/user/",
            likeInfo: null,
            me: localStorage.getItem("username"),
            likeStatus: false,
            offset: 0,
            lastBatch: 0,
        };
    },
    methods: {
        async refresh() {
            this.id = this.post
            this.errormsg = null
            await this.$axios.get(`/post/${this.id}/likes`).then(response => {
                this.likeInfo = response.data.likers;
                this.lastBatch = response.data.likers.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
            if (this.liked) {
                this.likeStatus = true
            }
            else{
                this.likeStatus = false
            }
        },
        async MORE() {
            this.errormsg = null
            this.offset = this.offset + 1;
            await this.$axios.get(`/post/${this.id}/likes?offset=${this.offset}`).then(response => {
                this.likeInfo = this.likeInfo.concat(response.data.likers);
                this.lastBatch = response.data.comments.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
    },
    mounted() {
        this.refresh();
    },
    watch: { 
      	liked: function(newVal, oldVal) {
            this.likeStatus = newVal;
        },
        post: function(newVal, oldVal) {
            this.refresh();
        },
    },
}
</script>
<template>
    <div class="like-header">
        <div v-if=this.likeStatus class="small-text">
            <RouterLink to="/me">
                    {{ me }}
            </RouterLink>
        </div>
        <div v-if=this.likeInfo>
            <div v-for="liker in this.likeInfo" :key="liker">
                <div v-if="liker!=me" class="small-text">
                    <RouterLink :to="this.userlink+liker">
                        {{ liker }}
                    </RouterLink>
                </div>
            </div>
        </div>
        <button v-if="lastBatch==30" class="woomy" @click="MORE">More</button>
    </div>
</template>
