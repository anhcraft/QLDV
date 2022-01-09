<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">THPT Dĩ An</span>
  </div>
  <div class="grid grid-cols-5 mt-36 mb-36">
    <article class="col-start-2 col-span-3" v-if="this.loaded">
      <header class="text-5xl">{{ post.title }}</header>
      <div class="mt-10 break-words" v-html="post.content"></div>
      <div class="border-t-2 border-t-slate-300 pt-10 mt-10 flex flex-row flex-wrap gap-3" v-if="post.attachments.length > 0">
        <img v-for="att in post.attachments" class="max-w-md cursor-pointer hover:opacity-80" :src="serverBaseURL + '/static/' + att.id" alt="" @click="previewImage(att.id)" />
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
  <div v-if="previewImageId !== undefined" @click="previewImage(undefined)">
    <div class="bg-black opacity-75 fixed top-0 left-0 w-screen h-screen"></div>
    <img class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10" :src="serverBaseURL + '/static/' + previewImageId" alt="" />
  </div>
</template>

<script>
import { ChevronDoubleUpIcon, ChevronDoubleDownIcon, HomeIcon } from '@heroicons/vue/solid'
import server from "../api/server";
import conf from "../conf";

export default {
  name: "Post",
  components: { ChevronDoubleUpIcon, ChevronDoubleDownIcon, HomeIcon },
  data() {
    return {
      post: {},
      loaded: false,
      previewImageId: undefined
    }
  },
  computed: {
    serverBaseURL() {
      return conf.server
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
    },
    previewImage(id) {
      this.previewImageId = id
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
