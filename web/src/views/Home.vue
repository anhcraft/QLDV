<template>
  <Header></Header>
  <div class="max-w-[1100px] m-auto flex flex-col gap-3 mt-5 px-3">
    <div class="flex flex-row place-items-center text-sm">
      <div class="whitespace-nowrap px-3 py-1 bg-blue-800 text-white">
        Thông báo
      </div>
      <div class="contents">
        <marquee-text class="py-1 bg-blue-200">
          Chào mừng kỷ niệm 91 năm Ngày thành lập Đoàn TNCS Hồ Chí Minh (26/3/1931 - 26/3/2022). Chung kết cấp tỉnh "Tự hào Việt Nam" lần thứ IV năm 2022.
        </marquee-text>
      </div>
    </div>
    <div class="mb-5 h-[300px]">
      <SlideshowWidget :images="slideshow"></SlideshowWidget>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-7 md:gap-16 mt-5">
      <div class="col-span-2 md:order-last">
        <Sidebar></Sidebar>
      </div>
      <div class="col-span-5 mt-5 md:mt-0 md:pr-10">
        <div class="flex flex-row gap-3 place-items-center">
          <NewspaperIcon class="w-8 text-gray-600"></NewspaperIcon>
          <span class="font-light text-xl">TIN TỨC</span>
        </div>
        <div class="w-full flex flex-col gap-4 mt-5" v-if="posts.length > 0">
          <PostWidget v-for="value in posts" :id="value.id" :title="value.title" :bg="getBg(value.attachments)"></PostWidget>
        </div>
        <div class="mt-10">
          <LoadingState ref="postLoadingState"></LoadingState>
        </div>
      </div>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import PostWidget from "../components/home/PostWidget.vue";
import {NewspaperIcon} from "@heroicons/vue/solid";
import server from "../api/server";
import conf from "../conf";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import auth from "../api/auth";
import LoadingState from "../components/LoadingState.vue";
import Sidebar from "../components/Sidebar.vue";
import MarqueeText from "vue-marquee-text-component";
import SlideshowWidget from "../components/home/SlideshowWidget.vue";

export default {
  name: "Home",
  components: {
    LoadingState, Header, PostWidget, FloatingMenu, Sidebar, NewspaperIcon, MarqueeText, SlideshowWidget
  },
  data() {
    return {
      slideshow: [
        "https://i.imgur.com/qK6gzc0.jpg",
        "https://i.imgur.com/CvXZJy4.jpg"
      ],
      postAvailable: true,
      dateOffset: 0,
      posts: []
    }
  },
  methods: {
    getBg(a) {
      return a.length === 0 ? "" : (conf.server + '/static/' + a[0].id)
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.postLoadingState.loading && this.postAvailable) {
          this.loadNextPosts()
        }
      }
    },
    loadNextPosts(){
      this.$refs.postLoadingState.activate()
      server.loadPosts(10, this.dateOffset, auth.getToken()).then(s => {
        if(s.posts.length === 0) {
          this.postAvailable = false
        } else {
          this.dateOffset = s.posts[s.posts.length - 1].date
        }
        this.posts = this.posts.concat(s.posts)
        this.$refs.postLoadingState.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải bài viết thất bại",
          text: e.message,
          type: "error"
        });
      })
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    this.dateOffset = new Date().getTime()
    this.loadNextPosts()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
