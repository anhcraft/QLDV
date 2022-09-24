<template>
  <notifications/>
  <router-view/>
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
import profileCoverDefaultImg from "./assets/profile-cover.webp"
import UserAPI from "./api/user-api";
import {ServerError} from "./api/server-error";
import conf from "./conf";
import {GetRoleGroup, IsManager, RoleGroupGlobalManager} from "./auth/roles";
import {nextTick} from "vue";

export default {
  data() {
    return {
      loadingProfile: false,
      initialized: false,
      initQueue: [],
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
      },
      mountQueue: []
    }
  },
  computed: {
    isManager() {
      return IsManager(this.user.profile.role)
    },
    isGlobalManager() {
      return GetRoleGroup(this.user.profile.role) >= RoleGroupGlobalManager
    }
  },
  methods: {
    isLoggedIn() {
      return auth.isLoggedIn()
    },
    popupError(e) {
      if (e instanceof ServerError) {
        this.$notify({
          title: "Đã xảy ra lỗi!",
          text: e.message,
          type: "error"
        });
      } else if (e.hasOwnProperty("message")) {
        this.$notify({
          title: "Đã xảy ra lỗi!",
          text: e["message"],
          type: "error"
        });
      }
    },
    pullQueue() {
      this.initQueue.forEach(v => v.call(null))
    },
    pushQueue(func) {
      if (this.initialized) {
        func.call(null)
      } else {
        this.initQueue.push(func)
      }
    }
  },
  async mounted() {
    await nextTick()
    if (!this.isLoggedIn()) {
      this.initialized = true
      this.pullQueue()
      return
    }
    this.loadingProfile = true
    UserAPI.getUser("", {
      profile: true,
      achievements: false,
      "annual-ranks": false
    }).then((res) => {
      this.loadingProfile = false
      if (res instanceof ServerError) {
        this.$root.popupError(res)
        this.initialized = true
        this.pullQueue()
        return
      }
      if (res.profile.profileCover === "") {
        res.profile.profileCover = profileCoverDefaultImg
      } else {
        res.profile.profileCover = conf.assetURL + "/" + res.profile.profileCover
      }
      Object.assign(this.user, res)
      this.$forceUpdate()
      this.initialized = true
      this.pullQueue()
    })
  }
}
</script>
