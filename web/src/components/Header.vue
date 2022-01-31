<template>
  <div class="py-16 relative">
    <div class="max-w-[1024px] m-auto relative">
      <div class="flex flex-row gap-3 place-items-center">
        <img src="../assets/youth_logo.png" class="w-12 h-12"  alt=""/>
        <img src="../assets/das_logo.png" class="w-12 h-12"  alt=""/>
        <router-link to="/" class="text-3xl ml-5 font-light border-b-2 border-b-2-white hover:border-b-gray-500">ĐOÀN THPT DĨ AN</router-link>
        <div class="absolute right-0">
          <button v-if="!$root.isLoggedIn" @click="logIn()" class="transition-all duration-500 rounded-md bg-indigo-500 hover:bg-indigo-600 px-4 py-2 text-white text-sm">Đăng nhập</button>
          <button v-else @click="viewProfile()" class="transition-all duration-500 rounded-md bg-indigo-500 hover:bg-indigo-600 px-4 py-2 text-white text-sm">Trang cá nhân</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import auth from "../api/auth";

export default {
  name: "Header",
  methods: {
    logIn() {
      auth.createSession()
    },
    viewProfile() {
      this.$router.push("/u/")
    }
  }
}
</script>
