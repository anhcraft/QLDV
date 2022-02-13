<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Quản lý bài viết" link="/pm"></Breadcrumb>
    <div class="mt-10">
      <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 text-center text-sm" @click="createPost">Tạo bài viết</button>
    </div>
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
    <div class="mt-10">
      <LoadingState ref="loadingState">
        <div v-if="postAvailable">
          <button class="rounded-md bg-blue-500 hover:bg-blue-600 cursor-pointer px-3 py-2 text-white text-center text-xs m-auto block" @click="loadNextPosts()">Xem thêm...</button>
        </div>
        <div v-else>Đã tải hết bài viết.</div>
      </LoadingState>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt @callback="removePostCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa bài viết này?</p><br> {{ postRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, TrashIcon} from '@heroicons/vue/solid'
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";

export default {
  name: "PostManage",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    PencilIcon, TrashIcon, Prompt
  },
  data() {
    return {
      postAvailable: true,
      posts: [],
      postRemoveId: '',
      postRemoveTitle: ''
    }
  },
  methods: {
    loadNextPosts(){
      this.$refs.loadingState.activate()
      const older = this.posts.length === 0 ? new Date().getTime() : this.posts[this.posts.length - 1].date
      server.loadPosts(20, older, auth.getToken()).then(s => {
        if(s.posts.length === 0) {
          this.postAvailable = false
        }
        this.posts = this.posts.concat(s.posts)
        this.$refs.loadingState.deactivate()
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
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    this.loadNextPosts()
  }
}
</script>
