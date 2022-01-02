<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">Quản lý bài viết</span>
  </div>
  <div class="grid grid-cols-5 mt-36 mb-36">
    <div class="col-start-2 col-span-3 flex flex-col gap-5 ">
      <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 w-32 text-center" @click="createPost">Tạo bài viết</button>
      <div class="w-full mt-10" v-for="post in posts">
        <div class="w-full flex flex-row gap-3">
          <p class="text-xl grow">{{ post.title }}</p>
          <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="editPost(post.id)"></PencilIcon>
          <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="removePost(post.id)"></TrashIcon>
          <p class="text-gray-500">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.date)) }}</p>
        </div>
      </div>
      <div v-if="loadingPosts">
        <svg class="animate-spin h-16 w-16 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      <div class="mt-10" v-else-if="postAvailable">
        <button class="bg-white hover:bg-blue-300 cursor-pointer border-2 border-blue-300 px-3 py-1 w-32 text-center" @click="loadNextPosts">Xem thêm...</button>
      </div>
      <div class="mt-10" v-else>Đã tải hết bài viết.</div>
    </div>
  </div>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2">
    <HomeIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="backToHome"></HomeIcon>
    <ChevronDoubleUpIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToTop"></ChevronDoubleUpIcon>
  </div>
</template>

<script>
import {ChevronDoubleUpIcon, HomeIcon, PencilIcon, TrashIcon} from '@heroicons/vue/solid'
import server from "../api/server";

export default {
  name: "Post",
  components: { ChevronDoubleUpIcon, HomeIcon, PencilIcon, TrashIcon },
  data() {
    return {
      loadingPosts: false,
      postAvailable: true,
      posts: []
    }
  },
  methods: {
    jumpToTop() {
      window.scrollTo(0, 0);
    },
    backToHome() {
      this.$router.push('/')
    },
    loadNextPosts(){
      this.loadingPosts = true
      const older = this.posts.length === 0 ? new Date().getTime() : this.posts[this.posts.length - 1].date
      server.loadPosts(5, older).then(s => {
        if(s.posts.length === 0) {
          this.postAvailable = false
        }
        this.posts = this.posts.concat(s.posts)
        this.loadingPosts = false
      })
    },
    createPost() {
      this.$router.push(`/pe/`)
    },
    editPost(id) {
      this.$router.push(`/pe/${id}/`)
    },
    removePost(id) {

    }
  },
  mounted() {
    this.loadNextPosts()
  }
}
</script>
