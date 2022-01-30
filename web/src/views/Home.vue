<template>
  <Header></Header>
  <div class="py-16">
    <div>
      <div v-for="(q, index) in quotes" :class="{'hidden': currentQuote !== index}">
        <img :src="q.img" alt="" class="h-[400px] m-auto transition-all duration-300 hover:opacity-80 cursor-pointer" @click="currentQuote = (currentQuote === quotes.length - 1) ? 0 : currentQuote + 1">
        <p class="max-w-[800px] h-[80px] m-auto mt-10 text-center">
          <q class="font-yomogi text-2xl">{{ q.text }}</q>
          <span class="font-serif text-gray-500 ml-5">― {{ q.author }}</span>
        </p>
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
  <div class="py-16">
    <div class="max-w-[1024px] m-auto grid grid-cols-6 gap-24">
      <div class="col-span-4">
        <div class="flex flex-row gap-3 place-items-center">
          <NewspaperIcon class="w-8 text-gray-600"></NewspaperIcon>
          <span class="font-light text-xl">TIN TỨC</span>
        </div>
        <div class="w-full mt-5 flex flex-col gap-4 mt-10" v-if="posts.length > 0">
          <PostWidget v-for="value in posts" :id="value.id" :title="value.title" :bg="getBg(value.attachments)"></PostWidget>
        </div>
        <div class="mt-10" v-if="loadingPosts">
          <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
      </div>
      <div class="col-span-2">
        <div class="flex flex-row gap-3 place-items-center">
          <CalendarIcon class="w-8 text-gray-600"></CalendarIcon>
          <span class="font-light text-xl">HOẠT ĐỘNG</span>
        </div>
        <div class="my-5">
          <div class="border-b-4 border-b-sky-300 pb-3 flex flex-row">
            <div class="text-xl grow">Tháng {{ eventCalendar.currentMonth + 1 }}/{{ eventCalendar.currentYear }}</div>
            <div class="p-1 cursor-pointer hover:bg-gray-300" @click="nextMonth(-1)">
              <ChevronLeftIcon class="w-4 text-gray-600"></ChevronLeftIcon>
            </div>
            <div class="p-1 cursor-pointer hover:bg-gray-300" @click="nextMonth(1)">
              <ChevronRightIcon class="w-4 text-gray-600"></ChevronRightIcon>
            </div>
          </div>
          <div class="grid grid-cols-7 gap-y-3 my-3">
            <div v-for="i in 35">
              <div class="m-auto w-7 h-7 rounded-full flex items-center justify-center" :class="{'bg-pink-400 text-white': isToday(i)}">
                <span v-if="i <= this.getDaysInMonth" class="text-sm">{{ i }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-3">
          <div class="border-2 border-dashed border-gray-400 rounded-xl px-5 py-2">
            <div class="text-lg">Event</div>
            <div class="text-sm text-gray-500">00/00 - 00/00</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import PostWidget from "../components/PostWidget.vue";
import {CalendarIcon, NewspaperIcon} from "@heroicons/vue/solid";
import server from "../api/server";
import conf from "../conf";
import treePlanting from "../assets/tree-planting.jpg"
import protectNature from "../assets/protect-nature.jpg"
import saveOcean from "../assets/save-ocean.jpg"
import springWallpaper from "../assets/spring-wallpaper.jpg"
import studyTogether from "../assets/study-together.jpg"
import {ChevronLeftIcon, ChevronRightIcon} from "@heroicons/vue/outline";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";

export default {
  name: "Home",
  components: {
    Header, PostWidget, FloatingMenu,
    NewspaperIcon, CalendarIcon, ChevronLeftIcon, ChevronRightIcon
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
      eventCalendar: {
        currentYear: 0,
        currentMonth: 0,
        events: []
      },
      loadingPosts: false,
      postAvailable: true,
      dateOffset: 0,
      posts: []
    }
  },
  computed: {
    getDaysInMonth() {
      return new Date(this.eventCalendar.currentYear, this.eventCalendar.currentMonth + 1, 0).getDate()
    }
  },
  methods: {
    getBg(a) {
      return a.length === 0 ? "" : (conf.server + '/static/' + a[0].id)
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.loadingPosts && this.postAvailable) {
          this.loadNextPosts()
        }
      }
    },
    loadNextPosts(){
      this.loadingPosts = true
      server.loadPosts(10, this.dateOffset).then(s => {
        if(s.posts.length === 0) {
          this.postAvailable = false
        } else {
          this.dateOffset = s.posts[s.posts.length - 1].date
        }
        this.posts = this.posts.concat(s.posts)
        this.loadingPosts = false
      })
    },
    nextMonth(delta){
      let d = this.eventCalendar.currentMonth + delta;
      if(d > 11) {
        d -= 12;
        this.eventCalendar.currentYear = this.eventCalendar.currentYear + 1
      } else if(d < 0) {
        d += 12;
        this.eventCalendar.currentYear = this.eventCalendar.currentYear - 1
      }
      this.eventCalendar.currentMonth = d
    },
    isToday(day) {
      const date = new Date()
      return day === date.getDate() && this.eventCalendar.currentMonth === date.getMonth() && this.eventCalendar.currentYear === date.getFullYear()
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    const date = new Date()
    this.eventCalendar.currentYear = date.getFullYear()
    this.eventCalendar.currentMonth = date.getMonth()
    this.dateOffset = date.getTime()
    this.loadNextPosts()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
