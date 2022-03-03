import {createApp} from 'vue'
import router from "./router";
import App from './App.vue'
import './index.css'
import './firebase'
import Notifications from '@kyvg/vue3-notification'

const app = createApp(App)
app.use(router)
app.use(Notifications)
app.mount('#app')
