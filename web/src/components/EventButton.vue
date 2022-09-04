<template>
  <div class="centered-horizontal shadow-lg shadow-slate-300 rounded-r-2xl transition-all duration-500"
               :class="(data.status === 'ongoing' ? 'bg-sky-50 hover:bg-sky-100 border border-zinc-300 hover:border-zinc-500' : (data.status === 'waiting' ? 'bg-violet-50 hover:bg-violet-100 opacity-50' : 'bg-zinc-50 hover:bg-zinc-100 opacity-50')) + (data.status === 'ongoing' && data.hasOwnProperty('contest') ? ' cursor-pointer' : '')"
                @click="openEvent(data)">
    <div class="w-16 h-16">
      <img src="https://i.ibb.co/zGYdGMX/qDCvJs8.jpg" class="object-cover object-center w-full h-full rounded-r-2xl" />
    </div>
    <div class="px-7">
      <p>{{ data.title }}</p>
      <p class="text-sm" v-if="data.status === 'finished'">Đã kết thúc</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "EventButton",
  props: {
    data: Object
  },
  methods: {
    openEvent(ev){
      if(ev.status === 'ongoing' && ev.hasOwnProperty("contest")) {
        this.$router.push(`/c/` + ev.link)
      }
    }
  }
}
</script>
