<template>
  <Header></Header>
  <section class="page-section px-10 py-8 lg:py-16">
    <Breadcrumb text="Trang cá nhân" :link="{name: 'profile', params: { id: this.$route.params.id } }"></Breadcrumb>
    <LoadingState ref="profileLoadingState" hidden>
      <div class="grid grid-cols-1 md:grid-cols-3 md:gap-16 mt-10">
        <div class="col-span-1 self-start">

          <div class="shadow-lg shadow-slate-400">
            <section>
              <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÔNG TIN</div>
              <div class="px-5 py-3">
                <p>Tên: {{ user.profile.name }}</p>
                <p v-if="user.profile.hasOwnProperty('class')">Chi đoàn: {{ user.profile.class }}</p>
                <p v-if="user.profile.hasOwnProperty('gender')">Giới tính: {{ user.profile.gender === "female" ? "Nữ" : "Nam" }}</p>
                <p v-if="user.profile.hasOwnProperty('role')">Chức vụ: {{ roleName }}</p>
                <p v-if="user.profile.hasOwnProperty('entryYear')">Niên khóa: {{ user.profile.entryYear }} - {{ user.profile.entryYear + 3 }}</p>
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
                    {{ v.level === 1 ? "Xuất sắc" : (v.level === 2 ? "Khá" : (v.level === 3 ? "Trung bình" : "-")) }} ({{ v.year }} - {{ v.year + 1 }})
                  </li>
                </ul>
              </div>
            </section>
          </div>

          <div class="shadow-lg shadow-slate-400 mt-7" v-if="isPersonalProfile">
            <div class="border-l-4 border-l-sky-400 bg-sky-200 px-4 py-2 shadow-lg shadow-slate-300">THIẾT LẬP</div>
            <div class="p-5">
              <div class="flex flex-row place-items-center gap-1">
                <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.profileLocked" @input="user.profile.settings.profileLocked = !user.profile.settings.profileLocked" />
                <p>Khóa trang cá nhân</p>
              </div>
              <div class="flex flex-row place-items-center gap-1" :class="{'opacity-50' : user.profile.settings.profileLocked}">
                <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.classPublic" @input="user.profile.settings.classPublic = !user.profile.settings.classPublic" :disabled="user.profile.settings.profileLocked" />
                <p>Công khai chi đoàn</p>
              </div>
              <div class="flex flex-row place-items-center gap-1" :class="{'opacity-50' : user.profile.settings.profileLocked}">
                <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.achievementPublic" @input="user.profile.settings.achievementPublic = !user.profile.settings.achievementPublic" :disabled="user.profile.settings.profileLocked" />
                <p>Công khai thành tích</p>
              </div>
              <div class="flex flex-row place-items-center gap-1" :class="{'opacity-50' : user.profile.settings.profileLocked}">
                <input type="checkbox" class="w-4 h-4" v-bind:checked="user.profile.settings.annualRankPublic" @input="user.profile.settings.annualRankPublic = !user.profile.settings.annualRankPublic" :disabled="user.profile.settings.profileLocked" />
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
        </div>

        <div class="col-span-2 mt-10 md:mt-0">
          <section v-if="user.profile.hasOwnProperty('profileCover')" class="w-full inline-block relative overflow-hidden shadow-lg shadow-slate-400" :class="isPersonalProfile ? ('border-4 border-dashed border-white hover:border-black ' + (savingProfileCover ? 'opacity-50' : 'hover:opacity-80')) : ''">
            <div :style="{ 'background-image': 'url(' + user.profile.profileCover + ')' }" class="w-full h-64 bg-cover bg-center bg-no-repeat" />
            <input type="file" class="absolute left-0 top-0 opacity-0 h-64 w-full cursor-pointer" @change="onProfileCoverChange" accept="image/png, image/jpeg" v-if="isPersonalProfile" />
          </section>

          <section v-if="user.profile.hasOwnProperty('profileBoard')" class="mt-7 p-5 shadow-lg shadow-slate-400">
            <div v-if="isPersonalProfile">
              <Editor
                  apiKey="r7g4lphizuprqmrjv0ooj15pn5qpcesynrg101ekc40avzlg"
                  :init="{
                    height: 500,
                    plugins: ['advlist', 'autolink', 'lists', 'link', 'image', 'insertdatetime', 'media', 'table', 'wordcount'],
                    toolbar: 'undo redo | styleselect | bold italic | forecolor backcolor emoticons link | ' +
                      'bullist numlist outdent indent | media | ' +
                      'alignleft aligncenter alignright alignjustify | help',
                    menubar: false,
                    branding: false
                  }"
                  v-model="user.profile.profileBoard"
              ></Editor>
              <div class="flex place-content-end">
                <button class="btn-success mt-5" :class="{'opacity-50' : savingProfileBoard}" @click="saveBoard">Lưu thay đổi</button>
              </div>
            </div>
            <div v-else class="break-words prose w-full" v-html="user.profile.profileBoard"></div>
          </section>
        </div>
      </div>
    </LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import Editor from '@tinymce/tinymce-vue'
import profileCoverDefaultImg from "../assets/profile-cover.jpg";
import LoadingState from "../components/LoadingState.vue";
import {GetRoleName} from "../auth/roles";
import UserAPI from "../api/user-api";
import {ServerError} from "../api/server-error";
import conf from "../conf";
import Footer from "../components/Footer.vue";
import auth from "../auth/auth";

export default {
  name: "Profile",
  components: {Footer, LoadingState, Breadcrumb, Header, Editor},
  data() {
    return {
      savingProfileSettings: false,
      savingProfileBoard: false,
      savingProfileCover: false,
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
    roleName() {
      return GetRoleName(this.user.profile.role)
    },
    userId() {
      return this.$route.params.id.toString()
    },
    isPersonalProfile() {
      if(/^\d+$/.test(this.userId)) {
        return this.$root.user.profile.id === this.userId
      } else {
        return this.$root.user.profile.email === this.userId + "@dian.sgdbinhduong.edu.vn"
      }
    }
  },
  methods: {
    onProfileCoverChange(e) {
      if(this.savingProfileCover) return
      if (e.target.files.length > 0) {
        if(e.target.files[0].size > 500000){
          this.$root.popupError(new ServerError("ERROR_PROFILE_COVER_TOO_LARGE"))
          return
        }
        this.savingProfileCover = true
        UserAPI.uploadProfileCover(e.target.files[0]).then(res => {
          this.savingProfileCover = false
          if(res instanceof ServerError) {
            this.$root.popupError(res)
          } else {
            this.user.profile.profileCover = conf.assetURL + "/" + res.name
            this.$notify({
              title: "Đã lưu ảnh bìa",
              text: "",
              type: "success"
            });
          }
        })
      }
    },
    saveBoard(){
      if(this.savingProfileBoard) return
      if(this.user.profile.profileBoard.length < 10){
        this.$root.popupError(new ServerError("ERROR_PROFILE_BOARD_TOO_SHORT"))
        return
      }
      if(this.user.profile.profileBoard.length > 10000){
        this.$root.popupError(new ServerError("ERROR_PROFILE_BOARD_TOO_LONG"))
        return
      }
      this.savingProfileBoard = true
      UserAPI.updateUser("", {
        profile: {
          profileBoard: this.user.profile.profileBoard
        },
        achievements: undefined,
        annualRanks: undefined,
      }).then(res => {
        this.savingProfileBoard = false
        if(res instanceof ServerError) {
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
    saveSettings() {
      if(this.savingProfileSettings) return
      this.savingProfileSettings = true
      UserAPI.updateUser("", {
        profile: {
          settings: this.user.profile.settings
        },
        achievements: undefined,
        annualRanks: undefined,
      }).then(res => {
        this.savingProfileSettings = false
        if(res instanceof ServerError) {
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
  },
  mounted() {
    UserAPI.getUser(this.userId, {
      profile: true,
      achievements: true,
      "annual-ranks": true
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
      this.$refs.profileLoadingState.deactivate()
    })
  }
}
</script>
