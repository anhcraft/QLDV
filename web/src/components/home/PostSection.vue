<template>
  <div class="centered-horizontal mb-5">
    <div class="centered-horizontal gap-3">
      <RssIcon class="w-8 h-8 text-rose-500"></RssIcon>
      <p class="text-3xl font-heading">Tin tức</p>
    </div>
    <router-link class="centered-horizontal gap-1 ml-auto text-slate-500 hover:text-black transition-all duration-300" to="/p">
      <p class="text-sm">Xem thêm</p>
      <ArrowRightIcon class="w-3 h-3"></ArrowRightIcon>
    </router-link>
  </div>
  <LoadingState ref="postLoadingState">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-y-5 lg:gap-x-5">
      <div class="col-span-2">
        <PostWidget :data="posts[0]" large v-if="posts.length > 0"></PostWidget>
      </div>
      <div class="col-span-1">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-none lg:grid-rows-2 gap-5" v-if="posts.length > 1">
          <PostWidget v-for="data in posts.slice(1)" :data="data"></PostWidget>
        </div>
      </div>
    </div>
  </LoadingState>
</template>

<script>
import PostWidget from "./PostWidget.vue";
import LoadingState from "../LoadingState.vue";
import server from "../../api/server";
import auth from "../../auth/auth";
import { RssIcon, ArrowRightIcon } from '@heroicons/vue/24/solid';

export default {
  name: "PostSection",
  components: {
    LoadingState,
    PostWidget,
    RssIcon,
    ArrowRightIcon
  },
  data(){
    return {
      posts: []
    }
  },
  methods: {
    loadPosts(){
      this.$refs.postLoadingState.activate()
      server.loadPosts(3, [], "", 0, 0, auth.getToken()).then(s => {
        this.posts = s.posts
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
  mounted() {
    this.loadPosts()
  }
}
</script>
