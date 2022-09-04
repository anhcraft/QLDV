<template>
  <div class="centered-horizontal gap-3">
    <FireIcon class="w-8 h-8 text-rose-500"></FireIcon>
    <p class="text-3xl font-heading">Hoạt động Đoàn</p>
  </div>
  <div class="lg:centered-horizontal gap-24 my-10">
    <div class="w-[500px]">
      <ActivityGallery></ActivityGallery>
    </div>
    <div>
      <div class="text-3xl font-heading mb-5 mt-10 lg:mt-0">Sự kiện tháng {{ new Date().getMonth() }}</div>
      <LoadingState ref="loadingStateOngoing">
        <div class="grid grid-cols-2 md:grid-cols-1 xl:grid-cols-2 gap-x-5 gap-y-1">
          <EventButton v-for="data in onGoingEvents" :data="data"></EventButton>
        </div>
      </LoadingState>
    </div>
  </div>
  <p class="text-3xl font-heading text-center mt-24 xl:mt-36">Cá nhân tiêu biểu</p>
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
import auth from "../../api/auth";
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
