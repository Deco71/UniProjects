<script>
export default {
	emits: ["login"],
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			username: "",
		}
	},
	methods: {
		async refresh() {
			this.errormsg = null;
		},
		login() {
			this.errormsg = null;
			this.$axios.post("/session", {
				username: this.username,
			}).then(
				res=> {
					localStorage.setItem("id", res.data.identifier);
					localStorage.setItem("username", this.username);
					this.$axios.defaults.headers.common["Authorization"] = `Bearer ${localStorage.getItem("id")}`;
					this.$emit("login");
					this.$router.push("/home")
				}
			).catch(
				err=> {
					this.errormsg = err.response.data.message;
				}	
			);
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
    <div class="login-block">
        <h1>DecoPhoto</h1>
        <input type="text" v-model="this.username"  placeholder="Username" id="username" />
        <button @click.prevent="login">Login</button>
    </div>
</template>

<style>
</style>
