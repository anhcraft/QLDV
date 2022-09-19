<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb :text="($route.params.id === undefined ? 'Tạo' : 'Sửa') + ' bài viết'" link="/pm" class="mb-10"></Breadcrumb>
    <LoadingState ref="loadingState">
      <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tiêu đề..." v-model="post.title">
      <div class="mt-10 centered-horizontal">
        <span>#</span>
        <input type="text" class="border-b-2 border-b-slate-300 w-full" placeholder="Hashtag" v-model="post.hashtag" list="hashtags">
        <datalist id="hashtags">
          <option v-for="v in hashtags" :value="v"/>
        </datalist>
      </div>
      <div class="mt-10">
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
            v-model="post.content"
        ></Editor>
        <div class="border border-gray-300 py-2 px-5 my-10">
          <div class="flex flex-row gap-5 place-items-center">
            <p>Chỉ cho thành viên xem</p>
            <input type="checkbox" class="w-4 h-4" v-bind:checked="(post.privacy & 1) === 1" @input="post.privacy = $event.target.value ? (post.privacy ^ 1) : (post.privacy | 1)">
          </div>
          <div class="flex flex-row gap-5 place-items-center">
            <p>Chỉ cho bí thư xem</p>
            <input type="checkbox" class="w-4 h-4" v-bind:checked="(post.privacy & 2) === 2" @input="post.privacy = $event.target.value ? (post.privacy ^ 2) : (post.privacy | 2)">
          </div>
          <div class="flex flex-row gap-5 place-items-center">
            <p>Chỉ cho quản trị viên xem</p>
            <input type="checkbox" class="w-4 h-4" v-bind:checked="(post.privacy & 4) === 4" @input="post.privacy = $event.target.value ? (post.privacy ^ 4) : (post.privacy | 4)">
          </div>
        </div>
      </div>
      <div class="mt-10">
        <p>Ảnh đính kèm:</p>
        <div class="my-10">
          <div class="flex flex-row flex-wrap gap-3">
            <img v-for="att in post.attachments" class="max-h-36" :class="{'border-2 border-slate-500 opacity-50' : removeAttachments.includes(att.id)}" :src="serverBaseURL + '/static/' + att.id" alt="" @click="removeAttachment(att.id)" />
          </div>
          <p class="text-red-500 mt-3" v-if="removeAttachments.length > 0">Sẽ xóa {{ removeAttachments.length }} ảnh được chọn.</p>
        </div>
        <p>Tải ảnh mới:</p>
        <input @change="onAttachmentChange" accept="image/*" multiple class="block px-3 py-1.5 text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none" type="file">
        <div class="flex flex-row flex-wrap gap-3 my-10">
          <img v-for="url in attachmentUploadPreviews" class="max-h-36" :src="url" alt=""/>
        </div>
      </div>
      <button class="bg-pink-400 hover:bg-pink-500 cursor-pointer px-4 py-2 text-white text-center text-sm" v-if="!submittingPost" @click="submitPost">{{ $route.params.id === undefined ? "Đăng bài" : "Lưu chỉnh sửa" }}</button>
    </LoadingState>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import server from "../api/server";
import auth from "../auth/auth";
import conf from "../conf";
import Editor from '@tinymce/tinymce-vue'
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

export default {
  "name": "PostEdit",
  components: {LoadingState, Header, FloatingMenu, Breadcrumb, Editor },
  data() {
    return {
      post: {
        title: "",
        content: "",
        hashtag: "",
        attachments: [],
        privacy: 0
      },
      hashtags: [],
      submittingPost: false,
      attachmentUpload: [],
      attachmentUploadQueue: 0,
      attachmentUploadPreviews: [],
      removeAttachments: []
    }
  },
  computed: {
    serverBaseURL() {
      return conf.server
    }
  },
  methods: {
    submitPost() {
      this.submittingPost = true
      const id = this.$route.params.id === undefined ? 0 : this.$route.params.id
      server.changePost(id, this.post.title, this.post.content, this.post.privacy, this.post.hashtag, this.removeAttachments, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.attachmentUploadQueue = this.attachmentUpload.length
          if(this.attachmentUploadQueue === 0) {
            this.submittingPost = false
            this.$router.push('/pm/')
          } else {
            for (let i = 0; i < this.attachmentUpload.length; i++) {
              server.uploadPostAttachment(s["id"], this.attachmentUpload[i], auth.getToken()).then(ss => {
                if (!ss.hasOwnProperty("error") && ss.hasOwnProperty("success") && ss["success"]) {
                  this.$notify({
                    title: "Tải ảnh thành công",
                    text: "",
                    type: "success"
                  });
                  this.attachmentUploadQueue--
                  if (this.attachmentUploadQueue === 0) {
                    this.submittingPost = false
                    this.$router.push('/pm/')
                  }
                } else {
                  this.$notify({
                    title: "Tải ảnh thất bại",
                    text: lookupErrorCode(ss["error"]),
                    type: "error"
                  });
                }
              }, (e) => {
                this.$notify({
                  title: "Tải ảnh thất bại",
                  text: e.message,
                  type: "error"
                });
              });
            }
          }
        } else {
          this.submittingPost = false
          this.$notify({
            title: "Lưu thay đổi thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        }
      }, (e) => {
        this.$notify({
          title: "Lưu thay đổi thất bại",
          text: e.message,
          type: "error"
        });
      });
    },
    onAttachmentChange(e) {
      const data = [];
      for(let i = 0; i < e.target.files.length; i++){
        data.push(URL.createObjectURL(e.target.files[i]))
      }
      this.attachmentUploadPreviews = data
      this.attachmentUpload = e.target.files
    },
    removeAttachment(id) {
      if(this.removeAttachments.includes(id)) {
        this.removeAttachments = this.removeAttachments.filter(a => a !== id)
      } else {
        this.removeAttachments = this.removeAttachments.concat(id)
      }
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    server.getHashtags().then(data => {
      this.hashtags = data
    })
    if(this.$route.params.id !== undefined) {
      server.loadPost(this.$route.params.id, auth.getToken()).then(s => {
        if (s.hasOwnProperty("error")) {
          this.$notify({
            title: "Tải bài viết thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
          return
        }
        this.post = s;
        this.$refs.loadingState.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải bài viết thất bại",
          text: e.message,
          type: "error"
        });
      });
    } else {
      this.$refs.loadingState.deactivate()
    }
  }
}
</script>
