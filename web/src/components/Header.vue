<template>
  <div class="page-header relative">
    <div class="z-[10] absolute top-14" v-if="showMenu">
      <div class="bg-white border border-slate-300">
        <div class="flex flex-col">
          <router-link class="px-10 py-3 hover:bg-slate-300" v-for="v in menu" :to="{name: v.route}">{{
              v.name
            }}
          </router-link>
        </div>
        <div class="flex flex-col border-t border-t-slate-300" v-if="$root.isManager">
          <router-link class="px-10 py-3 hover:bg-slate-300" v-for="v in adminMenu" :to="{name: v.route}">{{
              v.name
            }}
          </router-link>
        </div>
      </div>
    </div>

    <div>
      <div class="centered-horizontal gap-7 bg-indigo-500 text-sm text-white p-3 shadow-md shadow-slate-500"
           :class="{'md:hidden' : !$root.isManager}">
        <div class="w-8 h-8">
          <Bars3Icon class="w-8 h-8 cursor-pointer mr-5 md:hidden" @click="showMenu = !showMenu"></Bars3Icon>
        </div>
        <div class="grow hidden md:block" v-if="$root.isManager">
          <div class="centered-horizontal justify-center gap-5">
            <router-link class="border-b-2 border-b-transparent hover:border-b-white" v-for="v in adminMenu"
                         :to="{name: v.route}">{{ v.name }}
            </router-link>
          </div>
        </div>
        <div class="ml-auto mr-5 md:hidden">
          <button v-if="!$root.isLoggedIn()" @click="logIn()">Đăng nhập</button>
          <button v-else @click="viewProfile()">Trang cá nhân</button>
        </div>
      </div>

      <div class="pt-16 md:pt-24 pb-10 relative">
        <div class="page-section w-fit md:w-full relative">
          <div class="xl:centered-horizontal gap-3 xl:px-10">

            <div class="grow centered-horizontal justify-center xl:justify-start gap-3">
              <img src="../assets/youth-logo.png" class="w-12 h-12" alt=""/>
              <img src="../assets/das-logo.png" class="w-12 h-12" alt=""/>
              <router-link :to="{ name: 'home' }" class="text-3xl xl:text-4xl ml-5 font-light">
                <span class="block md:inline">ĐOÀN TRƯỜNG THPT</span> <span>DĨ AN</span>
                <p class="text-sm italic">Website đang hoạt động thử nghiệm</p>
              </router-link>
            </div>

            <div class="md:centered-horizontal justify-center gap-8 mt-10 xl:mt-0 hidden">
              <router-link class="border-b-2 border-b-transparent hover:border-b-slate-500" v-for="v in menu"
                           :to="{name: v.route}">{{ v.name }}
              </router-link>
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
      showMenu: false,
      menu: [
        {
          "name": "Trang chủ",
          "route": "home"
        },
        {
          "name": "Tin tức",
          "route": "listPosts"
        },
        {
          "name": "Hoạt động",
          "route": "listEvents"
        },
        {
          "name": "Tổ chức",
          "route": "committeePage"
        }
      ]
    }
  },
  computed: {
    adminMenu() {
      const s = [
        {
          "name": "Quản lý tài khoản",
          "route": "manageUsers"
        },
      ]
      if (this.$root.isGlobalManager) {
        s.push({
          "name": "Quản lý bài viết",
          "route": "managePosts"
        })
        s.push({
          "name": "Quản lý hoạt động",
          "route": "manageEvents"
        })
        s.push({
          "name": "Cài đặt hệ thống",
          "route": "manageSettings"
        })
      }
      return s
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
      this.$router.push("/u/" + this.$root.user.profile.email.substring(0, this.$root.user.profile.email.search("@")))
    }
  }
}
</script>
