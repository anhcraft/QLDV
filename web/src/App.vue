<template>
  <router-view />
</template>

<script>
import auth from "./api/auth";
import server from "./api/server";

export default {
  data() {
    return {
      loadingProfile: false,
      progressionLoading: false,
      progressionLoaded: false,
      profile: {
        email: "",
        name: "Guest",
        certified: false,
        admin: false,
        mod: false,
        gender: false,
        class: "XX",
        entryDate: 2022,
        endDate: 2022,
        studentId: "0000000000000000",
        achievements: [],
        rates: {}
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
    isLoggedIn() {
      return auth.getToken() != null
    }
  },
  methods: {
    loadProgression() {
      if(this.isLoggedIn) {
        this.progressionLoading = true
        server.loadProgression(auth.getToken(), "").then(s => {
          if (s.hasOwnProperty("error")) {
            if(s["error"] === "ERR_TOKEN_VERIFY") {
              auth.destroySession()
            }
            return
          }
          s["achievements"].forEach((value) => {
            this.profile.achievements.push(value["title"] + ` (${value["year"]})`)
          });
          s["rates"].forEach((value) => {
            this.profile.rates[value["year"]] = value["level"]
          })
          this.progressionLoading = false
          this.progressionLoaded = true
        })
      }
    }
  },
  mounted() {
    if(this.isLoggedIn) {
      this.loadingProfile = true
      server.loadProfile(auth.getToken()).then(s => {
        if (s.hasOwnProperty("error")) {
          if(s["error"] === "ERR_TOKEN_VERIFY") {
            auth.destroySession()
          }
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
        this.loadingProfile = false
      })
    }
  }
}
</script>
