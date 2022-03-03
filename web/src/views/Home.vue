<template>
  <Header></Header>
  <div class="md:py-16">
    <div>
      <div v-for="(q, index) in quotes" :class="{'hidden': currentQuote !== index}">
        <div class="h-[200px] md:h-[400px] m-auto transition-all duration-300 hover:opacity-80 cursor-pointer bg-bottom bg-no-repeat bg-[length:100%] md:bg-[length:auto_400px]" :style="`background-image: url(${q.img})`" @click="currentQuote = (currentQuote === quotes.length - 1) ? 0 : currentQuote + 1"></div>
        <div class="max-w-[800px] h-[120px] m-auto mt-7 text-center">
          <q class="font-yomogi text-xl md:text-2xl">{{ q.text }}</q>
          <p class="font-serif text-gray-500 italic mt-1">― {{ q.author }}</p>
        </div>
      </div>
    </div>
    <div class="flex flex-row place-content-center gap-3">
      <div class="cursor-pointer" v-for="(_, index) in quotes" @click="currentQuote = index">
        <svg height="10" width="10">
          <circle cx="5" cy="5" r="5" :fill="currentQuote === index ? '#555' : '#aaa'" />
        </svg>
      </div>
    </div>
  </div>
  <div class="p-5 md:p-10 py-16">
    <div class="max-w-[1024px] m-auto grid grid-cols-1 md:grid-cols-6 md:gap-24">
      <div class="col-span-2 md:order-last">
        <Sidebar></Sidebar>
      </div>
      <div class="col-span-4 mt-10 md:mt-0">
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
import PostWidget from "../components/PostWidget.vue";
import {NewspaperIcon} from "@heroicons/vue/solid";
import server from "../api/server";
import conf from "../conf";
import treePlanting from "../assets/tree-planting.jpg"
import protectNature from "../assets/protect-nature.jpg"
import saveOcean from "../assets/save-ocean.jpg"
import springWallpaper from "../assets/spring-wallpaper.jpg"
import studyTogether from "../assets/study-together.jpg"
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import auth from "../api/auth";
import LoadingState from "../components/LoadingState.vue";
import Sidebar from "../components/Sidebar.vue";

export default {
  name: "Home",
  components: {
    LoadingState, Header, PostWidget, FloatingMenu, Sidebar, NewspaperIcon
  },
  data() {
    return {
      quotes: [
        {
          text: "Trees are poems that the earth writes upon the sky.",
          author: "Kahlil Gebran",
          img: treePlanting
        },
        {
          text: "Nature is painting for us, day after day, pictures of infinite beauty if only we have the eyes to see them.",
          author: "John Ruskin",
          img: protectNature
        },
        {
          text: "No water, no life. No blue, no green.",
          author: "Dr. Sylvia Earle",
          img: saveOcean
        },
        {
          text: "Một năm khởi đầu từ mùa xuân. Một đời khởi đầu từ tuổi trẻ. Tuổi trẻ là mùa xuân của xã hội.",
          author: "Ho Chi Minh",
          img: springWallpaper
        },
        {
          text: "Chúng ta phải học, phải cố gắng học nhiều. Không chịu khó học thì không tiến bộ được...",
          author: "Ho Chi Minh",
          img: studyTogether
        }
      ],
      currentQuote: 0,
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
