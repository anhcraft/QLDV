<template>
  <div class="centered-horizontal gap-3">
    <FireIcon class="w-8 h-8 text-rose-500"></FireIcon>
    <p class="text-3xl font-heading">Hoạt động Đoàn</p>
  </div>
  <div class="centered-horizontal gap-24 my-10">
    <div class="w-[500px]">
      <ActivityGallery></ActivityGallery>
    </div>
    <div>
      <div class="text-3xl font-heading mb-5">Sự kiện tháng 7</div>
      <div class="grid grid-cols-2 gap-x-5 gap-y-1">
        <EventButton>Event Đoàn mùa xuân</EventButton>
        <EventButton>Event Đoàn mùa hạ</EventButton>
        <EventButton>Event Đoàn mùa thu</EventButton>
        <EventButton>Event Đoàn mùa đông</EventButton>
        <EventButton>Event Đoàn mùa xuân</EventButton>
        <EventButton>Event Đoàn mùa hạ</EventButton>
        <EventButton>Event Đoàn mùa thu</EventButton>
        <EventButton>Event Đoàn mùa đông</EventButton>
      </div>
    </div>
  </div>
  <p class="text-3xl font-heading text-center mt-36">Cá nhân tiêu biểu</p>
  <div class="max-w-[600px] m-auto">
    <KeyMemberSlideshow></KeyMemberSlideshow>
  </div>
</template>

<script>
import EventButton from "../EventButton.vue";
import ActivityGallery from "./ActivityGallery.vue";
import KeyMemberSlideshow from "./KeyMemberSlideshow.vue";
import {FireIcon} from "@heroicons/vue/solid";

export default {
  name: "ActivitySection",
  components: {
    KeyMemberSlideshow,
    EventButton,
    ActivityGallery,
    FireIcon
  }
}
</script>
