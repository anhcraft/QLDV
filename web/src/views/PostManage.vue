<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 text-center text-sm" @click="createPost">Tạo bài viết</button>
    <table class="w-full mt-10">
      <tbody>
        <tr v-for="post in posts">
          <td>{{ post.title }}</td>
          <td class="flex flex-row gap-3">
            <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="editPost(post.id)"></PencilIcon>
            <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="removePost(post.id, post.title)"></TrashIcon>
            <p class="text-gray-500">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.date)) }}</p>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="loadingPosts">
      <svg class="animate-spin h-6 w-6 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
    <div class="mt-10" v-else-if="postAvailable">
      <button class="bg-white hover:bg-blue-300 cursor-pointer border-2 border-blue-300 px-3 py-1 text-center text-sm" @click="loadNextPosts">Xem thêm...</button>
    </div>
    <div class="mt-10" v-else>Đã tải hết bài viết.</div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt :content="'<p class=font-bold>Bạn có muốn xóa bài viết này?</p><br>' + postRemoveTitle" @callback="removePostCallback" ref="removePrompt"></Prompt>
</template>

<script>
import {PencilIcon, TrashIcon} from '@heroicons/vue/solid'
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";

export default {
  name: "PostManage",
  components: {
    Header, FloatingMenu,
    PencilIcon, TrashIcon, Prompt
  },
  data() {
    return {
      loadingPosts: false,
      postAvailable: true,
      posts: [],
      postRemoveId: '',
      postRemoveTitle: ''
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
      server.loadPosts(20, older).then(s => {
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
    removePost(id, name) {
      this.postRemoveId = id
      this.postRemoveTitle = name
      this.$refs.removePrompt.toggle()
    },
    removePostCallback(b) {
      if(b) {
        server.removePost(this.postRemoveId, auth.getToken()).then(s => {
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.posts = this.posts.filter(p => p.id !== this.postRemoveId)
            this.postRemoveId = ""
            this.postRemoveTitle = ""
          } else {
            alert(`Lỗi xóa bài viết: ${s["error"]}`)
          }
        })
      }
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn) {
      this.$router.push(`/`)
    }
    this.loadNextPosts()
  }
}
</script>
