<template>
  <div class="fixed right-10 bottom-10 flex flex-col gap-2" v-if="!$root.loadingProfile">
    <router-link to="/pm" v-if="$root.isLoggedIn && $root.profile.admin">
      <NewspaperIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2"></NewspaperIcon>
    </router-link>
    <router-link to="/um" v-if="$root.isLoggedIn && ($root.profile.admin || $root.profile['mod'])">
      <UsersIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2"></UsersIcon>
    </router-link>
    <!--<CogIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2"></CogIcon>-->
    <LogoutIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2" v-if="$root.isLoggedIn" @click="logOut()"></LogoutIcon>
  </div>
</template>

<script>
import {
  CogIcon,
  LogoutIcon,
  NewspaperIcon,
  UsersIcon
} from "@heroicons/vue/solid";
import auth from "../api/auth";

export default {
  name: "FloatingMenu",
  components: {
    LogoutIcon,
    CogIcon,
    NewspaperIcon,
    UsersIcon
  },
  props: {
    pageNavigation: Boolean
  },
  methods: {
    logOut() {
      auth.destroySession()
      this.$router.push("/")
    }
  }
}
</script>

