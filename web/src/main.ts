import {createApp} from 'vue'
import router from "./router";
import App from './App.vue'
import './index.css'
import './firebase'
import Notifications from '@kyvg/vue3-notification'
import auth from "./auth/auth";
import VueViewer from 'v-viewer'

auth.init(() => {
    const app = createApp(App)
    app.use(router)
    app.use(Notifications)
    app.use(VueViewer)
    app.mount('#app')
})
