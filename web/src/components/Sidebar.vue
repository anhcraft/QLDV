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
      <div class="border-2 border-dashed border-gray-400 rounded-xl px-5 py-2" v-for="event in eventCalendar.events[eventCalendar.currentMonth+'.'+eventCalendar.currentYear]" :class="{'hover:border-gray-800 cursor-pointer' : event.hasOwnProperty('contest')}" @click="openEvent(event)">
        <div class="text-lg">{{ event.title }}</div>
        <div class="text-sm text-gray-500">
          {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.startDate)) }} -
          {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.endDate)) }}
        </div>
      </div>
    </div>
  </LoadingState>
  <div class="my-10 flex flex-col gap-3">
    <router-link to="/bch/" class="text-white text-xl px-10 py-5 shadow-lg shadow-slate-400" style="background: radial-gradient(circle, rgba(34,120,195,1) 0%, rgba(162,45,253,1) 100%);">BCH Công Đoàn</router-link>
    <router-link to="/" class="text-white text-xl px-10 py-5 shadow-lg shadow-slate-400" style="background: radial-gradient(circle, rgba(238,174,202,1) 0%, rgba(148,187,233,1) 100%);">Hội thi ATGT</router-link>
    <router-link to="/" class="text-white text-xl px-10 py-5 shadow-lg shadow-slate-400" style="background: linear-gradient(21deg, rgba(58,137,180,1) 0%, rgba(13,191,35,1) 100%);">Thư viện ảnh</router-link>
  </div>
</template>

<script>
import {ChevronLeftIcon, ChevronRightIcon} from "@heroicons/vue/outline";
import {CalendarIcon} from "@heroicons/vue/solid";
import LoadingState from "../components/LoadingState.vue";
import server from "../api/server";
import auth from "../api/auth";

export default {
  name: "Sidebar",
  components: {
    LoadingState, CalendarIcon, ChevronLeftIcon, ChevronRightIcon
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
      if(event.hasOwnProperty("contest")) {
        this.$router.push("/c/" + event.contest.id)
      }
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
