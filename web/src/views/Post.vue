<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">THPT Dĩ An</span>
  </div>
  <div class="grid grid-cols-5 mt-36 mb-36">
    <article class="col-start-2 col-span-3" v-if="this.loaded">
      <header class="text-5xl">L{{ this.post.title }}</header>
      <div class="mt-10 break-all">
        {{ this.post.content }}
      </div>
    </article>
    <div class="col-start-2 col-span-3" v-else>
      <svg class="animate-spin h-16 w-16 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
  </div>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2">
    <HomeIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="backToHome"></HomeIcon>
    <ChevronDoubleUpIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToTop"></ChevronDoubleUpIcon>
    <ChevronDoubleDownIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToBottom"></ChevronDoubleDownIcon>
  </div>
</template>

<script>
import { ChevronDoubleUpIcon, ChevronDoubleDownIcon, HomeIcon } from '@heroicons/vue/solid'
import server from "../api/server";

export default {
  name: "Post",
  components: { ChevronDoubleUpIcon, ChevronDoubleDownIcon, HomeIcon },
  data() {
    return {
      post: {},
      loaded: false
    }
  },
  methods: {
    jumpToTop() {
      window.scrollTo(0, 0);
    },
    jumpToBottom() {
      window.scrollTo(0, document.body.scrollHeight);
    },
    backToHome() {
      this.$router.push('/')
    }
  },
  mounted() {
    server.loadPost(this.$route.params.id).then(s => {
      this.post = s;
      this.loaded = true
    })
  }
}
</script>
