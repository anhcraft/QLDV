import {createApp} from 'vue'
import router from "./router";
import App from './App.vue'
import './index.css'
import './firebase'

const app = createApp(App)
app.use(router)
app.mount('#app')
