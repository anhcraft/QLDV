<template>
  <Header></Header>
  <div class="pb-16 max-w-[1024px] m-auto p-5 md:px-10">
    <Breadcrumb text="Trang cá nhân" :link="`/u/${this.$route.params.id}`"></Breadcrumb>
    <div class="grid grid-cols-1 md:grid-cols-3 md:gap-16 mt-10">
      <div class="col-span-1 self-start">
        <div class="shadow-lg shadow-slate-400">
          <LoadingState ref="profileLoadingState" hidden>
            <section>
              <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÔNG TIN</div>
              <div class="px-5 py-3">
                <p>Tên: {{ profile.name }}</p>
                <p v-if="profile.hasOwnProperty('class')">Chi đoàn: {{ profile.class }}</p>
                <p v-if="profile.hasOwnProperty('gender')">Giới tính: {{ profile.gender ? "Nữ" : "Nam" }}</p>
                <p v-if="profile.hasOwnProperty('certified')">Đoàn viên: {{ profile.certified ? "Đã kết nạp" : "Không" }}</p>
                <p v-if="profile.hasOwnProperty('entry')">Niên khóa: {{ profile.entryDate }} - {{ profile.endDate }}</p>
              </div>
            </section>
          </LoadingState>
          <LoadingState ref="progressionLoadingState" hidden>
            <section v-if="achievements.length > 0">
              <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÀNH TÍCH</div>
              <div class="px-5 py-3">
                <ul class="list-disc list-inside">
                  <li v-for="a in achievements">{{ a }}</li>
                </ul>
              </div>
            </section>
            <section v-if="Object.keys(rates).length > 0">
              <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">XẾP LOẠI</div>
              <div class="px-5 py-3">
                <ul class="list-disc list-inside">
                  <li v-for="(value, name) in rates">
                    {{ value === 1 ? "Xuất sắc" : (value === 2 ? "Khá" : (value === 3 ? "Trung bình" : "-")) }} ({{ name }} - {{ parseInt(name) + 1 }})
                  </li>
                </ul>
              </div>
            </section>
          </LoadingState>
        </div>
        <div class="shadow-lg shadow-slate-400 mt-7" v-if="isPersonalProfile">
          <div class="border-l-4 border-l-sky-400 bg-sky-200 px-4 py-2 shadow-lg shadow-slate-300">THIẾT LẬP</div>
          <div class="p-5">
            <div class="flex flex-row place-items-center gap-1">
              <input type="checkbox" class="w-4 h-4" v-bind:checked="isProfileLocked" @input="toggleProfileSetting(1, $event.target.value)" />
              <p>Khóa trang cá nhân</p>
            </div>
            <div class="flex flex-row place-items-center gap-1">
              <input type="checkbox" class="w-4 h-4" v-bind:checked="isClassPublic" @input="toggleProfileSetting(2, $event.target.value)" />
              <p>Công khai chi đoàn</p>
            </div>
            <div class="flex flex-row place-items-center gap-1">
              <input type="checkbox" class="w-4 h-4" v-bind:checked="isAchievementPublic" @input="toggleProfileSetting(4, $event.target.value)" />
              <p>Công khai thành tích</p>
            </div>
            <div class="flex flex-row place-items-center gap-1">
              <input type="checkbox" class="w-4 h-4" v-bind:checked="isRatePublic" @input="toggleProfileSetting(8, $event.target.value)" />
              <p>Công khai xếp loại</p>
            </div>
            <button class="btn-success mt-3" :class="{'opacity-50' : savingProfileSettings}" @click="saveSettings">Lưu thay đổi</button>
          </div>
        </div>
      </div>
      <div class="col-span-2 mt-10 md:mt-0">
        <LoadingState ref="profileCoverLoadingState">
          <section v-if="profile.hasOwnProperty('profileCover')" class="w-full inline-block relative overflow-hidden shadow-lg shadow-slate-400" :class="{'border-4 border-dashed border-white hover:opacity-80 hover:border-black' : isPersonalProfile}">
            <div :style="{ 'background-image': 'url(' + profile.profileCover + ')' }" class="w-full h-64 bg-cover bg-center bg-no-repeat" />
            <input type="file" class="absolute left-0 top-0 opacity-0 h-64 w-full cursor-pointer" @change="onProfileCoverChange" accept="image/*" v-if="isPersonalProfile" />
          </section>
        </LoadingState>
        <LoadingState hidden ref="profileBoardLoadingState">
          <section v-if="profile.hasOwnProperty('profileBoard')" class="mt-7 p-5 shadow-lg shadow-slate-400">
            <div v-if="isPersonalProfile">
              <Editor
                  apiKey="r7g4lphizuprqmrjv0ooj15pn5qpcesynrg101ekc40avzlg"
                  :init="{
                  height: 500,
                  plugins: [
                    'advlist autolink link image lists charmap print preview hr anchor pagebreak',
                    'searchreplace wordcount visualblocks visualchars code fullscreen insertdatetime media nonbreaking',
                    'table emoticons template paste help'
                  ],
                  toolbar: 'undo redo | styleselect | bold italic | forecolor backcolor emoticons link | ' +
                    'bullist numlist outdent indent | media | ' +
                    'alignleft aligncenter alignright alignjustify | help',
                  menubar: false,
                  branding: false
                }"
                  v-model="profile.profileBoard"
              ></Editor>
              <div class="flex place-content-end">
                <button class="btn-success mt-5" @click="saveBoard">Lưu thay đổi</button>
              </div>
            </div>
            <div v-else class="break-words prose w-full" v-html="profile.profileBoard"></div>
          </section>
        </LoadingState>
      </div>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import server from "../api/server";
import auth from "../api/auth";
import Editor from '@tinymce/tinymce-vue'
import conf from "../conf";
import profileCoverDefaultImg from "../assets/profile-cover.jpg";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "Profile",
  components: {LoadingState, Breadcrumb, Header, FloatingMenu, Editor},
  data() {
    return {
      profile:  {
        email: "",
        name: "",
        admin: false,
        mod: false,
        profileSettings: 0
      },
      achievements: [],
      rates: {},
      savingProfileSettings: false
    }
  },
  computed: {
    isPersonalProfile() {
      return this.$root.profile.email === this.getUserId()
    },
    isProfileLocked() {
      return (this.profile.profileSettings&1) === 1
    },
    isClassPublic() {
      return (this.profile.profileSettings&2) === 2
    },
    isAchievementPublic() {
      return (this.profile.profileSettings&4) === 4
    },
    isRatePublic() {
      return (this.profile.profileSettings&8) === 8
    }
  },
  methods: {
    getUserId() {
      return this.$route.params.id + "@dian.sgdbinhduong.edu.vn";
    },
    toggleProfileSetting(a, b) {
      this.profile.profileSettings = b ? (this.profile.profileSettings ^ a) : (this.profile.profileSettings | a)
    },
    onProfileCoverChange(e) {
      if (e.target.files.length > 0) {
        if(e.target.files[0].size > 500000){
          this.$notify({
            title: "Lưu ảnh bìa thất bại",
            text: lookupErrorCode("ERR_PROFILE_COVER_SIZE_EXCEEDED"),
            type: "error"
          });
          return
        }
        this.$refs.profileCoverLoadingState.activate()
        server.setProfileCover(e.target.files[0], auth.getToken()).then(s => {
          this.$refs.profileCoverLoadingState.deactivate()
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.$notify({
              title: "Đã lưu ảnh bìa",
              text: "",
              type: "success"
            });
            window.location.reload() // cover image wont be changed without page-reloading
          } else {
            this.$notify({
              title: "Lưu ảnh bìa thất bại",
              text: lookupErrorCode(s["error"]),
              type: "error"
            });
          }
        }, (e) => {
          this.$notify({
            title: "Lưu ảnh bìa thất bại",
            text: e.message,
            type: "error"
          });
        });
      }
    },
    saveBoard(){
      this.$refs.profileBoardLoadingState.activate()
      server.setProfileBoard(this.profile.profileBoard, auth.getToken()).then(s => {
        this.$refs.profileBoardLoadingState.deactivate()
        if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$notify({
            title: "Đã lưu tường nhà",
            text: "",
            type: "success"
          });
          return
        }
        this.$notify({
          title: "Lưu tường nhà thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
      }, (e) => {
        this.$notify({
          title: "Lưu tường nhà thất bại",
          text: e.message,
          type: "error"
        });
      });
    },
    saveSettings() {
      if(this.savingProfileSettings) return
      this.savingProfileSettings = true
      server.setProfileSettings(this.profile.profileSettings, auth.getToken()).then(s => {
        this.savingProfileSettings = false
        if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$notify({
            title: "Đã lưu thiết lập thành công",
            text: "",
            type: "success"
          });
          return
        }
        this.$notify({
          title: "Lưu thiết lập thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
      }, (e) => {
        this.$notify({
          title: "Lưu thiết lập thất bại",
          text: e.message,
          type: "error"
        });
      });
    }
  },
  mounted() {
    server.loadProfile(this.getUserId(), auth.getToken()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$notify({
          title: "Tải thông tin thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
        return
      }
      this.profile = s;
      if(s.hasOwnProperty("entry")) {
        this.profile.entryDate = parseInt(s["entry"]);
        this.profile.endDate = this.profile.entryDate + 3;
      }
      if(s.hasOwnProperty("profileCover")) {
        if (s["profileCover"].length > 0) {
          this.profile.profileCover = conf.server + "/static/" + s["profileCover"];
        } else {
          this.profile.profileCover = profileCoverDefaultImg
        }
      }
      this.$refs.profileLoadingState.deactivate()
      this.$refs.profileCoverLoadingState.deactivate()
      this.$refs.profileBoardLoadingState.deactivate()
    }, (e) => {
      this.$notify({
        title: "Tải thông tin thất bại",
        text: e.message,
        type: "error"
      });
    });
    server.loadProgression(auth.getToken(), this.getUserId()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$notify({
          title: "Tải thông tin thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
        return
      }
      s["achievements"].forEach((value) => {
        this.achievements.push(value["title"] + ` (${value["year"]})`)
      });
      s["rates"].forEach((value) => {
        this.rates[value["year"]] = value["level"]
      })
      this.$refs.progressionLoadingState.deactivate()
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
