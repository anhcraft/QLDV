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
import server from "./api/server";
import profileCoverDefaultImg from "./assets/profile-cover.jpg"
import conf from "./conf";
import lookupErrorCode from "./api/errorCode";

export default {
  data() {
    return {
      loadingProfile: false,
      progressionLoading: false,
      progressionLoaded: false,
      profile: {
        email: "",
        name: "",
        certified: false,
        admin: false,
        mod: false,
        gender: false,
        class: "",
        entryDate: 2022,
        endDate: 2022,
        studentId: "0000000000000000",
        achievements: [],
        rates: {},
        profileCover: undefined,
        profileBoard: ""
      }
    }
  },
  computed: {
    formattedStudentId() {
      const chunks = [];
      for (let i = 0, charsLength = this.profile.studentId.length; i < charsLength; i += 4) {
        chunks.push(this.profile.studentId.substring(i, i + 4));
      }
      return chunks.join(" ");
    },
  },
  methods: {
    isLoggedIn() {
      return auth.isLoggedIn()
    }
  },
  mounted() {
    if(!this.isLoggedIn()) {
      return
    }
    this.loadingProfile = true
    server.loadProfile('', auth.getToken()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$notify({
          title: "Tải thông tin thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
        auth.setAuthenticated(false)
        window.location.reload();
        return
      }
      this.profile.email = s["email"];
      this.profile.name = s["name"];
      this.profile.certified = s["certified"];
      this.profile.admin = s["admin"];
      this.profile.mod = s["mod"];
      this.profile.gender = s["gender"];
      this.profile.entryDate = parseInt(s["entry"]);
      this.profile.endDate = this.profile.entryDate + 3;
      this.profile.class = s["class"];
      this.profile.studentId = s["sid"];
      if (s["profileCover"].length > 0) {
        this.profile.profileCover = conf.server + "/static/" + s["profileCover"];
      } else {
        this.profile.profileCover = profileCoverDefaultImg
      }
      this.profile.profileBoard = s["profileBoard"]
      this.loadingProfile = false
    }, (e) => {
      this.$notify({
        title: "Tải thông tin thất bại",
        text: e.message,
        type: "error"
      });
    });
  }
}
</script>
