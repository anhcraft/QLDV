<template>
  <div class="centered-horizontal gap-3">
    <FireIcon class="w-8 h-8 text-rose-500"></FireIcon>
    <p class="text-3xl font-heading">Hoạt động Đoàn</p>
  </div>
  <div class="lg:flex content-start xl:gap-24 my-10">
    <div class="md:w-[600px] md:h-[450px] md:rotate-[-5deg] relative m-auto">
      <img src="../../assets/activity-gallery-bg.svg" class="absolute w-full h-full z-10 hidden md:block">
      <div class="md:absolute top-20 left-24 right-24 bottom-20">
        <ActivityGallery></ActivityGallery>
      </div>
    </div>
    <div class="mt-5 xl:mt-0">
      <div class="text-3xl font-heading text-center">Câu lạc bộ</div>
      <div class="grid grid-cols-2 gap-5 xl:gap-14 font-heading text-lg text-center">
        <div class="self-end">
          <img src="../../assets/thedayband.png" class="w-48 m-auto">
          <p>CLB Âm Nhạc</p>
        </div>
        <div>
          <img src="../../assets/Basketball-bro.png" class="w-48 h-48 m-auto">
          <p>CLB Bóng rổ</p>
        </div>
        <div>
          <img src="../../assets/Volleyball-bro.png" class="w-48 h-48 m-auto">
          <p>CLB Bóng chuyền</p>
        </div>
        <div>
          <img src="../../assets/Badminton-bro.png" class="w-48 h-48 m-auto">
          <p>CLB Bóng bàn</p>
        </div>
      </div>
    </div>
  </div>
  <LoadingState ref="loadingStateOngoing">
    <div v-if="onGoingEvents.length > 0">
      <p class="text-3xl font-heading text-center mt-16 xl:mt-24">Sự kiện tháng {{ new Date().getMonth() }}</p>
      <div class="max-w-[700px] md:h-[460px] m-auto mt-10 relative shadow-lg shadow-slate-300">
        <img src="../../assets/event-notes.png" class="absolute w-full h-full hidden md:block">
        <div class="relative z-5 max-h-[460px] md:px-16 py-10 overflow-auto">
          <div class="grid grid-cols-2 md:grid-cols-1 xl:grid-cols-2 gap-x-5 gap-y-1">
            <EventButton v-for="data in onGoingEvents" :data="data"></EventButton>
          </div>
        </div>
      </div>
    </div>
  </LoadingState>
  <p class="text-3xl font-heading text-center mt-16 xl:mt-24">Cá nhân tiêu biểu</p>
  <div class="max-w-[600px] m-auto">
    <KeyMemberSlideshow></KeyMemberSlideshow>
  </div>
</template>

<script>
import EventButton from "../EventButton.vue";
import ActivityGallery from "./ActivityGallery.vue";
import KeyMemberSlideshow from "./KeyMemberSlideshow.vue";
import { FireIcon } from '@heroicons/vue/24/solid';
import server from "../../api/server";
import auth from "../../auth/auth";
import LoadingState from "../LoadingState.vue";

export default {
  name: "ActivitySection",
  components: {
    KeyMemberSlideshow,
    EventButton,
    ActivityGallery,
    FireIcon,
    LoadingState
  },
  data() {
    return {
      onGoingEvents: []
    }
  },
  methods: {
    loadOngoingEvents(){
      this.$refs.loadingStateOngoing.activate()
      const t = new Date().getTime()
      server.loadEvents(8, 0, t, t, auth.getToken()).then(s => {
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
  },
  mounted() {
    this.loadOngoingEvents()
  }
}
</script>
