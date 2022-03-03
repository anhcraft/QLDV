<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb text="Tin tức" link="/"></Breadcrumb>
    <LoadingState ref="loadingState">
      <div class="flex flex-row gap-5 place-content-end place-items-center text-slate-500 mb-3">
        <p class="text-sm">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.date)) }}</p>
        <div class="flex flex-row gap-1 text-xs">
          <EyeIcon class="w-4"></EyeIcon>
          <p>{{ post.views }}</p>
        </div>
        <div class="flex flex-row gap-1 border-2 border-white rounded-md px-2 py-1 text-xs transition-all duration-300" :class="
          (post.liked ? 'bg-pink-500 text-white hover:bg-pink-300': '') + ' ' +
          (this.$root.isLoggedIn() ? 'cursor-pointer border-pink-500' : '')" @click="likePost()">
          <HeartIcon class="w-4" :class="post.liked ? 'text-white' : 'text-pink-500'"></HeartIcon>
          <p>{{ post.likes }}</p>
        </div>
      </div>
      <article class="border-y-2 border-y-slate-300 py-10">
        <header class="text-4xl">{{ post.title }}</header>
        <section class="mt-5 break-words prose max-w-max" v-html="post.content"></section>
      </article>
      <div class="mt-10 flex flex-row flex-wrap gap-3" v-if="post.attachments.length > 0">
        <img v-for="att in post.attachments" class="max-h-36 cursor-pointer hover:opacity-80" :src="serverBaseURL + '/static/' + att.id" alt="" @click="previewImage(att.id)" />
      </div>
    </LoadingState>
  </div>
  <FloatingMenu></FloatingMenu>
  <div class="select-none" v-if="previewImageId !== undefined">
    <div class="bg-black opacity-75 fixed top-0 left-0 w-screen h-screen" v-on:mouseenter="zoomControlShow = false" @click="previewImage(undefined)"></div>
    <div class="md:hidden">
      <img class="w-full fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10 m-auto transition-all duration-300" :src="serverBaseURL + '/static/' + previewImageId" alt="" />
    </div>
    <div class="hidden md:block">
      <img :style="`width: ${this.previewImageSize}%`" v-on:mouseenter="zoomControlShow = true" class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10 m-auto transition-all duration-300" :src="serverBaseURL + '/static/' + previewImageId" alt="" />
      <div class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10 flex flex-row justify-center mt-1" v-if="zoomControlShow">
        <ZoomInIcon class="w-7 cursor-pointer p-1 bg-gray-300 hover:bg-gray-400" @click="zoomPreviewImg(1)"></ZoomInIcon>
        <ZoomOutIcon class="w-7 cursor-pointer p-1 bg-gray-300 hover:bg-gray-400" @click="zoomPreviewImg(-1)"></ZoomOutIcon>
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
import {EyeIcon, HeartIcon} from "@heroicons/vue/solid";
import auth from "../api/auth";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "Post",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    ZoomInIcon, ZoomOutIcon, EyeIcon, HeartIcon
  },
  data() {
    return {
      post: {},
      previewImageId: undefined,
      previewImageSize: 0,
      zoomControlShow: false
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
      this.zoomControlShow = false
    },
    zoomPreviewImg(base) {
      this.previewImageSize = Math.max(Math.min(this.previewImageSize + base * 10, 80), 50)
    },
    likePost(){
      if(!this.$root.isLoggedIn()) return
      server.updatePostStat(this.$route.params.id, "like", auth.getToken()).then(s => {
        if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          if (this.post.liked) {
            this.post.liked = false
            this.post.likes = this.post.likes - 1
          } else {
            this.post.liked = true
            this.post.likes = this.post.likes + 1
          }
        } else {
          this.$notify({
            title: "Cập nhật bài viết thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        }
      }, (e) => {
        this.$notify({
          title: "Lỗi hệ thống",
          text: e.message,
          type: "error"
        });
      })
    }
  },
  mounted() {
    server.loadPost(this.$route.params.id, auth.getToken()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$notify({
          title: "Tải bài viết thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
        return
      }
      this.post = s;
      this.$refs.loadingState.deactivate()
      server.updatePostStat(this.$route.params.id, "view", auth.getToken()).then(s => {
        if (s.hasOwnProperty("error") && s["error"] !== "ERR_TOKEN_VERIFY") {
          this.$notify({
            title: "Tải bài viết thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        } else if(s.hasOwnProperty("success") && s["success"]) {
          this.post.views = this.post.views + 1
        }
      })
    }, (e) => {
      this.$notify({
        title: "Tải bài viết thất bại",
        text: e.message,
        type: "error"
      });
    })
  }
}
</script>
