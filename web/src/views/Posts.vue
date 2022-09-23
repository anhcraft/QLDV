<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <div class="centered-horizontal mb-5 gap-3">
      <RssIcon class="w-8 h-8 text-rose-500"></RssIcon>
      <p class="text-3xl font-heading">Tin tức</p>
    </div>
    <div class="centered-horizontal gap-10">
      <div class="grow centered-horizontal gap-3">
        <div class="btn-outline-sm"
             v-for="tag in Object.keys(pagination.hashtags)"
             :class="{'bg-violet-500 text-white' : pagination.hashtags[tag]}"
             @click="this.switchHashtag(tag)">#{{ tag }}</div>
      </div>
      <div class="centered-horizontal gap-1">
        <CalendarIcon class="w-8 h-8 btn-outline-sm"
                      :class="{'bg-violet-500 text-white' : pagination.sortBy === ''}"
                      @click="this.setSortBy('')"></CalendarIcon>
        <EyeIcon class="w-8 h-8 btn-outline-sm"
                 :class="{'bg-violet-500 text-white' : pagination.sortBy === 'view'}"
                 @click="this.setSortBy('view')"></EyeIcon>
        <HeartIcon class="w-8 h-8 btn-outline-sm"
                   :class="{'bg-violet-500 text-white' : pagination.sortBy === 'like'}"
                   @click="this.setSortBy('like')"></HeartIcon>
      </div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5 mt-10">
      <PostWidget v-for="data in posts" :data="data"></PostWidget>
    </div>
    <LoadingState ref="loadingState"></LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import Footer from "../components/Footer.vue";
import LoadingState from "../components/LoadingState.vue";
import PostWidget from "../components/home/PostWidget.vue";
import {CalendarIcon, EyeIcon, HeartIcon, RssIcon} from "@heroicons/vue/24/solid";
import PostAPI from "../api/post-api";
import {ServerError} from "../api/server-error";

export default {
  name: "Posts",
  components: {
    Header, Footer, PostWidget, LoadingState, RssIcon, HeartIcon, EyeIcon, CalendarIcon
  },
  data() {
    return {
      posts: [],
      pagination: {
        belowId: 0,
        lowerThan: 0,
        sortBy: "",
        available: true,
        hashtags: {}
      }
    }
  },
  methods: {
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.loadingState.loading && this.pagination.available) {
          this.loadNextPosts()
        }
      }
    },
    loadNextPosts(){
      this.$refs.loadingState.activate()
      let hashtags = []
      Object.keys(this.pagination.hashtags).forEach(t => {
        if(this.pagination.hashtags[t]) {
          hashtags.push(t)
        }
      })
      PostAPI.listPosts({
        limit: 15,
        "below-id": this.pagination.belowId,
        "filter-hashtags": hashtags,
        "sort-by": this.pagination.sortBy,
        "lower-than": this.pagination.lowerThan
      }).then((res) => {
        if(res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        if(res.length < 10) {
          this.pagination.available = false
        }
        if(res.length > 0) {
          this.pagination.belowId = res[res.length - 1].id
          this.posts = this.posts.concat(res)
        }
        this.$refs.loadingState.deactivate()
      })
    },
    resetPosts() {
      this.posts = []
      this.pagination.belowId = 0
      this.pagination.lowerThan = 0
      this.pagination.available = true
      this.loadNextPosts()
    },
    switchHashtag(tag) {
      this.pagination.hashtags[tag] = !this.pagination.hashtags[tag]
      this.resetPosts();
    },
    setSortBy(w){
      if(w !== this.pagination.sortBy) {
        this.pagination.sortBy = w
        this.resetPosts();
      }
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    PostAPI.getHashtags().then(data => {
      if(data instanceof ServerError) {
        this.$root.popupError(data)
        return
      }
      data.forEach(e => this.pagination.hashtags[e] = false)
      if(this.$route.query.hasOwnProperty("tag")) {
        let tag = this.$route.query["tag"]
        if(this.pagination.hashtags.hasOwnProperty(tag)) {
          this.pagination.hashtags[tag] = true
        }
      }
      this.loadNextPosts()
    })
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
