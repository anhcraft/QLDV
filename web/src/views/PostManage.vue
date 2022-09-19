<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb text="Quản lý bài viết" link="/pm"></Breadcrumb>
    <div class="mt-10">
      <button class="btn-success" @click="createPost">Tạo bài viết</button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="w-max md:w-full">
        <thead class="text-left">
          <tr>
            <th>Tên bài</th>
            <th>Hashtag</th>
            <th>Ngày đăng</th>
            <th>Lượt xem</th>
            <th>Lượt thích</th>
            <th>Thao tác</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="post in posts">
            <td class="max-w-xs break-words">{{ post.title }}</td>
            <td class="max-w-xs break-words">#{{ post.hashtag }}</td>
            <td class="max-w-xs break-words">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(post.date)) }}</td>
            <td class="max-w-xs break-words">{{ post.views }}</td>
            <td class="max-w-xs break-words">{{ post.likes }}</td>
            <td class="ml-5 flex flex-row gap-5">
              <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="editPost(post.id)"></PencilIcon>
              <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="removePost(post.id, post.title)"></TrashIcon>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="mt-10">
      <LoadingState ref="loadingState">
        <div v-if="!pagination.available">Đã tải hết bài viết.</div>
      </LoadingState>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt @callback="removePostCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa bài viết này?</p><br> {{ postRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, TrashIcon} from '@heroicons/vue/24/solid'
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../auth/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "PostManage",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    PencilIcon, TrashIcon, Prompt
  },
  data() {
    return {
      posts: [],
      pagination: {
        belowId: 0,
        lowerThan: 0,
        sortBy: "",
        available: true
      },
      postRemoveId: -1,
      postRemoveTitle: ''
    }
  },
  methods: {
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.loadingState.loading && this.pagination.available) {
          this.loadNextPosts()
        }
      }
    },
    loadNextPosts(){
      this.$refs.loadingState.activate()
      server.loadPosts(10, [], this.pagination.sortBy, this.pagination.lowerThan, this.pagination.belowId, auth.getToken()).then(s => {
        if(s.posts.length < 10) {
          this.pagination.available = false
        }
        if(s.posts.length > 0) {
          this.pagination.belowId = s.posts[s.posts.length - 1].id
          this.posts = this.posts.concat(s.posts)
        }
        this.$refs.loadingState.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải bài viết thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    createPost() {
      this.$router.push(`/pe/`)
    },
    editPost(id) {
      this.$router.push(`/pe/${id}/`)
    },
    removePost(id, name) {
      this.postRemoveId = id
      this.postRemoveTitle = name
      this.$refs.removePrompt.toggle()
    },
    removePostCallback(b) {
      if(b) {
        server.removePost(this.postRemoveId, auth.getToken()).then(s => {
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.posts = this.posts.filter(p => p.id !== this.postRemoveId)
            this.postRemoveId = -1
            this.postRemoveTitle = ''
          } else {
            this.$notify({
              title: "Xóa bài viết thất bại",
              text: lookupErrorCode(s["error"]),
              type: "error"
            });
          }
        }, (e) => {
          this.$notify({
            title: "Xóa bài viết thất bại",
            text: e.message,
            type: "error"
          });
        });
      }
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    this.loadNextPosts()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
