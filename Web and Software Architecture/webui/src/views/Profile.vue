<script>
import Feed from '../components/Feed.vue';

export default {
    emits: ["login"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            me: localStorage.getItem("username"),
            profileName: null,
            photos: null,
            offset: 0,
			lastBatch: 0,
			newThings: Array(),
            followers: null,
            following: null,
            isFollowing: null,
            banned: null,
        };
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            this.profileName = this.$route.path.replace("/user/", "");
            if ((this.me == this.profileName)) {
                this.$router.replace("/me");
                return
            }
            else {
                await this.$axios.get(`/user/${this.profileName}`).then(response => {
                if (response.data.images != null) {
                    this.some_data = response.data.images;
                    this.lastBatch = response.data.images.length;
                    this.loading = false;
                }
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                    this.loading = false;
                });

                await this.$axios.get(`/user/${this.profileName}/follow`, ).then(response => {
                    this.following = response.data.followers.length;
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                });

                await this.$axios.get(`/user/${this.profileName}/followers`, ).then(response => {
                    this.followers = response.data.followers.length;
                    for(var i = 0; i < response.data.followers.length; i++) {
                        if (response.data.followers[i] == this.me) {
                            this.isFollowing = true;
                            break
                        }
                    }
                    if (this.isFollowing == null) {
                        this.isFollowing = false;
                    }
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                });

                await this.$axios.get(`/user/${this.me}/ban`, ).then(response => {
                    for(var i = 0; i < response.data.ban.length; i++) {
                        if (response.data.ban[i] == this.profileName) {
                            this.banned = true;
                            break
                        }
                    }
                    if (this.banned == null) {
                        this.banned = false;
                    }
                }).catch(err => {
                    this.errormsg = err.response.data.message;
                });
            }
            
        },
        async MORE() {
            this.errormsg = null
			this.newThings = Array();
            this.offset = this.offset + 1;
            await this.$axios.get(`/user/${this.profileName}?offset=${this.offset}`).then(response => {
				this.newThings = response.data.images;
                this.lastBatch = response.data.images.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async follow() {
            this.errormsg = null;
            await this.$axios.put(`/user/${this.me}/follow/${this.profileName}`).then(response => {
                this.isFollowing = true;
                this.followers = this.followers + 1;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async unfollow() {
            this.errormsg = null;
            await this.$axios.delete(`/user/${this.me}/follow/${this.profileName}`).then(response => {
                this.isFollowing = false;
                this.followers = this.followers - 1;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async ban() {
            this.errormsg = null;
            await this.$axios.put(`/user/${this.me}/ban/${this.profileName}`).then(response => {
                this.banned = true;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
        async unban() {
            this.errormsg = null;
            await this.$axios.delete(`/user/${this.me}/ban/${this.profileName}`).then(response => {
                this.banned = false;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
    },
    mounted() {
        this.refresh();
    },
    watch: {
    '$route': {
        handler() {
            if (!this.$route.path.replace("/user/", "").startsWith('/')) {
                this.refresh();
            }
        },
        }
	},
    components: { Feed }
}
</script>

<template>
	<div class="home-page">
        <button v-if="!banned" type="button" @click="ban" class="btn btn-sm btn-center primary red">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
			<div class="medium-text">Ban</div>
		</button>
        <button v-if="banned" type="button" @click="unban" class="btn btn-sm btn-center primary">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
			<div class="medium-text">Unban</div>
		</button>
		<h1 class="h2 me-2 big-text clickable" @click="refresh">{{ profileName }}</h1>
		<button v-if="!isFollowing" type="button" @click="follow" class="btn btn-sm btn-center primary">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
			<div class="medium-text">Follow</div>
		</button>
        <button v-if="isFollowing" type="button" @click="unfollow" class="btn btn-sm btn-center primary">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
			<div class="medium-text">Unfollow</div>
		</button>
	</div>
	    <div class="feed">
            <div class="follow-section">
            <div class="follow-indicator">
                <div class="small-text">Followers</div>
                <div class="medium-text"> {{ this.followers }}</div>
            </div>
            <div class="follow-indicator">
                <div class="small-text">Following</div>
                <div class="medium-text"> {{ this.following }}</div>
            </div>
        </div>
		<Feed v-if="some_data" :post_data="some_data" :newThings="this.newThings"></Feed>
		<button v-if="lastBatch==30" class="woomy" @click="MORE">More</button>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
	
</template>

<style>
</style>