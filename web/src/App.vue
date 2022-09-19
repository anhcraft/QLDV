<template>
  <notifications />
  <router-view />
</template>

<style>
.vue-notification {
  font-size: 0.9rem;
}
.prose img {
  display: inline;
}
</style>

<script>
import auth from "./auth/auth";
import profileCoverDefaultImg from "./assets/profile-cover.jpg"
import UserAPI from "./api/user-api";
import {ServerError} from "./api/server-error";

export default {
  data() {
    return {
      loadingProfile: false,
      user: {
        profile: {
          id: 0,
          email: "",
          role: 0,
          name: "",
          gender: false,
          birthday: 0,
          entryYear: 0,
          phone: "",
          class: "",
          settings: {
            profileLocked: false,
            classPublic: false,
            achievementPublic: false,
            annualRankPublic: false
          },
          profileCover: "",
          profileBoard: "",
          updateDate: 0,
          createDate: 0
        },
        achievements: [],
        annualRanks: []
      }
    }
  },
  methods: {
    isLoggedIn(){
      return auth.isLoggedIn()
    },
    popupError(e){
      if(e.hasOwnProperty("message")){
        this.$notify({
          title: "Đã xảy ra lỗi!",
          text: e["message"],
          type: "error"
        });
      }
    }
  },
  mounted() {
    if(!this.isLoggedIn()) return
    this.loadingProfile = true
    UserAPI.getUser("", {
      profile: true,
      achievements: false,
      "annual-ranks": false
    }).then((res) => {
      if(res instanceof ServerError) {
        this.$root.popupError(res)
      }
      if(res.profile.profileCover === undefined) {
        res.profile.profileCover = profileCoverDefaultImg
      }
      Object.assign(this.user, res)
      this.$forceUpdate()
      this.loadingProfile = false
    })
  }
}
</script>
