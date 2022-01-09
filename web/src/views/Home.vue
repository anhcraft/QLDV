<template>
  <div class="h-screen grid grid-cols-6 gap-10 items-center bg-slate-300">
    <div class="col-span-2">
      <img src="src/assets/youth_logo.png" alt="" class="w-1/3 m-auto">
      <img src="src/assets/das_logo.png" alt="" class="w-1/3 mt-5 m-auto">
      <div class="border-l-4 border-l-slate-400 ml-20 mt-20">
        <div class="font-light leading-normal text-3xl ml-5">
          <p>CREATIVITY</p>
          <p>RESPONSIBILITY</p>
          <p>FOR COMMUNITY</p>
        </div>
      </div>
    </div>
    <div class="col-span-3">
      <div v-if="loadingPosts">
        <svg class="animate-spin h-16 w-16 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      <div v-else-if="posts.length > 0">
        <PostWidget :id="posts[0].id" :title="posts[0].title" large></PostWidget>
        <div class="grid grid-cols-2 gap-4 mt-10">
          <PostWidget v-for="value in posts.slice(1)" :id="value.id" :title="value.title"></PostWidget>
        </div>
      </div>
    </div>
  </div>
  <div class="h-screen grid grid-cols-1 gap-5 place-items-center">
    <div v-if="$root.isLoggedIn">
      <div class="grid grid-cols-2 gap-20" v-if="!$root.loadingProfile && !$root.progressionLoading">
        <CardWidget :pink="$root.profile.gender"></CardWidget>
        <div>
          <section class="mt-5" v-if="$root.profile.rates.length > 0">
            <p class="text-xl">Xếp hạng</p>
            <ul class="list-disc list-inside">
              <li v-for="(value, name) in $root.profile.rates">
                {{ value ? "Tốt" : "Khá" }} ({{ name }})
              </li>
            </ul>
          </section>
          <section class="mt-5" v-if="$root.profile.achievements.length > 0">
            <p class="text-xl">Thành tích</p>
            <ul class="list-disc list-inside">
              <li v-for="value in $root.profile.achievements">
                {{ value }}
              </li>
            </ul>
          </section>
        </div>
      </div>
      <div v-else>
        <svg class="animate-spin h-16 w-16 text-sky-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>
    <div v-else>
      <img src="../assets/elements.svg" alt="" class="w-1/2 m-auto">
      <button class="bg-white hover:bg-blue-300 cursor-pointer shadow-xl shadow-slate-400 border-4 border-blue-300 rounded-md text-xl px-5 py-3 w-64 text-center mt-20 block m-auto" @click="logIn">
        Đăng nhập bằng <img src="../assets/google.svg" alt="" class="inline-flex">
      </button>
    </div>
  </div>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2" v-if="$root.isLoggedIn">
    <NewspaperIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="managePosts" v-if="$root.profile.admin"></NewspaperIcon>
    <UsersIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="manageUsers" v-if="$root.profile.admin"></UsersIcon>
    <!--<CogIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2"></CogIcon>-->
    <LogoutIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="logOut"></LogoutIcon>
  </div>
</template>

<script>
import PostWidget from "../components/PostWidget.vue";
import auth from "../api/auth";
import CardWidget from "../components/CardWidget.vue";
import {CogIcon, LogoutIcon, NewspaperIcon, UsersIcon} from "@heroicons/vue/solid";
import server from "../api/server";

export default {
  name: "Home",
  components: {CardWidget, PostWidget, LogoutIcon, CogIcon, NewspaperIcon, UsersIcon},
  data() {
    return {
      loadingPosts: false,
      posts: []
    }
  },
  methods: {
    logIn() {
      auth.createSession()
    },
    logOut() {
      auth.destroySession()
    },
    managePosts() {
      this.$router.push('/pm/')
    },
    manageUsers() {
      this.$router.push('/um/')
    }
  },
  mounted() {
    if(!this.$root.progressionLoaded) {
      this.$root.loadProgression()
    }
    this.loadingPosts = true
    server.loadPosts(7, new Date().getTime()).then(s => {
      this.posts = s.posts
      this.loadingPosts = false
    })
  }
}
</script>
