<template>
  <div class="w-96 h-56 m-auto bg-red-100 rounded-xl relative text-white shadow-2xl shadow-slate-500 transition-transform transform hover:scale-110">
    <img class="relative object-cover w-full h-full rounded-xl" :class="{grayscale : !$root.profile.certified}" src="../assets/card_pink.png" alt="" v-if="pink">
    <img class="relative object-cover w-full h-full rounded-xl" :class="{grayscale : !$root.profile.certified}" src="../assets/card.png" alt="" v-else>
    <div class="w-full px-8 absolute top-8">
      <div class="flex justify-between">
        <div class="">
          <p class="font-light">Name</p>
          <p class="font-medium tracking-widest">{{ $root.profile.name.toUpperCase() }}</p>
        </div>
        <div class="flex gap-2">
          <img class="w-14 h-14" src="src/assets/youth_logo.png" alt="">
        </div>
      </div>
      <div class="pt-1">
        <p class="font-light">Card Number</p>
        <p class="font-medium tracking-more-wider">{{ $root.formattedStudentId }}</p>
      </div>
      <div class="pt-6 pr-6">
        <div class="flex justify-between">
          <div class="">
            <p class="font-light text-xs">Entry</p>
            <p class="font-medium tracking-wider text-sm">{{ $root.profile.entryDate }}</p>
          </div>
          <div class="">
            <p class="font-light text-xs text-xs">End</p>
            <p class="font-medium tracking-wider text-sm">{{ $root.profile.endDate }}</p>
          </div>
          <div class="">
            <p class="font-light text-xs">Class</p>
            <p class="font-medium tracking-more-wider text-sm">{{ $root.profile.class }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CardWidget",
  props: {
    certified: Boolean,
    pink: Boolean
  }
}
</script>
