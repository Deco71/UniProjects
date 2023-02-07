<script>
import { RouterLink, RouterView } from 'vue-router'
export default {
    data() {
        return {
            logged: false,
        };
    },
    mounted() {
        localStorage.clear();
        if (localStorage.getItem("id") !== null) {
            this.logged = true;
        }
        else {
            this.logged = false;
        }
    },
    methods: {
        logout() {
            localStorage.removeItem("id");
            localStorage.removeItem("username");
            this.$axios.defaults.headers.common["Authorization"] = null;
            this.logged = false;
        },
        login() {
            this.logged = true;
        }
    },
}
</script>

<template>


	<header v-if="this.logged">
		<div class="wrap">
			<RouterLink to="/home" class="nav-link header-link">
				<div class="title">DECOPHOTO</div>
			</RouterLink>
			<SearchBar searchStyle="search"/>
			<RouterLink to="/home" class="nav-link header-link">
							<svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
							Home
			</RouterLink>
			<RouterLink to="/search" class="nav-link header-link search-button">
							<svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
							Search
			</RouterLink>
			<RouterLink to="/me" class="nav-link header-link">
							<svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
							Profile
			</RouterLink>
			<RouterLink to="/" @click=logout class="nav-link header-link">
							<svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
							Logout
			</RouterLink>
		</div>
	</header>
	<main class="main">
			<RouterView @login="login"/>
	</main>
</template>


<style>
</style>
