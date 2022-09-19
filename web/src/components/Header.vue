<template>
  <div class="page-header relative">
    <div class="z-[10] absolute top-14" v-if="showMenu">
      <div class="flex flex-col bg-white border border-slate-300">
        <router-link class="px-10 py-3 hover:bg-slate-300" :to="{name: 'home'}">Trang chủ</router-link>
        <router-link class="px-10 py-3 hover:bg-slate-300" :to="{name: 'listPosts'}">Tin tức</router-link>
        <router-link class="px-10 py-3 hover:bg-slate-300" :to="{name: 'listEvents'}">Hoạt động</router-link>
        <router-link class="px-10 py-3 hover:bg-slate-300" :to="{name: 'committeePage'}">Tổ chức</router-link>
      </div>
    </div>

    <div>
      <div class="centered-horizontal bg-indigo-500 text-sm text-white p-3 pr-5 shadow-md shadow-slate-500 md:hidden">
        <div class="grow">
          <Bars3Icon class="w-8 h-8 cursor-pointer" @click="showMenu = !showMenu"></Bars3Icon>
        </div>
        <button v-if="!$root.isLoggedIn()" @click="logIn()">Đăng nhập</button>
        <button v-else @click="viewProfile()">Trang cá nhân</button>
      </div>

      <div class="pt-24 pb-10 relative">
        <div class="page-section w-fit md:w-full relative">
          <div class="xl:centered-horizontal gap-3 xl:px-10">

            <div class="grow centered-horizontal justify-center xl:justify-start gap-3">
              <img src="../assets/youth-logo.png" class="w-12 h-12"  alt=""/>
              <img src="../assets/das-logo.png" class="w-12 h-12"  alt=""/>
              <router-link to="/" class="text-3xl xl:text-4xl ml-5 font-light">
                <span class="block md:inline">ĐOÀN TRƯỜNG THPT</span> <span>DĨ AN</span>
                <p class="text-sm italic">Website đang hoạt động thử nghiệm</p>
              </router-link>
            </div>

            <div class="md:centered-horizontal justify-center gap-8 mt-10 xl:mt-0 hidden">
              <router-link class="border-b-2 border-b-transparent hover:border-b-slate-500" to="/">Trang chủ</router-link>
              <router-link class="border-b-2 border-b-transparent hover:border-b-slate-500" to="/p">Tin tức</router-link>
              <router-link class="border-b-2 border-b-transparent hover:border-b-slate-500" to="/e">Hoạt động</router-link>
              <router-link class="border-b-2 border-b-transparent hover:border-b-slate-500" to="/bch">Tổ chức</router-link>
              <button v-if="!$root.isLoggedIn()" @click="logIn()" class="btn-primary ml-10">Đăng nhập</button>
              <button v-else @click="viewProfile()" class="btn-primary ml-10">Trang cá nhân</button>
            </div>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import auth from "../auth/auth";
import {Bars3Icon} from '@heroicons/vue/24/solid';

export default {
  name: "Header",
  components: {
    Bars3Icon
  },
  data() {
    return {
      showMenu: false
    }
  },
  methods: {
    logIn() {
      auth.requestAuth(() => {
        window.location.reload()
      }, (e) => {
        this.$root.popupError(e)
      })
    },
    viewProfile() {
      this.$router.push("/u/" + this.$root.user.profile.id).then(() => {
        this.$router.go(0)
      })
    }
  }
}
</script>
