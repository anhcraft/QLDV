<template>
  <div v-if="loading">
    <div v-if="!hidden" class="flex flex-row place-items-center m-auto">
      <svg class="animate-spin h-5 w-5 text-sky-400 block m-auto" xmlns="http://www.w3.org/2000/svg" fill="none"
           viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <div class="text-sm">{{ text }}</div>
    </div>
  </div>
  <slot v-else></slot>
</template>

<script>
export default {
  name: "LoadingState",
  props: {
    hidden: {
      type: Boolean,
      default: false
    },
    text: String
  },
  data() {
    return {
      loading: true
    }
  },
  methods: {
    activate() {
      this.loading = true
    },
    deactivate() {
      this.loading = false
    }
  }
}
</script>
