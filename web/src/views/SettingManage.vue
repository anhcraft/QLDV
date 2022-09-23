<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <LoadingState ref="homepageLoadingState">
      <div class="text-2xl">Trang chủ</div>
      <div class="flex flex-col gap-5 mt-5">
        <div class="centered-horizontal">
          <p>Số lượng đoàn viên tiêu biểu</p>
          <input type="number" class="border border-slate-400 ml-5" v-model.number="homepageSettings.featuredUserLimit" min="0" max="10">
        </div>
        <div class="centered-horizontal">
          <p>Số lượng thành tựu tiêu biểu sẽ hiển thị</p>
          <input type="number" class="border border-slate-400 ml-5" v-model.number="homepageSettings.featuredAchievementLimit" min="0" max="10">
        </div>
        <div>
          <div class="centered-horizontal">
            <p>Link ảnh hoạt động</p>
            <PlusCircleIcon class="w-6 h-6 text-slate-500 cursor-pointer" @click="homepageSettings.activitySlideshow.push('')"></PlusCircleIcon>
          </div>
          <ul class="list-disc list-inside">
            <li class="mt-1" v-for="(val, index) in homepageSettings.activitySlideshow">
              <input type="text" class="border border-slate-400 px-2 py-1 rounded-md w-[400px]" v-model.trim="homepageSettings.activitySlideshow[index]">
            </li>
          </ul>
        </div>
        <div>
          <button class="btn-success" :class="{'opacity-50' : updatingSettings}" @click="saveHomepageSettings">Cập nhật</button>
        </div>
      </div>
    </LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import LoadingState from "../components/LoadingState.vue";
import Footer from "../components/Footer.vue";
import {PlusCircleIcon} from "@heroicons/vue/24/solid";
import SettingAPI from "../api/setting-api";
import {ServerError} from "../api/server-error";

export default {
  name: "SettingManage",
  components: {
    LoadingState, Header, Footer, PlusCircleIcon
  },
  data() {
    return {
      homepageSettings: {
        activitySlideshow: [],
        featuredUserLimit: 0,
        featuredAchievementLimit: 0
      },
      updatingSettings: false
    }
  },
  methods: {
    saveHomepageSettings() {
      if(this.updatingSettings) return
      this.updatingSettings = true
      SettingAPI.updateSetting("homepage", this.homepageSettings).then(res => {
        this.updatingSettings = false
        if(res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        this.$notify({
          title: "Đã lưu thay đổi",
          text: "",
          type: "success"
        });
      })
    }
  },
  mounted() {
    const f = () => {
      if(!this.$root.isLoggedIn() || !this.$root.isGlobalManager) {
        this.$router.push({name: "home"})
        return
      }
      SettingAPI.getSetting("homepage").then(data => {
        if(data instanceof ServerError) {
          this.$root.popupError(data)
        } else {
          if (!data.hasOwnProperty("activitySlideshow")) {
            data["activitySlideshow"] = []
          }
          this.homepageSettings = data
        }
        this.$refs.homepageLoadingState.deactivate()
      })
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
