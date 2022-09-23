<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <LoadingState ref="loadingState">

      <div class="flex flex-col gap-7">
        <input type="text" class="border-b-2 border-b-slate-300 w-full text-3xl" placeholder="Tiêu đề..."
               v-model.trim="post.title">

        <div class="centered-horizontal">
          <span>#</span>
          <input type="text" class="border-b-2 border-b-slate-300 w-full" placeholder="Hashtag"
                 v-model.trim="post.hashtag" list="hashtags">
          <datalist id="hashtags">
            <option v-for="v in hashtags" :value="v"/>
          </datalist>
        </div>


        <div class="border border-gray-400 py-2 px-5">
          <p class="text-xl">Đề mục:</p>
          <textarea class="border-b-2 border-b-slate-300 w-full" v-model.trim="post.headline"></textarea>
          <p class="italic text-sm">Để trống để tạo đề mục tự động từ nội dung bài viết</p>
        </div>

        <Editor
            apiKey="r7g4lphizuprqmrjv0ooj15pn5qpcesynrg101ekc40avzlg"
            :init="{
                  height: 500,
                  plugins: ['advlist', 'autolink', 'lists', 'link', 'image', 'insertdatetime', 'media', 'table', 'wordcount', 'emoticons', 'table'],
                  toolbar: 'undo redo | styleselect | bold italic | forecolor backcolor emoticons link | ' +
                    'bullist numlist outdent indent | media | ' +
                    'alignleft aligncenter alignright alignjustify | help',
                  menubar: false,
                  branding: false
                }"
            v-model="post.content"
        ></Editor>

        <div class="border border-gray-400 py-2 px-5">
          <p class="text-xl">Giới hạn người xem</p>
          <select class="mt-5 text-sm" v-model.number="post.privacy">
            <option v-for="v in roleTables" :value="v.role">Từ {{ v.name }} trở lên</option>
          </select>
          <p class="mt-5 text-sm italic">
            {{ roleTables.filter(v => v.role >= post.privacy).map(v => v.name).join(", ") }} có thể xem
          </p>
        </div>

        <div class="border border-gray-400 py-2 px-5">
          <p class="text-xl">Ảnh đính kèm:</p>
          <div class="my-10">
            <div class="flex flex-row flex-wrap gap-3">
              <img v-for="att in post.attachments" class="max-h-36"
                   :class="{'border-2 border-slate-500 opacity-50' : removeAttachments.includes(att.id)}"
                   :src="assetURL + '/' + att.id" alt="" @click="toggleAttachment(att.id)"/>
            </div>
            <p class="text-red-500 mt-3" v-if="removeAttachments.length > 0">Sẽ xóa {{ removeAttachments.length }} ảnh
              được chọn.</p>
          </div>

          <p class="font-bold">Tải ảnh mới:</p>
          <input @change="onAttachmentChange" accept="image/*" multiple
                 class="block px-3 py-1.5 text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
                 type="file">
          <div class="flex flex-row flex-wrap gap-3 my-10">
            <img v-for="data in attachmentUpload" class="max-h-36" :src="getImageData(data)" alt=""/>
          </div>
        </div>

        <div>
          <button class="btn-success float-right" :class="{'opacity-50' : submittingPost || submittingAttachments}"
                  @click="submitPost">{{ $route.params.id === undefined ? "Đăng bài" : "Lưu chỉnh sửa" }}
          </button>

          <div
              v-if="submittingAttachments || this.submittedAttachmentSuccessCount < this.submittedAttachmentExpectedCount">
            <progress id="file" :value="submittedAttachmentCount" :max="submittedAttachmentExpectedCount"></progress>
            <p>Đã tải lên {{ submittedAttachmentSuccessCount }} / {{ submittedAttachmentExpectedCount }} ảnh thành
              công.</p>
          </div>
        </div>
      </div>

    </LoadingState>
  </section>
  <Footer></Footer>
</template>

<script>
import Editor from '@tinymce/tinymce-vue'
import Header from "../components/Header.vue";
import LoadingState from "../components/LoadingState.vue";
import PostAPI from "../api/post-api";
import {ServerError} from "../api/server-error";
import {GetRoleTable} from "../auth/roles";
import conf from "../conf";
import Footer from "../components/Footer.vue";

export default {
  "name": "PostEdit",
  components: {LoadingState, Header, Footer, Editor},
  data() {
    return {
      post: {
        title: "",
        content: "",
        headline: "",
        hashtag: "",
        privacy: 0,
        attachments: [] // existing attachments
      },
      hashtags: [],
      submittingPost: false,
      submittingAttachments: false,
      removeAttachments: [], // selected attachments to remove
      attachmentUpload: [], // new attachments waiting to be uploaded
      submittedAttachmentCount: 0,
      submittedAttachmentSuccessCount: 0,
      submittedAttachmentExpectedCount: 0,
      attachmentFailedUpload: []
    }
  },
  computed: {
    assetURL() {
      return conf.assetURL
    },
    roleTables() {
      return GetRoleTable().filter(v => v.role <= this.$root.user.profile.role)
    }
  },
  methods: {
    getImageData(data) {
      return URL.createObjectURL(data)
    },
    submitPost() {
      if (this.submittingPost) return
      if (this.post.title.length > 300) {
        this.$root.popupError(new ServerError("ERROR_POST_TITLE_TOO_LONG"))
        return
      }
      if (this.post.title.length < 10) {
        this.$root.popupError(new ServerError("ERROR_POST_TITLE_TOO_SHORT"))
        return
      }
      if (this.post.headline.length > 250) {
        this.$root.popupError(new ServerError("ERROR_POST_HEADLINE_TOO_LONG"))
        return
      }
      if (this.post.headline.length > 0 && this.post.headline.length < 30) {
        this.$root.popupError(new ServerError("ERROR_POST_HEADLINE_TOO_SHORT"))
        return
      }
      if (this.post.hashtag.length > 20) {
        this.$root.popupError(new ServerError("ERROR_POST_HASHTAG_TOO_LONG"))
        return
      }
      if (this.post.hashtag.length < 5) {
        this.$root.popupError(new ServerError("ERROR_POST_HASHTAG_TOO_SHORT"))
        return
      }
      if (this.post.content.length > 100000) {
        this.$root.popupError(new ServerError("ERROR_POST_CONTENT_TOO_LONG"))
        return
      }
      if (this.post.content.length < 100) {
        this.$root.popupError(new ServerError("ERROR_POST_CONTENT_TOO_SHORT"))
        return
      }
      this.submittingPost = true
      const id = this.$route.params.id === undefined ? "" : this.$route.params.id
      PostAPI.updatePost(id, {
        title: this.post.title,
        privacy: this.post.privacy,
        hashtag: this.post.hashtag,
        headline: this.post.headline,
        content: this.post.content
      }).then(res => {
        if (res instanceof ServerError) {
          this.submittingPost = false
          this.$root.popupError(res)
          return
        }
        return PostAPI.deleteAttachment(this.removeAttachments)
      }).then(res => {
        this.submittingPost = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        if (this.attachmentUpload.length === 0) {
          this.$router.push({name: "managePosts"})
          return;
        }
        this.submittingAttachments = true
        this.submittedAttachmentCount = 0
        this.submittedAttachmentSuccessCount = 0
        this.submittedAttachmentExpectedCount = this.attachmentUpload.length
        this.attachmentFailedUpload = []
        for (let i = 0; i < this.attachmentUpload.length; i++) {
          let a = this.attachmentUpload[i]
          PostAPI.uploadAttachment(id, a).then(res => {
            if (res instanceof ServerError) {
              this.$root.popupError(res)
              this.attachmentFailedUpload.push(a)
            } else {
              this.submittedAttachmentSuccessCount += 1
              this.post.attachments.push(res.id)
            }
            if (++this.submittedAttachmentCount === this.submittedAttachmentExpectedCount) {
              this.submittingAttachments = false
              if (this.submittedAttachmentSuccessCount === this.submittedAttachmentExpectedCount) {
                this.$router.push({name: "managePosts"})
              } else {
                this.attachmentUpload = this.attachmentFailedUpload
                this.attachmentFailedUpload = []
              }
            }
          })
        }
      })
    },
    onAttachmentChange(e) {
      this.attachmentUpload = e.target.files
    },
    toggleAttachment(id) {
      if (this.removeAttachments.includes(id)) {
        this.removeAttachments = this.removeAttachments.filter(a => a !== id)
      } else {
        this.removeAttachments = this.removeAttachments.concat(id)
      }
    }
  },
  mounted() {
    const f = () => {
      if (!this.$root.isLoggedIn() || !this.$root.isGlobalManager) {
        this.$router.push({name: "managePosts"})
        return
      }
      PostAPI.getHashtags().then(data => {
        if (data instanceof ServerError) {
          this.$root.popupError(data)
          return
        }
        this.hashtags = data
      })
      if (this.$route.params.id !== undefined) {
        PostAPI.getPost(this.$route.params.id).then(res => {
          if (res instanceof ServerError) {
            this.$root.popupError(res)
            return
          }
          this.post = res;
          this.$refs.loadingState.deactivate()
        });
      } else {
        this.$refs.loadingState.deactivate()
      }
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
