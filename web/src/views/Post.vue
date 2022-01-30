<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Tin tức" link="/"></Breadcrumb>
    <div v-if="this.loaded">
      <div class="text-right text-slate-500 mb-3">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.date)) }}</div>
      <article class="border-y-2 border-y-slate-300 py-10">
        <header class="text-4xl">{{ post.title }}</header>
        <section id="content" class="mt-5 break-words" v-html="post.content"></section>
      </article>
      <div class="mt-10 flex flex-row flex-wrap gap-3" v-if="post.attachments.length > 0">
        <img v-for="att in post.attachments" class="max-w-xs cursor-pointer hover:opacity-80" :src="serverBaseURL + '/static/' + att.id" alt="" @click="previewImage(att.id)" />
      </div>
    </div>
    <div v-else>
      <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
  <div v-if="previewImageId !== undefined">
    <div class="bg-black opacity-75 fixed top-0 left-0 w-screen h-screen" @click="previewImage(undefined)"></div>
    <div class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10 select-none w-4/5">
      <img :style="`width: ${this.previewImageSize}%`" class="m-auto transition-all	duration-300" :src="serverBaseURL + '/static/' + previewImageId" alt="" />
      <div class="flex flex-row justify-center mt-1">
        <ZoomInIcon class="w-7 cursor-pointer p-1 bg-gray-300 hover:bg-gray-400" @click="this.previewImageSize = Math.min(this.previewImageSize + 10, 80)"></ZoomInIcon>
        <ZoomOutIcon class="w-7 cursor-pointer p-1 bg-gray-300 hover:bg-gray-400" @click="this.previewImageSize = Math.max(this.previewImageSize - 10, 50)"></ZoomOutIcon>
      </div>
    </div>
  </div>
</template>

<script>
import server from "../api/server";
import conf from "../conf";
import Header from "../components/Header.vue";
import {ZoomInIcon, ZoomOutIcon} from "@heroicons/vue/outline";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
  name: "Post",
  components: {
    Header, FloatingMenu, Breadcrumb,
    ZoomInIcon, ZoomOutIcon
  },
  data() {
    return {
      post: {},
      loaded: false,
      previewImageId: undefined,
      previewImageSize: 0
    }
  },
  computed: {
    serverBaseURL() {
      return conf.server
    }
  },
  methods: {
    previewImage(id) {
      this.previewImageId = id
      this.previewImageSize = 50
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

<style>
#content a {
  color: rgb(38 143 207);
}
#content ol {
  display: block;
  list-style-type: decimal;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  padding-inline-start: 40px;
}
#content ul {
  display: block;
  list-style-type: disc;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  padding-inline-start: 40px;
}
</style>
