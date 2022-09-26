<template>
  <div class="shadow-lg shadow-slate-400">
    <section class="centered-horizontal">
      <div class="relative w-full w-32 h-32">
        <img :src="user.profile.profileAvatar" class="
            md:absolute top-[-15px] left-[-15px]
            w-32 h-32 rounded-full shadow-md shadow-slate-400
            border-4 border-dashed border-white bg-white"
             :class="isPersonalProfile ? ('cursor-pointer hover:border-black ' + (savingProfileAvatar ? 'opacity-50' : 'hover:opacity-80')) : ''"
             @click="openProfileAvatarEditor()">
      </div>
      <div class="text-sm ml-5 md:ml-0">
        <p class="text-xl mb-1">
          <span v-if="!isPersonalProfile && user.profile.settings.profileLocked">???</span>
          <span v-else-if="user.profile.hasOwnProperty('name')">{{ user.profile.name }}</span>
        </p>
        <p v-if="user.profile.hasOwnProperty('class')">Chi đoàn: {{ user.profile.class }}</p>
        <p v-if="user.profile.hasOwnProperty('role')">Chức vụ: {{ roleName }}</p>
        <p v-if="user.profile.hasOwnProperty('entryYear')">Niên khóa: {{ user.profile.entryYear }} -
          {{ user.profile.entryYear + 3 }}</p>
      </div>
    </section>
    <section v-if="user.achievements.length > 0">
      <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÀNH TÍCH</div>
      <div class="px-5 py-3">
        <ul class="list-disc list-inside">
          <li v-for="v in user.achievements">{{ v.title }} ({{ v.year }})</li>
        </ul>
      </div>
    </section>
    <section v-if="user.annualRanks.length > 0">
      <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">XẾP LOẠI</div>
      <div class="px-5 py-3">
        <ul class="list-disc list-inside">
          <li v-for="v in user.annualRanks">
            {{ v.level === 1 ? "Xuất sắc" : (v.level === 2 ? "Khá" : (v.level === 3 ? "Trung bình" : "-")) }}
            ({{ v.year }} - {{ v.year + 1 }})
          </li>
        </ul>
      </div>
    </section>
  </div>

  <div class="shadow-lg shadow-slate-400 mt-7" v-if="isPersonalProfile">
    <div class="border-l-4 border-l-sky-400 bg-sky-200 px-4 py-2 shadow-lg shadow-slate-300">THIẾT LẬP</div>
    <div class="p-5">
      <div class="flex flex-row place-items-center gap-1">
        <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.profileLocked"
               @input="user.profile.settings.profileLocked = !user.profile.settings.profileLocked"/>
        <p>Khóa trang cá nhân</p>
      </div>
      <div class="flex flex-row place-items-center gap-1"
           :class="{'opacity-50' : user.profile.settings.profileLocked}">
        <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.classPublic"
               @input="user.profile.settings.classPublic = !user.profile.settings.classPublic"
               :disabled="user.profile.settings.profileLocked"/>
        <p>Công khai chi đoàn</p>
      </div>
      <div class="flex flex-row place-items-center gap-1"
           :class="{'opacity-50' : user.profile.settings.profileLocked}">
        <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.achievementPublic"
               @input="user.profile.settings.achievementPublic = !user.profile.settings.achievementPublic"
               :disabled="user.profile.settings.profileLocked"/>
        <p>Công khai thành tích</p>
      </div>
      <div class="flex flex-row place-items-center gap-1"
           :class="{'opacity-50' : user.profile.settings.profileLocked}">
        <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.annualRankPublic"
               @input="user.profile.settings.annualRankPublic = !user.profile.settings.annualRankPublic"
               :disabled="user.profile.settings.profileLocked"/>
        <p>Công khai xếp loại</p>
      </div>
      <button class="btn-success mt-3" :class="{'opacity-50' : savingProfileSettings}" @click="saveSettings">Lưu thay đổi</button>
    </div>
  </div>

  <div class="shadow-lg shadow-slate-400 mt-7" v-if="isPersonalProfile">
    <div class="border-l-4 border-l-gray-400 bg-gray-200 px-4 py-2 shadow-lg shadow-slate-300">TÀI KHOẢN</div>
    <div class="p-5">
      <button class="btn-outline-sm ml-auto" @click="logout">Đăng xuất</button>
    </div>
  </div>

  <div v-if="showProfileAvatarEditor">
    <div class="z-[100] fixed w-full h-full top-0 left-0 bg-black opacity-50"
         @click="showProfileAvatarEditor = false"></div>
    <div class="z-[200] fixed-center bg-white px-5 pt-5 max-w-max">
      <ImgCutter
          ref="imgCutterModal"
          :crossOrigin="true"
          crossOriginHeader="*"
          rate=""
          toolBgc="none"
          :isModal="false"
          :showChooseBtn="true"
          :lockScroll="true"
          :boxWidth="800"
          :boxHeight="500"
          :cutWidth="512"
          :cutHeight="512"
          :DoNotDisplayCopyright="true"
          :sizeChange="true"
          :moveAble="true"
          :imgMove="true"
          :originalGraph="false"
          :smallToUpload="true"
          :saveCutPosition="true"
          :scaleAble="false"
          :previewMode="true"
          :quality="1"
          :toolBoxOverflow="true"
          @cutDown="cutDownProfileAvatar"></ImgCutter>
    </div>
  </div>
</template>

<script>
import {GetRoleName} from "../../auth/roles";
import UserAPI from "../../api/user-api";
import {ServerError} from "../../api/server-error";
import auth from "../../auth/auth";
import ImgCutter from 'vue-img-cutter'

export default {
  name: "ProfileSidebar",
  props: {
    user: Object,
    isPersonalProfile: Boolean
  },
  components: {
    ImgCutter
  },
  data() {
    return {
      savingProfileSettings: false,
      showProfileAvatarEditor: false,
      savingProfileAvatar: false
    }
  },
  computed: {
    roleName() {
      return GetRoleName(this.user.profile.role)
    }
  },
  methods: {
    openProfileAvatarEditor() {
      if(!this.isPersonalProfile) return
      this.showProfileAvatarEditor = true
      this.$nextTick(() => {
        this.$refs.imgCutterModal.handleOpen({
          name: new URL(this.user.profile.profileAvatar).pathname.split("/").pop(),
          src: this.user.profile.profileAvatar,
        });
      })
    },
    cutDownProfileAvatar(e){
      if (!this.isPersonalProfile || this.savingProfileAvatar) return
      if (e.blob.size > 1000000) {
        this.$root.popupError(new ServerError("ERROR_PROFILE_AVATAR_TOO_LARGE"))
        return
      }
      this.savingProfileAvatar = true
      UserAPI.uploadProfileAvatar(e.blob).then(res => {
        this.savingProfileAvatar = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          window.location.reload()
        }
      })
    },
    saveSettings() {
      if (this.savingProfileSettings) return
      this.savingProfileSettings = true
      UserAPI.updateUser("", {
        profile: {
          settings: this.user.profile.settings
        },
        achievements: undefined,
        annualRanks: undefined,
      }).then(res => {
        this.savingProfileSettings = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          this.$notify({
            title: "Đã lưu thay đổi",
            text: "",
            type: "success"
          });
        }
      })
    },
    logout() {
      auth.logout().then(() => {
        this.$router.push({name: "home"}).then(() => {
          window.location.reload()
        })
      })
    }
  }
}
</script>
