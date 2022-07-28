<script lang="ts">
export default {
  props: {
    id: String,
    title: String,
    bg: String
  }
}
</script>

<template>
  <router-link class="pw-container block cursor-pointer w-full relative h-36 md:h-42" :to="'/p/' + id">
    <img v-if="bg === undefined || bg.length === 0" src="../../assets/vn-coast.jpg" class="object-cover object-center w-full h-36" />
    <img v-else :src="bg" class="object-cover object-center w-full h-36" />
    <div class="pw-image bg-black absolute w-full h-full top-0 left-0 opacity-20 transition delay-75 duration-300 ease-in-out"></div>
    <div class="w-full absolute bottom-3 pl-5 pr-5 z-10">
      <p class="w-full text-white md:text-xl truncate">{{ title }}</p>
    </div>
  </router-link>
</template>

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
.pw-container:hover > .pw-image {
  opacity: 0.5 !important;
}
</style>
