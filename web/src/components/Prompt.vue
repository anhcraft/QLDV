<template>
  <div v-if="show">
    <div class="bg-black opacity-75 fixed top-0 left-0 w-screen h-screen"></div>
    <div class="fixed w-full md:w-auto top-1/3 md:top-1/2 md:left-1/2 -translate-y-1/3 md:-translate-x-1/2 md:-translate-y-1/2 z-10 bg-white p-5">
      <div class="flex flex-col">
        <div>
          <slot></slot>
        </div>
        <div class="mt-5 flex flex-row gap-3">
          <button class="btn-success" @click="returnResult(true)">Có</button>
          <button class="btn-danger" @click="returnResult(false)">Không</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      show: false
    }
  },
  emits: ["callback"],
  methods: {
    toggle() {
      this.show = !this.show
    },
    returnResult(b) {
      this.toggle()
      this.$emit('callback', b)
    }
  }
}
</script>
