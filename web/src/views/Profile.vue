<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <LoadingState ref="profileLoadingState" hidden>
      <div class="grid grid-cols-1 md:grid-cols-3 md:gap-8 lg:gap-10 xl:gap-16">
        <div class="col-span-1 self-start">
          <ProfileSidebar :user="user" :is-personal-profile="isPersonalProfile"></ProfileSidebar>
        </div>

        <div class="col-span-2 mt-10 md:mt-0">
          <section v-if="user.profile.hasOwnProperty('profileCover')"
                   class="w-full inline-block relative overflow-hidden shadow-lg shadow-slate-400"
                   :class="isPersonalProfile ? ('cursor-pointer border-4 border-dashed border-white hover:border-black ' + (savingProfileCover ? 'opacity-50' : 'hover:opacity-80')) : ''">
            <div :style="{ 'background-image': 'url(' + user.profile.profileCover + ')' }"
                 @click="openProfileCoverEditor()" class="w-full h-64 bg-cover bg-center bg-no-repeat"/>
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
                <button class="btn-success mt-5" :class="{'opacity-50' : savingProfileBoard}" @click="saveBoard">Lưu
                  thay đổi
                </button>
              </div>
            </div>
            <div v-else class="break-words prose w-full" v-html="user.profile.profileBoard"></div>
          </section>

        </div>
      </div>
    </LoadingState>
  </section>
  <Footer></Footer>

  <div v-if="showProfileCoverEditor">
    <div class="z-[100] fixed w-full h-full top-0 left-0 bg-black opacity-50"
         @click="showProfileCoverEditor = false"></div>
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
          :cutWidth="600"
          :cutHeight="256"
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
          @cutDown="cutDownProfileCover"></ImgCutter>
    </div>
  </div>

</template>

<script>
import Header from "../components/Header.vue";
import Editor from '@tinymce/tinymce-vue'
import profileCoverDefaultImg from "../assets/profile-cover.webp";
import profileFemaleAvatarDefaultImg from "../assets/avatar-female.webp";
import profileMaleAvatarDefaultImg from "../assets/avatar-male.webp";
import LoadingState from "../components/LoadingState.vue";
import UserAPI from "../api/user-api";
import {ServerError} from "../api/server-error";
import conf from "../conf";
import Footer from "../components/Footer.vue";
import ImgCutter from 'vue-img-cutter'
import ProfileSidebar from "../components/profile/ProfileSidebar.vue";

export default {
  name: "Profile",
  components: {ProfileSidebar, Footer, LoadingState, Header, Editor, ImgCutter},
  data() {
    return {
      savingProfileBoard: false,
      savingProfileCover: false,
      showProfileCoverEditor: false,
      user: {
        profile: {
          id: 0,
          pid: "",
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
          profileAvatar: "",
          updateDate: 0,
          createDate: 0
        },
        achievements: [],
        annualRanks: []
      }
    }
  },
  computed: {
    userId() {
      return this.$route.params.id.toString()
    },
    isPersonalProfile() {
      return this.$root.user.profile.pid === this.userId
    }
  },
  methods: {
    openProfileCoverEditor() {
      if(!this.isPersonalProfile) return
      this.showProfileCoverEditor = true
      this.$nextTick(() => {
        this.$refs.imgCutterModal.handleOpen({
          name: new URL(this.user.profile.profileCover).pathname.split("/").pop(),
          src: this.user.profile.profileCover,
        });
      })
    },
    cutDownProfileCover(e) {
      if (!this.isPersonalProfile || this.savingProfileCover) return
      if (e.blob.size > 3000000) {
        this.$root.popupError(new ServerError("ERROR_PROFILE_COVER_TOO_LARGE"))
        return
      }
      this.savingProfileCover = true
      UserAPI.uploadProfileCover(e.blob).then(res => {
        this.savingProfileCover = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          window.location.reload()
        }
      })
    },
    saveBoard() {
      if (this.savingProfileBoard) return
      if (this.user.profile.profileBoard.length < 10) {
        this.$root.popupError(new ServerError("ERROR_PROFILE_BOARD_TOO_SHORT"))
        return
      }
      if (this.user.profile.profileBoard.length > 10000) {
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
    }
  },
  mounted() {
    UserAPI.getUser(this.userId, {
      profile: true,
      achievements: true,
      "annual-ranks": true
    }).then((res) => {
      if (res instanceof ServerError) {
        this.$root.popupError(res)
        return
      }
      if (res.profile.profileCover === "") {
        res.profile.profileCover = profileCoverDefaultImg
      } else {
        res.profile.profileCover = conf.assetURL + "/" + res.profile.profileCover
      }
      if (res.profile.profileAvatar === "") {
        if(res.profile.hasOwnProperty("gender") && res.profile["gender"] === "female") {
          res.profile.profileAvatar = profileFemaleAvatarDefaultImg
        } else {
          res.profile.profileAvatar = profileMaleAvatarDefaultImg
        }
      } else {
        res.profile.profileAvatar = conf.assetURL + "/" + res.profile.profileAvatar
      }
      Object.assign(this.user, res)
      this.$forceUpdate()
      this.$refs.profileLoadingState.deactivate()
    })
  }
}
</script>
