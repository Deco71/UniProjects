<script>
import LoadingSpinner from '../components/LoadingSpinner.vue';

export default {
    emits: ["login"],
    data: function () {
        return {
            errormsg_username: null,
            errormsg_bans: null,
            newUsername: "",
            loading: false,
            some_data: null,
        };
    },
    methods: {
        async refresh() {
            this.errormsg_bans = null;
            this.errormsg_username = null;
        },
        async changeUsername() {
            this.loading = true;
            this.errormsg_username = null;
            var newname = this.newUsername;
            this.$axios.put("/setMyUsername", {
                username: this.newUsername,
            }).then((response) => {
                localStorage.setItem("username", newname);
                this.loading = false;
                this.$router.go(-1);
            }).catch((error) => {
                this.errormsg_username = error.response.data.message;;
                this.loading = false;
            });
            this.newUsername = "";
        },
    },
    components: { LoadingSpinner }
}
</script>

<template>
	<div class="settings">
		<div class="big-text">Settings</div>
		<div class="settings-section">
			<div class="medium-text settings-title">Change your username</div>
        	<div class="search-mobile change-username">
      			<input type="text" v-model="this.newUsername" class="searchTerm change-username" placeholder="New fantastic username here">
      			<button type="submit" class="searchButton" @click="changeUsername">
        			<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
     			</button>
   			</div>
			<LoadingSpinner :load="this.loading"></LoadingSpinner>
			<ErrorMsg v-if="errormsg_username" :msg="errormsg_username"></ErrorMsg>
		</div>
	</div>
</template>

<style>
</style>
