<template>
  <div class="fixed z-[99] right-5 lg:right-10 bottom-10 flex flex-col gap-2" v-if="!$root.loadingProfile">
    <router-link to="/em" v-if="$root.isLoggedIn() && $root.profile.admin">
      <CalendarIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2 bg-white"></CalendarIcon>
    </router-link>
    <router-link to="/pm" v-if="$root.isLoggedIn() && $root.profile.admin">
      <NewspaperIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2 bg-white"></NewspaperIcon>
    </router-link>
    <router-link to="/um" v-if="$root.isLoggedIn() && ($root.profile.admin || $root.profile['mod'])">
      <UsersIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2 bg-white"></UsersIcon>
    </router-link>
    <!--<CogIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2"></CogIcon>-->
    <ArrowRightOnRectangleIcon class="w-12 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-2 bg-white" v-if="$root.isLoggedIn()" @click="logOut()"></ArrowRightOnRectangleIcon>
  </div>
</template>

<script>
import {
  CalendarIcon,
  CogIcon,
  ArrowRightOnRectangleIcon,
  NewspaperIcon,
  UsersIcon
} from '@heroicons/vue/24/solid';
import auth from "../api/auth";

export default {
  name: "FloatingMenu",
  components: {
    ArrowRightOnRectangleIcon,
    CogIcon,
    NewspaperIcon,
    UsersIcon,
    CalendarIcon
  },
  props: {
    pageNavigation: Boolean
  },
  methods: {
    logOut() {
      auth.setAuthenticated(false)
      window.location.reload()
    }
  }
}
</script>

