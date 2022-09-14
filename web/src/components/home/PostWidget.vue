<template>
  <router-link class="pw-container block cursor-pointer rounded-2xl before:rounded-2xl w-full max-h-[300px] lg:max-h-full relative" :to="'/p/' + data.link">
    <img :src="getBackground(data)" class="object-cover object-center rounded-2xl w-full max-h-[300px] lg:max-h-full" />
    <div class="absolute text-white top-3 right-5 z-10">
      <div class="flex flex-row gap-1">
        <EyeIcon class="w-4"></EyeIcon>
        <p>{{ data.views }}</p>
      </div>
      <div class="flex flex-row gap-1">
        <HeartIcon class="w-4"></HeartIcon>
        <p>{{ data.likes }}</p>
      </div>
    </div>
    <div class="w-full absolute text-white break-words z-10 bottom-5 px-5" :class="large ? ' md:bottom-10 lg:px-10' : ''">
      <p class="font-heading" :class="large ? 'text-xl md:text-3xl xl:text-4xl md:font-bold' : 'text-xl font-semibold'">{{ data.title }}</p>
      <p class="font-heading">#{{ data.hashtag }}</p>
      <p class="text-sm mt-5 hidden md:block" v-if="large">
        {{ data.headline }}
      </p>
    </div>
  </router-link>
</template>

<script>
import conf from "../../conf";
import { HeartIcon, EyeIcon } from '@heroicons/vue/24/solid';

export default {
  name: "PostWidget",
  props: {
    data: Object,
    large: Boolean
  },
  components: {
    HeartIcon,
    EyeIcon
  },
  methods: {
    getBackground(data) {
      if(data.hasOwnProperty("attachments")) {
        if(data.attachments.length > 0) {
          return conf.server + '/static/' + data.attachments[0].id
        }
      }
      return "https://i.ibb.co/WFJ07mW/aH4g7pj.jpg";
    }
  }
}
</script>

<style scoped>
.pw-container:before {
  content: "";
  position: absolute;
  height: 50%;
  width: 100%;
  bottom: 0;
  left: 0;
  right: 0;
  background-image: linear-gradient(180deg, transparent, #333);
  z-index: 1;
}
.pw-container:hover::before  {
  height: 75%;
  background-image: linear-gradient(180deg, transparent, #111);
  transition: all 1s ease;
}
</style>
