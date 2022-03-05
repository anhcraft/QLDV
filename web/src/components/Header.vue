<template>
  <div class="w-full bg-indigo-500 text-right text-sm text-white p-3 pr-5 shadow-md shadow-slate-500 md:hidden">
    <button v-if="!$root.isLoggedIn()" @click="logIn()">Đăng nhập</button>
    <button v-else @click="viewProfile()">Trang cá nhân</button>
  </div>
  <div class="py-8 md:py-10 relative">
    <div class="max-w-[1024px] w-fit md:w-full m-auto relative">
      <div class="flex flex-row gap-3 place-items-center md:p-10">
        <img src="../assets/youth-logo.png" class="w-12 h-12"  alt=""/>
        <img src="../assets/das-logo.png" class="w-12 h-12"  alt=""/>
        <router-link to="/" class="text-3xl ml-5 font-light border-b-2 border-b-white hover:border-b-gray-500">
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
</template>

<script>
import {getAuth, GoogleAuthProvider, signInWithPopup} from "firebase/auth";
import server from "../api/server";
import Cookies from "js-cookie";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "Header",
  methods: {
    logIn() {
      const provider = new GoogleAuthProvider();
      provider.addScope('https://www.googleapis.com/auth/userinfo.email');
      const auth = getAuth();
      signInWithPopup(auth, provider)
          .then((result) => {
            const credential = GoogleAuthProvider.credentialFromResult(result);
            if(credential != null){
              if(result.user.email?.endsWith("@dian.sgdbinhduong.edu.vn")) {
                auth.currentUser?.getIdToken().then(token => {
                  server.loadProfile('', token).then((s) => {
                    if(s.hasOwnProperty("error")) {
                      this.$notify({
                        title: "Đăng nhập thất bại",
                        text: lookupErrorCode(s["error"]),
                        type: "error"
                      });
                    } else {
                      Cookies.set('qldvtkn', token, {expires: 3})
                      window.location.reload();
                    }
                  }, (e) => {
                    this.$notify({
                      title: "Đăng nhập thất bại",
                      text: e.message,
                      type: "error"
                    });
                  })
                })
              } else {
                this.$notify({
                  title: "Đăng nhập thất bại",
                  text: "Tài khoản không thuộc nội bộ nhà trường",
                  type: "error"
                });
              }
            } else {
              this.$notify({
                title: "Đăng nhập thất bại",
                text: "Lỗi hệ thống xác thực",
                type: "error"
              });
            }
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
