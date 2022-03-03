<template>
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
  <LoadingState ref="eventCalendarLoadingState">
    <div class="flex flex-col gap-3">
      <div class="border-2 border-dashed border-gray-400 rounded-xl px-5 py-2" v-for="event in eventCalendar.events[eventCalendar.currentMonth+'.'+eventCalendar.currentYear]" :class="{'hover:border-gray-800 cursor-pointer' : (event.hasOwnProperty('contest') && $root.isLoggedIn())}" @click="openEvent(event)">
        <div class="text-lg break-words">
          <FireIcon class="w-6 text-gray-600 text-rose-400 float-left mr-1" v-if="event.hasOwnProperty('contest')"></FireIcon>
          <p>{{ event.title }}</p>
        </div>
        <div class="text-sm text-gray-500">{{ getEventStatus(event) }}</div>
      </div>
    </div>
  </LoadingState>
  <div class="my-5 md:my-10 flex flex-col gap-3">
    <router-link to="/bch/" class="bg-blue-800 hover:opacity-80 transition transition-all text-white text-lg py-3 px-5 shadow-lg shadow-slate-400 flex flex-row place-items-center gap-3">
      <img src="../assets/youth-logo.png" class="w-12 h-12"  alt=""/>
      <p>BCH Đoàn trường</p>
    </router-link>
    <a href="https://youtube.com/channel/UCPtKpRuGva1y2RwJCooVu5Q" target="_blank" class="bg-indigo-800 hover:opacity-80 transition transition-all text-white text-lg py-3 px-5 shadow-lg shadow-slate-400 flex flex-row place-items-center gap-3">
      <img src="../assets/youtube-btn.png" class="w-12 h-12"  alt=""/>
      <p>Video hoạt động</p>
    </a>
  </div>
</template>

<script>
import {ChevronLeftIcon, ChevronRightIcon} from "@heroicons/vue/outline";
import {CalendarIcon, FireIcon} from "@heroicons/vue/solid";
import LoadingState from "../components/LoadingState.vue";
import server from "../api/server";
import auth from "../api/auth";

export default {
  name: "Sidebar",
  components: {
    LoadingState, CalendarIcon, ChevronLeftIcon, ChevronRightIcon, FireIcon
  },
  data() {
    return {
      eventCalendar: {
        currentYear: 0,
        currentMonth: 0,
        events: {}
      }
    }
  },
  computed: {
    getDaysInMonth() {
      return new Date(this.eventCalendar.currentYear, this.eventCalendar.currentMonth + 1, 0).getDate()
    }
  },
  methods: {
    openEvent(event){
      if(event.hasOwnProperty("contest") && this.$root.isLoggedIn()) {
        this.$router.push("/c/" + event.contest.id)
      }
    },
    getEventStatus(event){
      if((event.endDate - new Date().getTime()) >= 60*60*3*1000) {
        return ""
      }
      return "Sắp kết thúc: " + new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.endDate))
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

      const key = this.eventCalendar.currentMonth + "." + this.eventCalendar.currentYear;
      if(this.eventCalendar.events.hasOwnProperty(key)) return
      this.$refs.eventCalendarLoadingState.activate()
      const a = new Date(this.eventCalendar.currentYear, this.eventCalendar.currentMonth, 1, 0, 0, 0)
      const b = new Date(this.eventCalendar.currentYear, this.eventCalendar.currentMonth + 1, 1, 0, 0, 0)
      server.loadEvents(10, new Date().getTime(), a.getTime(), b.getTime() - 1000, auth.getToken()).then(s => {
        const v = this.eventCalendar.events
        v[key] = s.events
        this.eventCalendar.events = v
        this.$refs.eventCalendarLoadingState.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    isToday(day) {
      const date = new Date()
      return day === date.getDate() && this.eventCalendar.currentMonth === date.getMonth() && this.eventCalendar.currentYear === date.getFullYear()
    }
  },
  mounted() {
    const date = new Date()
    this.eventCalendar.currentYear = date.getFullYear()
    this.eventCalendar.currentMonth = date.getMonth()
    this.nextMonth(0)
  }
}
</script>

<style scoped>

</style>
