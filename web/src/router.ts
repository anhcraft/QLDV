import {createRouter, createWebHashHistory} from "vue-router";

const routes = [
    { path: '/', component: () => import("./views/Home.vue") },
    { path: '/p/:id', component: () => import("./views/Post.vue") },
    { path: '/pe/', component: () => import("./views/PostEdit.vue") },
    { path: '/pe/:id', component: () => import("./views/PostEdit.vue") },
    { path: '/pm/', component: () => import("./views/PostManage.vue") }
]

const router = createRouter({
    history: createWebHashHistory(),
    scrollBehavior (to: any, from: any, savedPosition: any) {
        return { left: 0, top: 0 }
    },
    routes,
})

export default router
