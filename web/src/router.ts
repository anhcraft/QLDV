import {createRouter, createWebHashHistory} from "vue-router";

const routes = [
    { name: 'home', path: '/', component: () => import("./views/Home.vue") },
    { name: 'listPosts', path: '/bai-viet/', component: () => import("./views/Posts.vue") },
    { name: 'readPost', path: '/bai-viet/:id', component: () => import("./views/Post.vue") },
    { name: 'createPost', path: '/quan-ly/tao-bai-viet/', component: () => import("./views/PostEdit.vue") },
    { name: 'updatePost', path: '/quan-ly/sua-bai-viet/:id', component: () => import("./views/PostEdit.vue") },
    { name: 'managePosts', path: '/quan-ly/bai-viet/', component: () => import("./views/PostManage.vue") },
    { name: 'manageUsers', path: '/quan-ly/tai-khoan/', component: () => import("./views/UserManage.vue") },
    { name: 'profile', path: '/u/:id', component: () => import("./views/Profile.vue") },
    { name: 'listEvents', path: '/su-kien/', component: () => import("./views/Events.vue") },
    { name: 'manageEvents', path: '/quan-ly/su-kien/', component: () => import("./views/EventManage.vue") },
    { name: 'createEvent', path: '/quan-ly/tao-su-kien/', component: () => import("./views/EventEdit.vue") },
    { name: 'updateEvent', path: '/quan-ly/sua-su-kien/:id', component: () => import("./views/EventEdit.vue") },
    { name: 'committeePage', path: '/ban-chap-hanh/', component: () => import("./views/CommitteePage.vue") },
    { name: 'manageSettings', path: '/cai-dat/', component: () => import("./views/SettingManage.vue") },
    /*
    { name: 'home', path: '/c/', component: () => import("./views/Home.vue") }, // fallback
    { name: 'home', path: '/c/:id', component: () => import("./views/Contest.vue") },
    { name: 'home', path: '/mc/', component: () => import("./views/EventManage.vue") }, // fallback
    { name: 'home', path: '/mc/:id', component: () => import("./views/ContestManage.vue") },
    { name: 'home', path: '/mcs/', component: () => import("./views/EventManage.vue") }, // fallback
    { name: 'home', path: '/mcs/:id', component: () => import("./views/ContestSessionManage.vue") },*/
    { name: 'notFound', path: '/:pathMatch(.*)*', component: () => import("./views/NotFoundPage.vue") }
]

const router = createRouter({
    history: createWebHashHistory(),
    scrollBehavior (to: any, from: any, savedPosition: any) {
        return { left: 0, top: 0 }
    },
    routes,
})

export default router
