<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">{{ $route.params.id === undefined ? "Tạo" : "Sửa" }} bài viết</span>
  </div>
  <div class="grid grid-cols-5 mt-36 mb-36">
    <div class="col-start-2 col-span-3 flex flex-col gap-5">
      <div v-if="postLoaded">
        <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tiêu đề..." v-model="post.title">
        <div class="mt-10">
          <TiptapEditor :content="post.content" @onChange="this.onContentChange"></TiptapEditor>
        </div>
        <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 w-36 text-center" v-if="!submittingPost" @click="submitPost">{{ $route.params.id === undefined ? "Đăng bài" : "Lưu chỉnh sửa" }}</button>
      </div>
      <div v-else>
        <svg class="animate-spin h-16 w-16 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>
  </div>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2">
    <ArrowLeftIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="backToManage"></ArrowLeftIcon>
    <ChevronDoubleUpIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToTop"></ChevronDoubleUpIcon>
    <ChevronDoubleDownIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToBottom"></ChevronDoubleDownIcon>
  </div>
</template>

<script>
import {ArrowLeftIcon, ChevronDoubleDownIcon, ChevronDoubleUpIcon} from '@heroicons/vue/solid'
import TiptapEditor from "../components/TiptapEditor.vue";
import server from "../api/server";
import auth from "../api/auth";

export default {
  "name": "PostEdit",
  components: {TiptapEditor, ChevronDoubleUpIcon, ChevronDoubleDownIcon, ArrowLeftIcon },
  data() {
    return {
      post: {
        title: "",
        content: ""
      },
      postLoaded: false,
      submittingPost: false
    }
  },
  methods: {
    jumpToTop() {
      window.scrollTo(0, 0);
    },
    jumpToBottom() {
      window.scrollTo(0, document.body.scrollHeight);
    },
    backToManage() {
      this.$router.push('/pm/')
    },
    submitPost() {
      this.submittingPost = true
      server.changePost(this.$route.params.id, this.post.title, this.post.content, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$router.push('/pm/')
        } else {
          alert(`Lỗi lưu bài viết: ${s["error"]}`)
        }
      })
    },
    onContentChange(content) {
      this.post.content = content
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn) {
      this.$router.push(`/`)
    }
    if(this.$route.params.id !== undefined) {
      server.loadPost(this.$route.params.id).then(s => {
        this.post = s;
        this.postLoaded = true;
      });
    } else {
      this.postLoaded = true;
    }
  }
}
</script>
