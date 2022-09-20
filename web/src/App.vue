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
import conf from "./conf";
import {IsManager} from "./auth/roles";

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
          gender: "male",
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
  computed: {
    isManager() {
      return IsManager(this.user.profile.role)
    }
  },
  methods: {
    isLoggedIn(){
      return auth.isLoggedIn()
    },
    popupError(e){
      if(e instanceof ServerError){
        this.$notify({
          title: "Đã xảy ra lỗi!",
          text: e.message,
          type: "error"
        });
      } else if(e.hasOwnProperty("message")){
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
        return
      }
      if(res.profile.profileCover === "") {
        res.profile.profileCover = profileCoverDefaultImg
      } else {
        res.profile.profileCover = conf.assetURL + "/" + res.profile.profileCover
      }
      Object.assign(this.user, res)
      this.$forceUpdate()
      this.loadingProfile = false
    })
  }
}
</script>
