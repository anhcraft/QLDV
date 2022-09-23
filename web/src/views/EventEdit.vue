<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <LoadingState ref="loadingState">
      <div class="flex flex-col gap-5 mb-10">
        <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tên sự kiện..." v-model.trim="event.title">

        <div class="flex flex-row gap-5 place-items-center">
          <p>Ngày bắt đầu</p>
          <Datepicker v-model="event.beginDate" locale="vi-VN" format="dd/MM/yyyy HH:mm"></Datepicker>
        </div>

        <div class="flex flex-row gap-5 place-items-center">
          <p>Ngày kết thúc</p>
          <Datepicker v-model="event.endDate" locale="vi-VN" format="dd/MM/yyyy HH:mm"></Datepicker>
        </div>

        <div class="border border-gray-400 py-2 px-5">
          <p class="text-xl">Giới hạn người xem</p>
          <select class="mt-5 text-sm" v-model.number="event.privacy">
            <option v-for="v in roleTables" :value="v.role">Từ {{ v.name }} trở lên</option>
          </select>
          <p class="mt-5 text-sm italic">
            {{ roleTables.filter(v => v.role >= event.privacy).map(v => v.name).join(", ") }} có thể xem
          </p>
        </div>

      </div>
      <button class="btn-success" :class="{'opacity-50' : submittingEvent}" @click="submit()">{{ $route.params.id === undefined ? "Thêm sự kiện" : "Lưu chỉnh sửa" }}</button>
    </LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import Datepicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
import LoadingState from "../components/LoadingState.vue";
import Footer from "../components/Footer.vue";
import {GetRoleTable} from "../auth/roles";
import {ServerError} from "../api/server-error";
import EventAPI from "../api/event-api";

export default {
  "name": "EventEdit",
  components: {Footer, LoadingState, Header, Breadcrumb, Datepicker},
  data() {
    return {
      event: {
        title: "",
        beginDate: null,
        endDate: null,
        privacy: 0
      },
      submittingEvent: false
    }
  },
  computed: {
    roleTables() {
      return GetRoleTable().filter(v => v.role <= this.$root.user.profile.role)
    }
  },
  methods: {
    submit() {
      if(this.submittingEvent) return
      if(this.event.title.length > 300){
        this.$root.popupError(new ServerError("ERROR_EVENT_TITLE_TOO_LONG"))
        return
      }
      if(this.event.title.length < 10){
        this.$root.popupError(new ServerError("ERROR_EVENT_TITLE_TOO_SHORT"))
        return
      }
      if(this.event.beginDate.getTime() > this.event.endDate.getTime()){
        this.$root.popupError(new ServerError("ERROR_EVENT_INVALID_DURATION"))
        return
      }
      this.submittingEvent = true
      const id = this.$route.params.id === undefined ? "" : this.$route.params.id
      EventAPI.updateEvent(id, {
        title: this.event.title,
        beginDate: this.event.beginDate.getTime(),
        endDate: this.event.endDate.getTime(),
        privacy: this.event.privacy,
      }).then(res => {
        this.submittingEvent = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          this.$router.push({name: "manageEvents"})
        }
      })
    }
  },
  mounted() {
    const f = () => {
      if(!this.$root.isLoggedIn() || !this.$root.isGlobalManager) {
        this.$router.push({name: "manageEvents"})
        return
      }
      if(this.$route.params.id !== undefined) {
        EventAPI.getEvent(this.$route.params.id).then(res => {
          if(res instanceof ServerError) {
            this.$root.popupError(res)
            return
          }
          this.event = res
          this.event.beginDate = new Date(res.beginDate)
          this.event.endDate = new Date(res.endDate)
          this.$refs.loadingState.deactivate()
        });
      } else {
        this.event.beginDate = new Date()
        this.event.endDate = new Date(this.event.beginDate.getTime() + 60 * 60 * 24 * 1000)
        this.$refs.loadingState.deactivate()
      }
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
