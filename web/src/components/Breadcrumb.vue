<template>
  <div class="flex flex-row gap-1 place-items-center">
    <HomeIcon class="w-4 cursor-pointer" @click="this.$router.push({ name: 'home' })"></HomeIcon>
    <ChevronRightIcon class="w-4"></ChevronRightIcon>
    <div class="cursor-pointer" @click="this.$router.push(link).then(() => {this.$router.go(0)})">{{ text }}</div>
    <ChevronRightIcon class="w-4"></ChevronRightIcon>
  </div>
</template>

<script>
import {ChevronRightIcon, HomeIcon} from '@heroicons/vue/24/solid';

export default {
  name: "Breadcrumb",
  components: {
    ChevronRightIcon,
    HomeIcon,
  },
  props: {
    text: String,
    link: Object
  }
}
</script>

