<script lang="ts">
export default {
  props: {
    id: String,
    title: String,
    large: Boolean
  }
}
</script>

<template>
  <router-link class="pw-container block cursor-pointer w-full relative" :class="[large ? 'h-72 shadow-xl shadow-slate-500' : 'h-32']" :to="'/p/' + id">
    <img src="../assets/blank.jpg" class="object-cover object-center w-full" :class="[large ? 'h-72' : 'h-32']" />
    <div class="bg-black absolute w-full h-full top-0 left-0 opacity-20 transition delay-75 duration-00 ease-in-out"></div>
    <div class="w-full absolute bottom-3 pl-5 pr-5 z-10">
      <p class="w-full text-white text-xl truncate">{{ title }}</p>
    </div>
  </router-link>
</template>

<style>
.pw-container:hover > div {
  opacity: 0.5 !important;
}
</style>
