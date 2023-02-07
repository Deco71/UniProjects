import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import MyProfile from '../views/MyProfile.vue'
import Search from '../views/Search.vue'
import Post from '../views/PostPage.vue'
import NewPost from '../views/NewPost.vue'
import Settings from '../views/Settings.vue'
import Profile from '../views/Profile.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: Login},
		{path: '/home', component: HomeView},
		{path: '/me', component: MyProfile},
		{path: '/user/:name', component: Profile},
		{path: '/search', component: Search},
		{path: '/post/:id', component: Post},
		{path: '/upload', component: NewPost},
		{path: '/settings', component: Settings},
	],
})

router.beforeEach((to, from, next) => {
	if ((localStorage.getItem('id') == null) && (to.path != "/")) {
		next({
		  path: "/",
		});
	} else {
	  next();
	}
});

router.beforeEach((to, from, next) => {
	if ((localStorage.getItem('id') != null) && (to.path == "/")) {
		next({
		  path: "/home",
		});
	} else {
	  next();
	}
});


export default router
