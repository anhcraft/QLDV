<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb :text="($route.params.id === undefined ? 'Tạo' : 'Sửa') + ' sự kiện'" link="/em" class="mb-10"></Breadcrumb>
    <div v-if="eventLoaded && !submittingEvent">
      <div class="flex flex-col gap-5 mb-20">
        <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tên sự kiện..." v-model="event.title">
        <div class="flex flex-row gap-5 place-items-center">
          <p>Ngày bắt đầu</p>
          <Datepicker v-model="event.startDate" locale="vi-VN" format="dd/MM/yyyy HH:mm"></Datepicker>
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p>Ngày kết thúc</p>
          <Datepicker v-model="event.endDate" locale="vi-VN" format="dd/MM/yyyy HH:mm"></Datepicker>
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p>Chỉ cho thành viên xem</p>
          <input type="checkbox" class="w-4 h-4" v-bind:checked="(event.privacy & 1) === 1" @input="event.privacy = $event.target.value ? (event.privacy ^ 1) : (event.privacy | 1)">
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p>Chỉ cho bí thư xem</p>
          <input type="checkbox" class="w-4 h-4" v-bind:checked="(event.privacy & 2) === 2" @input="event.privacy = $event.target.value ? (event.privacy ^ 2) : (event.privacy | 2)">
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p>Chỉ cho quản trị viên xem</p>
          <input type="checkbox" class="w-4 h-4" v-bind:checked="(event.privacy & 4) === 4" @input="event.privacy = $event.target.value ? (event.privacy ^ 4) : (event.privacy | 4)">
        </div>
      </div>
      <button class="bg-pink-400 hover:bg-pink-500 cursor-pointer px-4 py-2 text-white text-center text-sm" v-if="!submittingEvent" @click="submit()">{{ $route.params.id === undefined ? "Thêm sự kiện" : "Lưu chỉnh sửa" }}</button>
    </div>
    <div v-else>
      <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import server from "../api/server";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import Datepicker from 'vue3-date-time-picker';
import 'vue3-date-time-picker/dist/main.css'

export default {
  "name": "EventEdit",
  components: {Header, FloatingMenu, Breadcrumb, Datepicker},
  data() {
    return {
      event: {
        title: "",
        startDate: null,
        endDate: null,
        date: null,
        privacy: 0
      },
      eventLoaded: false,
      submittingEvent: false
    }
  },
  methods: {
    submit() {
      this.submittingEvent = true
      server.changeEvent(this.$route.params.id, this.event, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.submittingEvent = false
          this.$router.push('/em/')
        } else {
          alert(`Lỗi lưu sự kiện: ${s["error"]}`)
        }
      })
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    if(this.$route.params.id !== undefined) {
      server.loadEvent(this.$route.params.id, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error")) {
          this.event = s;
          this.event.startDate = new Date(s.startDate);
          this.event.endDate = new Date(s.endDate);
          this.eventLoaded = true;
        } else {
          alert(`Lỗi tải sự kiện: ${s["error"]}`)
        }
      });
    } else {
      this.eventLoaded = true;
      this.event.startDate = new Date()
      this.event.endDate = new Date()
    }
  }
}
</script>
