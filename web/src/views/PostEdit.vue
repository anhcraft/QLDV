<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">{{ $route.params.id === undefined ? "Tạo" : "Sửa" }} bài viết</span>
  </div>
  <div class="grid grid-cols-5 mt-36 mb-36">
    <div class="col-start-2 col-span-3 flex flex-col gap-5 ">
      <input type="text" class="border-b-2 border-b-slate-300 text-3xl" placeholder="Tiêu đề...">
      <TiptapEditor></TiptapEditor>
      <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 w-36 text-center" @click="submitPost">{{ $route.params.id === undefined ? "Đăng bài" : "Lưu chỉnh sửa" }}</button>
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

export default {
  "name": "Post",
  components: {TiptapEditor, ChevronDoubleUpIcon, ChevronDoubleDownIcon, ArrowLeftIcon },
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
    jumpToBottom() {
      window.scrollTo(0, document.body.scrollHeight);
    },
    backToManage() {
      this.$router.push('/pm/')
    }
  },
  mounted() {

  }
}
</script>
