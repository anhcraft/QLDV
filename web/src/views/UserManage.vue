<template>
  <div class="bg-white shadow-md shadow-slate-300 fixed z-10 left-0 top-0 w-screen p-3">
    <img src="src/assets/das_logo.png" alt="" class="h-10 inline-flex" />
    <span class="text-xl ml-5">Quản lý thành viên</span>
  </div>
  <div class="grid grid-cols-7 mt-36 mb-36">
    <div class="col-start-2 col-span-5 flex flex-col gap-5">
      <div class="flex flex-row gap-5">
        <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 w-52 text-center">Nhập dữ liệu thành viên</button>
        <button class="bg-white hover:bg-orange-300 cursor-pointer border-2 border-orange-300 px-3 py-1 w-44 text-center" @click="saveChanges">Lưu thay đổi ({{ countUserChanges }})</button>
      </div>
      <table class="w-full mt-10">
        <thead>
          <tr>
            <th>Tên</th>
            <th>Email</th>
            <th>Lớp</th>
            <th>N.Sinh</th>
            <th>Giới</th>
            <th>SĐT</th>
            <th>Mã</th>
            <th>ĐV</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><input placeholder="..." class="border-2 border-pink-300 px-2 py-0.5 w-full" v-model="filter.name"></td>
            <td><input placeholder="..." class="border-2 border-pink-300 px-2 py-0.5 w-full" v-model="filter.email"></td>
            <td><input placeholder="..." class="border-2 border-pink-300 px-2 py-0.5 w-full" v-model="filter.class"></td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td><button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-2 py-0.5 text-center" @click="search" v-if="!loadingUsers">Tìm</button></td>
          </tr>
          <tr class="border-b-2 border-b-slate-400">
            <td>{{ this.users.length }} thành viên</td>
            <td></td>
            <td></td>
            <td></td>
            <td>{{ this.users.filter(u => u.gender).length }} nữ</td>
            <td></td>
            <td></td>
            <td>{{ this.users.filter(u => u.certified).length }}/{{ this.users.length }}</td>
          </tr>
          <tr v-for="user in users">
            <td :class="{'text-red-500' : user.admin}">{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.class }}</td>
            <td>{{ new Intl.DateTimeFormat("vi-VN" , {dateStyle: "short"}).format(new Date(user.birth)) }}</td>
            <td>{{ user.gender ? "Nữ" : "Nam" }}</td>
            <td>{{ user.phone }}</td>
            <td>{{ user.sid }}</td>
            <td>
              <BadgeCheckIcon class="w-6" :class="user.certified ? 'text-sky-400' : 'text-gray-300'" @click="toggleCertified(user)"></BadgeCheckIcon>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="loadingUsers">
        <svg class="animate-spin h-16 w-16 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      <div class="mt-10" v-if="!userAvailable">Đã tải hết thành viên.</div>
    </div>
  </div>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2">
    <HomeIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="backToHome"></HomeIcon>
    <ChevronDoubleUpIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" @click="jumpToTop"></ChevronDoubleUpIcon>
  </div>
</template>

<script>
import {BadgeCheckIcon, ChevronDoubleUpIcon, HomeIcon} from '@heroicons/vue/solid'
import server from "../api/server";
import auth from "../api/auth";

export default {
  name: "UserManage",
  components: { ChevronDoubleUpIcon, HomeIcon, BadgeCheckIcon },
  data() {
    return {
      loadingUsers: false,
      userAvailable: true,
      users: [],
      dataOffset: 0,
      certChanges: {},
      filter: {
        name: "",
        email: "",
        class: ""
      }
    }
  },
  methods: {
    jumpToTop() {
      window.scrollTo(0, 0);
    },
    backToHome() {
      this.$router.push('/')
    },
    loadNextUsers(){
      this.loadingUsers = true
      server.loadUsers(50, this.dataOffset, this.filter, auth.getToken()).then(s => {
        if(s.users.length === 0) {
          this.userAvailable = false
        } else {
          this.dataOffset += s.users.length
        }
        this.users = this.users.concat(s.users)
        this.loadingUsers = false
      })
    },
    toggleCertified(user) {
      user.certified = !user.certified
      if(this.certChanges.hasOwnProperty(user.email)) {
        delete this.certChanges[user.email]
      } else {
        this.certChanges[user.email] = user.certified
      }
    },
    saveChanges() {
      if(this.countUserChanges === 0) {
        return
      }
      server.saveUserChanges({
        certified: this.certChanges
      }, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          window.location.reload();
        } else {
          alert(`Lỗi lưu thay đổi: ${s["error"]}`)
        }
      })
    },
    search() {
      this.userAvailable = true
      this.users = []
      this.loadingUsers = false
      this.dataOffset = 0
      this.certChanges = {}
      this.loadNextUsers()
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.loadingUsers && this.userAvailable) {
          this.loadNextUsers()
        }
      }
    }
  },
  computed: {
    countUserChanges() {
      return Object.keys(this.certChanges).length
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    if(!this.$root.isLoggedIn) {
      this.$router.push(`/`)
    }
    this.loadNextUsers()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
