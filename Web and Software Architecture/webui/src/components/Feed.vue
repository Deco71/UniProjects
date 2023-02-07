<script>
export default {
    props: ["post_data", "newThings"],
	data: function() {
		return {
			errormsg: null,
			loading: false,
			me: localStorage.getItem("username"),
            profileName: null,
			photos: null,
            some_data: null,
		}
	},
	methods: {
        refresh() {
            this.errormsg = null
            this.some_data = this.post_data
        },
        removePost(id) {
			this.some_data = this.some_data.filter(function(photo) {
    			return photo != id;
			})
		}
	},
    watch: {
    post_data: {
        handler(newVal, oldVal) {
			if (oldVal !== undefined) {
            	this.refresh();
			}
        },
        },
	newThings: {
		handler(newVal, oldVal) {
			if (newVal != oldVal && this.some_data != null) {
				this.some_data = this.some_data.concat(this.newThings)
			}
		},
    	},
	},
	mounted() {
        this.refresh();
    },
}
</script>

<template>
	<div class="feed">
		<div v-if=some_data>
			<div v-for="images in this.some_data" :key="images" class="post">
				<Post :post="images" @removed="removePost"></Post>
			</div>
		</div>
		<div v-if="!some_data" class="flex">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#loader"/></svg>
			<div class="medium-text">Loading...</div>
		</div>
		<div v-if="some_data&&some_data.length==0" class="flex">
			<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#alert-triangle"/></svg>
			<div class="medium-text">No posts yet</div>
		</div>
	</div>
	
</template>

<style>
</style>