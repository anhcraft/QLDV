<template>
  <div class="py-16 relative">
    <div class="max-w-[1024px] m-auto relative">
      <div class="flex flex-row gap-3 place-items-center">
        <img src="../assets/youth_logo.png" class="w-12 h-12"  alt=""/>
        <img src="../assets/das_logo.png" class="w-12 h-12"  alt=""/>
        <p class="text-3xl ml-5 font-light">ĐOÀN THPT DĨ AN</p>
        <div class="absolute right-0">
          <button v-if="!$root.isLoggedIn" @click="logIn()" class="transition-all duration-300 border-2 border-sky-500 rounded-md font-light hover:text-white hover:bg-sky-500 px-5 py-2 flex flex-row gap-2 place-items-center">
            <img src="../assets/google.svg" class="w-6 h-6" alt="">
            <span>Đăng nhập</span>
          </button>
          <button v-else @click="viewProfile()" class="transition-all duration-300 border-2 border-indigo-500 rounded-md font-light hover:text-white hover:bg-indigo-500 px-5 py-2">
            <span>Trang cá nhân</span>
          </button>
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
  }
}
</script>
