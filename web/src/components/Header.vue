<template>
  <div class="page-header relative shadow-2xl shadow-gray-400" style="background-image: url('https://i.imgur.com/qK6gzc0.jpg')">
    <div class="z-[5] relative text-white">
      <div class="w-full bg-indigo-500 text-right text-sm text-white p-3 pr-5 shadow-md shadow-slate-500 md:hidden">
        <button v-if="!$root.isLoggedIn()" @click="logIn()">Đăng nhập</button>
        <button v-else @click="viewProfile()">Trang cá nhân</button>
      </div>
      <div class="py-8 md:py-10 relative">
        <div class="max-w-[1024px] w-fit md:w-full m-auto relative">
          <div class="flex flex-row gap-3 place-items-center md:p-10">
            <img src="../assets/youth-logo.png" class="w-12 h-12"  alt=""/>
            <img src="../assets/das-logo.png" class="w-12 h-12"  alt=""/>
            <router-link to="/" class="text-4xl ml-5 font-light">
              <span class="block md:inline">ĐOÀN THPT</span> <span>DĨ AN</span>
              <p class="text-sm italic">Website đang hoạt động thử nghiệm</p>
            </router-link>
            <div class="absolute right-10 hidden md:block">
              <button v-if="!$root.isLoggedIn()" @click="logIn()" class="btn-primary">Đăng nhập</button>
              <button v-else @click="viewProfile()" class="btn-primary">Trang cá nhân</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import auth from "../api/auth";
import server from "../api/server";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "Header",
  methods: {
    postLogIn() {
      server.loadProfile('', auth.getToken()).then((s) => {
        if(s.hasOwnProperty("error")) {
          this.$notify({
            title: "Đăng nhập thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
          return
        }
        auth.setAuthenticated(true)
        window.location.reload();
      }, (e) => {
        this.$notify({
          title: "Đăng nhập thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    logIn() {
      auth.requestAuth((result) => {
        if(result === null){
          return
        }
        if(!result.user.email?.endsWith("@dian.sgdbinhduong.edu.vn")) {
          this.$notify({
            title: "Đăng nhập thất bại",
            text: "Tài khoản không thuộc nội bộ nhà trường",
            type: "error"
          });
        }
        setTimeout(this.postLogIn, 1000) // delay a little so #getToken can work later
      }, (e) => {
        this.$notify({
          title: "Đăng nhập thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    viewProfile() {
      this.$router.push("/u/" + this.$root.profile.email.substring(0, this.$root.profile.email.search("@"))).then(() => {
        this.$router.go(0)
      })
    }
  }
}
</script>

<style scoped>
.page-header:before {
  content: "";
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  left: 0;
  right: 0;
  background-image: linear-gradient(0deg, transparent, #222);
  z-index: 1;
}
</style>
