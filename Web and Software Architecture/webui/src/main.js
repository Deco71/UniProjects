import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Immagine from './components/Immagine.vue'
import Post from './components/Post.vue'
import Comments from './components/Comments.vue'
import Feed from './components/Feed.vue'
import Likes from './components/Likes.vue'
import SearchBar from './components/SearchBar.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Immagine", Immagine);
app.component("Post", Post);
app.component("Comments", Comments);
app.component("Feed", Feed);
app.component("Likes", Likes);
app.component("SearchBar", SearchBar);
app.use(router)
app.mount('#app')