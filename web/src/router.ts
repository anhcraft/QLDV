import {createRouter, createWebHashHistory} from "vue-router";

const routes = [
    { path: '/', component: () => import("./views/Home.vue") },
    { path: '/p/', component: () => import("./views/Home.vue") }, // fallback
    { path: '/p/:id', component: () => import("./views/Post.vue") },
    { path: '/pe/', component: () => import("./views/PostEdit.vue") },
    { path: '/pe/:id', component: () => import("./views/PostEdit.vue") },
    { path: '/pm/', component: () => import("./views/PostManage.vue") },
    { path: '/um/', component: () => import("./views/UserManage.vue") },
    { path: '/u/', component: () => import("./views/Home.vue") }, // fallback
    { path: '/u/:id', component: () => import("./views/Profile.vue") },
    { path: '/em/', component: () => import("./views/EventManage.vue") },
    { path: '/ee/', component: () => import("./views/EventEdit.vue") },
    { path: '/ee/:id', component: () => import("./views/EventEdit.vue") },
    { path: '/bch/', component: () => import("./views/CommitteePage.vue") },
    { path: '/c/', component: () => import("./views/Home.vue") }, // fallback
    { path: '/c/:id', component: () => import("./views/Contest.vue") },
    { path: '/mc/', component: () => import("./views/EventManage.vue") }, // fallback
    { path: '/mc/:id', component: () => import("./views/ContestManage.vue") },
    { path: '/mcs/', component: () => import("./views/EventManage.vue") }, // fallback
    { path: '/mcs/:id', component: () => import("./views/ContestSessionManage.vue") },
    { path: '/:pathMatch(.*)*', component: () => import("./views/NotFoundPage.vue") }
]

const router = createRouter({
    history: createWebHashHistory(),
    scrollBehavior (to: any, from: any, savedPosition: any) {
        return { left: 0, top: 0 }
    },
    routes,
})

export default router
