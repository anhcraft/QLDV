<template>
  <div class="w-full bg-indigo-500 text-right text-sm text-white p-3 pr-5 shadow-md shadow-slate-500 md:hidden">
    <button v-if="!$root.isLoggedIn()" @click="logIn()">Đăng nhập</button>
    <button v-else @click="viewProfile()">Trang cá nhân</button>
  </div>
  <div class="py-8 md:py-10 relative">
    <div class="max-w-[1024px] w-fit md:w-full m-auto relative">
      <div class="flex flex-row gap-3 place-items-center md:p-10">
        <img src="../assets/youth_logo.png" class="w-12 h-12"  alt=""/>
        <img src="../assets/das_logo.png" class="w-12 h-12"  alt=""/>
        <router-link to="/" class="text-3xl ml-5 font-light border-b-2 border-b-white hover:border-b-gray-500">
          <span class="block md:inline">ĐOÀN THPT</span> <span>DĨ AN</span>
        </router-link>
        <div class="absolute right-10 hidden md:block">
          <button v-if="!$root.isLoggedIn()" @click="logIn()" class="btn-primary">Đăng nhập</button>
          <button v-else @click="viewProfile()" class="btn-primary">Trang cá nhân</button>
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
      this.$router.push("/u/" + this.$root.profile.email.substring(0, this.$root.profile.email.search("@"))).then(() => {
        this.$router.go(0)
      })
    }
  }
}
</script>
