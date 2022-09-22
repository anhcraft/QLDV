<template>
  <Header></Header>
  <section class="page-section px-10 py-8 lg:py-16">
    <LoadingState ref="loadingState">
      <div class="centered-horizontal gap-3 text-slate-500 mb-3">
        <div class="grow"></div>
        <p>
          <router-link class="text-cyan-500 hover:underline" :to="{ name: 'listPosts', query: { tag: post.hashtag} }">#{{ post.hashtag }}</router-link>
        </p>
        <p class="text-sm">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.updateDate)) }}</p>
        <div class="flex flex-row gap-1 text-xs">
          <EyeIcon class="w-4"></EyeIcon>
          <p>{{ post.stats.views }}</p>
        </div>
        <div class="flex flex-row gap-1 border-2 border-white rounded-md px-2 py-1 text-xs transition-all duration-300" :class="
          (post.stats.liked ? 'bg-pink-500 text-white hover:bg-pink-300': '') + ' ' +
          (this.$root.isLoggedIn() ? 'cursor-pointer border-pink-500' : '')" @click="likePost()">
          <HeartIcon class="w-4" :class="post.stats.liked ? 'text-white' : 'text-pink-500'"></HeartIcon>
          <p>{{ post.stats.likes }}</p>
        </div>
      </div>
      <article class="border-t-2 border-t-slate-300 py-10">
        <header class="text-3xl md:text-4xl">{{ post.title }}</header>
        <section class="mt-5 break-words prose max-w-max" v-html="post.content"></section>
      </article>
      <div class="mt-10">
        <viewer :images="attachments">
          <div class="grid grid-cols-6 gap-5">
            <img v-for="src in attachments" :key="src" :src="src" class="w-full h-[150px] object-contain border border-slate-400 hover:opacity-70 cursor-pointer">
          </div>
        </viewer>
      </div>
    </LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import conf from "../conf";
import Header from "../components/Header.vue";
import {EyeIcon, HeartIcon} from '@heroicons/vue/24/solid';
import LoadingState from "../components/LoadingState.vue";
import Footer from "../components/Footer.vue";
import PostAPI from "../api/post-api";
import {ServerError} from "../api/server-error";
import VueViewer from 'v-viewer'
import 'viewerjs/dist/viewer.css'

export default {
  name: "Post",
  components: {
    LoadingState, Header, Footer, EyeIcon, HeartIcon, VueViewer
  },
  data() {
    return {
      post: {}
    }
  },
  computed: {
    postId() {
      let s = this.$route.params.id.split(".")
      s = s[s.length - 1]
      s = s.replace(/\D/i, s)
      return s
    },
    attachments() {
      return this.post.attachments.map(v => conf.assetURL + "/" + v.id)
    }
  },
  methods: {
    likePost(){
      if(!this.$root.isLoggedIn()) return
      PostAPI.updatePostStat(this.postId, {
        like: !this.post.stats.liked,
        view: undefined
      }).then(s => {
        if (s instanceof ServerError) {
          this.$root.popupError(s)
        } else {
          this.post.stats.likes = s.likes
          this.post.stats.liked = !this.post.stats.liked
        }
      })
    }
  },
  mounted() {
    const f = () => {
      PostAPI.getPost(this.postId).then(s => {
        if(s instanceof ServerError) {
          this.$root.popupError(s)
          return
        }
        this.post = s;
        this.$refs.loadingState.deactivate()
        PostAPI.updatePostStat(this.postId, {
          like: undefined,
          view: true
        }).then(s => {
          if (s instanceof ServerError) {
            this.$root.popupError(s)
          } else {
            this.post.stats.views = s.views
            this.post.stats.viewed = true
          }
        })
      })
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
