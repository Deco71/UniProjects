<script>

export default {
	emits: ['login'],
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			me: localStorage.getItem("username"),
			offset: 0,
			lastBatch: 0,
			newThings: Array(),
			followers: null,
			following: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.me = localStorage.getItem("username");
			await this.$axios.get(`/user/${this.me}`).then(response => {
                if (response.data.images != null) {
					this.some_data = response.data.images
					this.lastBatch = response.data.images.length;
					this.loading = false;
				}
            }).catch(err => {
                this.errormsg = err.response.data.message;
				this.loading = false;
            });

			await this.$axios.get(`/user/${this.me}/follow`, ).then(response => {
                this.following = response.data.followers.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });

            await this.$axios.get(`/user/${this.me}/followers`, ).then(response => {
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
		},
		async MORE() {
            this.errormsg = null
			this.newThings = Array();
            this.offset = this.offset + 1;
            await this.$axios.get(`/user/${this.me}?offset=${this.offset}`).then(response => {
				this.newThings = response.data.images;
                this.lastBatch = response.data.images.length;
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div class="home-page">
		<RouterLink to="/settings">
			<button type="button" class="btn btn-sm btn-center border border-secondary" @click="newItem">
				<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
				<div class="medium-text">Settings</div>
			</button>
		</RouterLink>
		<h1 class="h2 me-2 big-text clickable" @click="refresh">{{ this.me }}</h1>
		<RouterLink to="/upload">
		<button type="button" class="btn btn-sm btn-center primary" @click="newItem">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
			<div class="medium-text">New Post</div>
		</button>
		</RouterLink>	
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
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>