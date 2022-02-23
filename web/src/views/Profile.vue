<template>
  <Header></Header>
  <div class="pb-16 max-w-[1024px] m-auto">
    <Breadcrumb text="Trang cá nhân" :link="`/u/${this.$route.params.id}`"></Breadcrumb>
    <div class="grid grid-cols-3 gap-16 mt-10">
      <div class="col-span-1 shadow-lg shadow-slate-400 self-start">
        <LoadingState ref="profileLoadingState" hidden>
          <section>
            <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÔNG TIN</div>
            <div class="px-5 py-3">
              <p>Tên: {{ profile.name }}</p>
              <p>Lớp: {{ profile.class }}</p>
              <p>Giới tính: {{ profile.gender ? "Nữ" : "Nam" }}</p>
              <p>Đoàn viên: {{ profile.certified ? "Đã kết nạp" : "Không" }}</p>
              <p>Niên khóa: {{ profile.entryDate }} - {{ profile.endDate }}</p>
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
                  {{ value === 1 ? "Tốt" : (value === 2 ? "Khá" : "-") }} ({{ name }} - {{ parseInt(name) + 1 }})
                </li>
              </ul>
            </div>
          </section>
        </LoadingState>
      </div>
      <div class="col-span-2">
        <LoadingState ref="profileCoverLoadingState">
          <section class="w-full inline-block relative overflow-hidden border-4 border-dashed border-white" :class="{'hover:opacity-80 hover:border-black' : isPersonalProfile}">
            <div :style="{ 'background-image': 'url(' + profile.profileCover + ')' }" class="w-full h-64 bg-cover bg-center bg-no-repeat" />
            <input type="file" class="absolute left-0 top-0 opacity-0 h-64 w-full cursor-pointer" @change="onProfileCoverChange" accept="image/*" v-if="isPersonalProfile" />
          </section>
        </LoadingState>
        <LoadingState hidden ref="profileBoardLoadingState">
          <section class="mt-10">
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
                  menubar: false
                }"
                  v-model="profile.profileBoard"
              ></Editor>
              <button class="float-right bg-pink-400 hover:bg-pink-500 cursor-pointer px-3 py-1 text-white text-center text-sm mt-5" @click="saveBoard">Lưu lại</button>
            </div>
            <div v-else id="content" class="break-words" v-html="profile.profileBoard"></div>
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

export default {
  name: "Profile",
  components: {LoadingState, Breadcrumb, Header, FloatingMenu, Editor},
  data() {
    return {
      profile:  {
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
        profileCover: undefined,
        profileBoard: ""
      },
      achievements: [],
      rates: {}
    }
  },
  computed: {
    isPersonalProfile() {
      return this.$root.profile.email === this.getUserId()
    }
  },
  methods: {
    getUserId() {
      return this.$route.params.id + "@dian.sgdbinhduong.edu.vn";
    },
    onProfileCoverChange(e) {
      if (e.target.files.length > 0) {
        this.$refs.profileCoverLoadingState.activate()
        server.setProfileCover(e.target.files[0], auth.getToken()).then(s => {
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.$refs.profileCoverLoadingState.deactivate()
            window.location.reload()
          } else {
            alert(`Lỗi lưu ảnh bìa: ${s["error"]}`)
          }
        })
      }
    },
    saveBoard(){
      this.$refs.profileBoardLoadingState.activate()
      server.setProfileBoard(this.profile.profileBoard, auth.getToken()).then(s => {
        if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$refs.profileBoardLoadingState.deactivate()
        } else {
          alert(`Lỗi lưu tường nhà: ${s["error"]}`)
        }
      })
    }
  },
  mounted() {
    server.loadProfile(this.getUserId(), auth.getToken()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$router.push("/")
        return
      }
      this.profile = s;
      this.profile.entryDate = parseInt(s["entry"]);
      this.profile.endDate = this.profile.entryDate + 3;
      if (s["profileCover"].length > 0) {
        this.profile.profileCover = conf.server + "/static/" + s["profileCover"];
      } else {
        this.profile.profileCover = profileCoverDefaultImg
      }
      this.$refs.profileLoadingState.deactivate()
      this.$refs.profileCoverLoadingState.deactivate()
      this.$refs.profileBoardLoadingState.deactivate()
    })
    server.loadProgression(auth.getToken(), this.getUserId()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$router.push("/")
        return
      }
      s["achievements"].forEach((value) => {
        this.achievements.push(value["title"] + ` (${value["year"]})`)
      });
      s["rates"].forEach((value) => {
        this.rates[value["year"]] = value["level"]
      })
      this.$refs.progressionLoadingState.deactivate()
    })
  }
}
</script>

<style>
#content a {
  color: rgb(38 143 207);
}
#content ol {
  display: block;
  list-style-type: decimal;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  padding-inline-start: 40px;
}
#content ul {
  display: block;
  list-style-type: disc;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  padding-inline-start: 40px;
}
#content img, #content svg, #content video, #content canvas, #content audio, #content iframe, #content embed, #content object {
  display: inline;
  vertical-align: middle;
}
</style>
