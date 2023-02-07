<script>

export default {
	emits: ['login'],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            photos: null,
			pop: null,
			offset: 0,
			lastBatch: 0,
        };
    },
    methods: {
        async refresh() {
			this.loading = true;
            await this.$axios.get("/feed").then(response => {
                if (response.data.posts != null) {
					this.photos = new Array()
					for (var i = 0; i < response.data.posts.length; i++) {
						this.photos.push(response.data.posts[i].imageId)
					}
					this.lastBatch = response.data.posts.length;
				} else {
					this.photos = Array()
					this.lastBatch = 0;
				}
				this.loading = false;
            }).catch(err => {
                this.errormsg = err.response.data.message;
				this.loading = false;
            });
        },
		async MORE() {
            this.errormsg = null
            this.offset = this.offset + 1;
            await this.$axios.get(`/feed?offset=${this.offset}`).then(response => {
				for (var i = 0; i < response.data.posts.length; i++) {
						this.photos.push(response.data.posts[i].imageId)
				}
			if (response.data.posts != null) {
                this.lastBatch = response.data.posts.length;
			} else {
				this.lastBatch = 0;
			}
            }).catch(err => {
                this.errormsg = err.response.data.message;
            });
        },
    },
    mounted() {
        this.refresh();
    },
}
</script>

<template>
	<div>
		<div class="home-page">
			<button type="button" class="btn btn-sm btn-outline-secondary btn-center" @click="refresh">
				<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#repeat"/></svg>
				<div class="medium-text">Refresh</div>
			</button>
			<h1 class="h2 me-2 big-text">Home</h1>
			<RouterLink to="/upload">
				<button type="button" class="btn btn-sm btn-outline-primary btn-center primary" @click="newItem">
					<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
					<div class="medium-text">New Post</div>
				</button>
			</RouterLink>
		</div>
	</div>
	<div>
		<div class="feed">
			<LoadingSpinner v-if="loading"></LoadingSpinner>
			<Feed v-if="this.photos" :post_data="this.photos"></Feed>
			<button v-if="lastBatch==30" class="woomy" @click="MORE">More</button>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>

<style>
</style>
