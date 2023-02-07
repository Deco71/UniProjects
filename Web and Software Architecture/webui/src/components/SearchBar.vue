<script>
export default {
	props: ['searchStyle'],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            username: '',
        };
    },
    methods: {
        async search() {
            this.loading = true;
            this.errormsg = null;
            await this.$axios.get(`/user/${this.username}`).then(response => {
                this.$router.push(`/user/${this.username}`);
                this.loading = false;
            }).catch(err => {
                this.loading = false;
                if (err = "AxiosError: Request failed with status code 404") {
                    this.errormsg = "User not found";
                }
                else {
                    this.errormsg = err.response.data.message;
                }
            });
            this.username = '';
            
        },
        cleanErr() {
            this.errormsg = null;
        },
    },
}
</script>

<template>
	<div :class="this.searchStyle">
    	<input type="text" v-model="username" class="searchTerm" placeholder="Search users here!">
    	<button type="submit" class="searchButton" @click="search">
    		<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
    	</button>
   	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg" timer=5 @timerEnd="cleanErr"></ErrorMsg>
</template>

<style>
</style>
