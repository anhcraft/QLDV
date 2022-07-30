<template>
  <router-link class="centered-horizontal shadow-lg shadow-slate-300 rounded-r-2xl bg-sky-50 hover:bg-sky-100 transition-all duration-500" to="">
    <div class="w-16 h-16">
      <img src="https://i.imgur.com/qDCvJs8.jpg" class="object-cover object-center w-full h-full rounded-r-2xl" />
    </div>
    <div class="px-7">
      <slot></slot>
    </div>
  </router-link>
</template>

<script>
export default {
  name: "EventButton"
}
</script>
