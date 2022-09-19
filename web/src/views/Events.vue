<template>
  <Header></Header>
  <section class="page-section px-10 py-8 lg:py-16">
    <div class="centered-horizontal mb-5 gap-3">
      <FireIcon class="w-8 h-8 text-rose-500"></FireIcon>
      <p class="text-3xl font-heading">Đang diễn ra</p>
    </div>
    <LoadingState ref="loadingStateOngoing">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5 mt-10">
        <EventButton v-for="data in onGoingEvents" :data="data"></EventButton>
      </div>
    </LoadingState>
    <div class="centered-horizontal mt-16 xl:mt-24 gap-3">
      <CircleStackIcon class="w-8 h-8 text-rose-500"></CircleStackIcon>
      <p class="text-3xl font-heading">Tất cả</p>
    </div>
    <div class="mt-10 flex flex-col gap-10">
      <div v-for="ent in Object.keys(scheduler)">
        <p class="text-2xl font-heading">{{ent}}</p>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5 mt-5">
          <EventButton v-for="data in scheduler[ent]" :data="events[data]"></EventButton>
        </div>
      </div>
    </div>
    <LoadingState ref="loadingStateAll"></LoadingState>
  </section>
  <Footer></Footer>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import Footer from "../components/Footer.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import EventButton from "../components/EventButton.vue";
import LoadingState from "../components/LoadingState.vue";
import {CircleStackIcon, FireIcon} from '@heroicons/vue/24/solid';
import server from "../api/server";
import auth from "../auth/auth";

export default {
  name: "Events",
  components: {
    Header, Footer, FloatingMenu, LoadingState, FireIcon, EventButton, CircleStackIcon
  },
  data() {
    return {
      events: [],
      scheduler: {},
      onGoingEvents: [],
      pagination: {
        belowId: 0,
        available: true
      }
    }
  },
  methods: {
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.loadingStateAll.loading && this.pagination.available) {
          this.loadNextEvents()
        }
      }
    },
    loadOngoingEvents(){
      this.$refs.loadingStateOngoing.activate()
      const t = new Date().getTime()
      server.loadEvents(20, 0, t, t, auth.getToken()).then(s => {
        this.onGoingEvents = s.events
        this.$refs.loadingStateOngoing.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    loadNextEvents(){
      this.$refs.loadingStateAll.activate()
      server.loadEvents(10, this.pagination.belowId, 0, 0, auth.getToken()).then(s => {
        if(s.events.length < 10) {
          this.pagination.available = false
        }
        if(s.events.length > 0) {
          this.pagination.belowId = s.events[s.events.length - 1].id
          for(let i = 0; i < s.events.length; i++){
            this.indexEvent(s.events[i], i + this.events.length)
          }
          this.events = this.events.concat(s.events)
        }
        this.$refs.loadingStateAll.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    indexEvent(ev, index) {
      let target = new Date(ev.beginDate)
      let cursor = new Date(ev.endDate)
      if(target.getTime() > cursor.getTime()) return
      while (true) {
        const m = cursor.getMonth()
        const y = cursor.getFullYear()
        const k = m + "/" + y
        let arr = []
        if(this.scheduler.hasOwnProperty(k)) {
          arr = this.scheduler[k]
        }
        arr.push(index)
        this.scheduler[k] = arr
        if(m === target.getMonth() && y === target.getFullYear()) {
          break
        }
        cursor.setMonth(cursor.getMonth() - 1)
      }
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    this.loadOngoingEvents()
    this.loadNextEvents()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
