<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb :text="($route.params.id === undefined ? 'Tạo' : 'Sửa') + ' bài viết'" link="/pm" class="mb-10"></Breadcrumb>
    <div v-if="postLoaded && !submittingPost">
      <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tiêu đề..." v-model="post.title">
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
                  toolbar: 'undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | ' +
                    'bullist numlist outdent indent | link media | ' +
                    'forecolor backcolor emoticons | help',
                  menubar: false
                }"
            v-model="post.content"
        ></Editor>
      </div>
      <div class="mt-10">
        <p>Ảnh đính kèm:</p>
        <div class="my-10">
          <div class="flex flex-row flex-wrap gap-3">
            <img v-for="att in post.attachments" class="max-w-md" :class="{'border-2 border-slate-500 opacity-50' : removeAttachments.includes(att.id)}" :src="serverBaseURL + '/static/' + att.id" alt="" @click="removeAttachment(att.id)" />
          </div>
          <p class="text-red-500 mt-3" v-if="removeAttachments.length > 0">Sẽ xóa {{ removeAttachments.length }} ảnh được chọn.</p>
        </div>
        <p>Tải ảnh mới:</p>
        <input @change="onAttachmentChange" accept="image/*" multiple class="block px-3 py-1.5 text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none" type="file">
        <div class="flex flex-row flex-wrap gap-3 my-10">
          <img v-for="url in attachmentUploadPreviews" class="max-w-md" :src="url" alt=""/>
        </div>
      </div>
      <button class="bg-pink-400 hover:bg-pink-500 cursor-pointer px-4 py-2 text-white text-center text-sm" v-if="!submittingPost" @click="submitPost">{{ $route.params.id === undefined ? "Đăng bài" : "Lưu chỉnh sửa" }}</button>
    </div>
    <div v-else>
      <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import server from "../api/server";
import auth from "../api/auth";
import conf from "../conf";
import Editor from '@tinymce/tinymce-vue'
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
  "name": "PostEdit",
  components: {Header, FloatingMenu, Breadcrumb, Editor },
  data() {
    return {
      post: {
        title: "",
        content: "",
        attachments: []
      },
      postLoaded: false,
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
    jumpToTop() {
      window.scrollTo(0, 0);
    },
    jumpToBottom() {
      window.scrollTo(0, document.body.scrollHeight);
    },
    backToManage() {
      this.$router.push('/pm/')
    },
    submitPost() {
      this.submittingPost = true
      server.changePost(this.$route.params.id, this.post.title, this.post.content, this.removeAttachments, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.attachmentUploadQueue = this.attachmentUpload.length
          if(this.attachmentUploadQueue === 0) {
            this.submittingPost = false
            this.$router.push('/pm/')
          } else {
            for (let i = 0; i < this.attachmentUpload.length; i++) {
              server.uploadPostAttachment(s["id"], this.attachmentUpload[i], auth.getToken()).then(s => {
                if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
                  this.attachmentUploadQueue--
                  if (this.attachmentUploadQueue === 0) {
                    this.submittingPost = false
                    this.$router.push('/pm/')
                  }
                } else {
                  alert(`Lỗi lưu bài viết: ${s["error"]}`)
                }
              })
            }
          }
        } else {
          alert(`Lỗi lưu bài viết: ${s["error"]}`)
        }
      })
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
    if(!this.$root.isLoggedIn) {
      this.$router.push(`/`)
    }
    if(this.$route.params.id !== undefined) {
      server.loadPost(this.$route.params.id).then(s => {
        this.post = s;
        this.postLoaded = true;
      });
    } else {
      this.postLoaded = true;
    }
  }
}
</script>
