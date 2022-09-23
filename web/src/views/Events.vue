<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
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
</template>

<script>
import Header from "../components/Header.vue";
import Footer from "../components/Footer.vue";
import EventButton from "../components/EventButton.vue";
import LoadingState from "../components/LoadingState.vue";
import {CircleStackIcon, FireIcon} from '@heroicons/vue/24/solid';
import EventAPI from "../api/event-api";
import {ServerError} from "../api/server-error";

export default {
  name: "Events",
  components: {
    Header, Footer, LoadingState, FireIcon, EventButton, CircleStackIcon
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
      EventAPI.listEvents({
        limit: 10,
        "below-id": 0,
        "begin-date": t,
        "end-date": t,
      }).then((res) => {
        if(res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        this.onGoingEvents = res
        this.$refs.loadingStateOngoing.deactivate()
      })
    },
    loadNextEvents(){
      this.$refs.loadingStateAll.activate()
      EventAPI.listEvents({
        limit: 15,
        "below-id": this.pagination.belowId,
        "begin-date": 0,
        "end-date": 0,
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
          for(let i = 0; i < res.length; i++){
            this.indexEvent(res[i], i + this.events.length)
          }
          this.events = this.events.concat(res)
        }
        this.$refs.loadingStateAll.deactivate()
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
