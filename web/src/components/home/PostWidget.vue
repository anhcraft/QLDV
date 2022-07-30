<template>
  <router-link class="pw-container block cursor-pointer rounded-2xl before:rounded-2xl w-full h-full relative" :to="'/p/' + data.id">
    <img :src="getBackground(data)" class="object-cover object-center rounded-2xl w-full h-full" />
    <div class="w-full absolute text-white break-words z-10" :class="large ? 'bottom-10 px-10' : 'bottom-5 px-5'">
      <p class="font-heading" :class="large ? 'text-4xl font-bold' : 'text-xl font-semibold'">{{ data.title }}</p>
      <p class="text-sm mt-5" v-if="large">
        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc sed eros metus. Aenean id porttitor lorem. Nunc rhoncus leo accumsan, facilisis sem quis, sodales dolor
      </p>
    </div>
  </router-link>
</template>

<script>
import conf from "../../conf";

export default {
  name: "PostWidget",
  props: {
    data: Object,
    large: Boolean
  },
  methods: {
    getBackground(data) {
      if(data.hasOwnProperty("attachments")) {
        if(data.attachments.length > 0) {
          return conf.server + '/static/' + data.attachments[0].id
        }
      }
      return "https://i.imgur.com/aH4g7pj.jpg";
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
