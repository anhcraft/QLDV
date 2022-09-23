<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <div>
      <button class="btn-success" @click="createPost">Tạo bài viết</button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="w-max md:w-full">
        <thead class="text-left">
        <tr>
          <th>Tên bài</th>
          <th>Hashtag</th>
          <th>Ảnh đính kèm</th>
          <th>Ngày cập nhật</th>
          <th>Ngày tạo</th>
          <th>Lượt xem</th>
          <th>Lượt thích</th>
          <th>Thao tác</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="post in posts">
          <td class="max-w-xs break-words">{{ post.title }}</td>
          <td class="max-w-xs break-words">#{{ post.hashtag }}</td>
          <td class="max-w-xs break-words">{{ post.attachments.length }}</td>
          <td class="max-w-xs break-words">{{
              new Intl.DateTimeFormat("vi-VN", {
                timeStyle: "medium",
                dateStyle: "short"
              }).format(new Date(post.updateDate))
            }}
          </td>
          <td class="max-w-xs break-words">{{
              new Intl.DateTimeFormat("vi-VN", {
                timeStyle: "medium",
                dateStyle: "short"
              }).format(new Date(post.createDate))
            }}
          </td>
          <td class="max-w-xs break-words">{{ post.stats.views }}</td>
          <td class="max-w-xs break-words">{{ post.stats.likes }}</td>
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
  </section>
  <Footer></Footer>
  <Prompt @callback="removePostCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa bài viết này?</p><br> {{ postRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, TrashIcon} from '@heroicons/vue/24/solid';
import Prompt from "../components/Prompt.vue";
import Header from "../components/Header.vue";
import LoadingState from "../components/LoadingState.vue";
import PostAPI from "../api/post-api";
import {ServerError} from "../api/server-error";
import Footer from "../components/Footer.vue";

export default {
  name: "PostManage",
  components: {
    Footer, LoadingState, Header, PencilIcon, TrashIcon, Prompt
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
      postRemoveTitle: '',
      deletingPost: false
    }
  },
  methods: {
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if (!this.$refs.loadingState.loading && this.pagination.available) {
          this.loadNextPosts()
        }
      }
    },
    loadNextPosts() {
      this.$refs.loadingState.activate()
      const limit = 15
      PostAPI.listPosts({
        limit: limit,
        "below-id": this.pagination.belowId,
        "filter-hashtags": [],
        "sort-by": this.pagination.sortBy,
        "lower-than": this.pagination.lowerThan
      }).then((res) => {
        if (res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        if (res.length < limit) {
          this.pagination.available = false
        }
        if (res.length > 0) {
          this.pagination.belowId = res[res.length - 1].id
          this.posts = this.posts.concat(res)
        }
        this.$refs.loadingState.deactivate()
      })
    },
    createPost() {
      this.$router.push({name: "createPost"})
    },
    editPost(id) {
      this.$router.push({name: "updatePost", params: {id: id}})
    },
    removePost(id, name) {
      this.postRemoveId = id
      this.postRemoveTitle = name
      this.$refs.removePrompt.toggle()
    },
    removePostCallback(b) {
      if (!b || this.deletingPost) return
      this.deletingPost = true
      PostAPI.deletePost(this.postRemoveId).then(s => {
        if (s instanceof ServerError) {
          this.deletingPost = false
          this.$root.popupError(s)
          return
        }
        this.posts = this.posts.filter(p => p.id !== this.postRemoveId)
        this.postRemoveId = -1
        this.postRemoveTitle = ''
      })
    }
  },
  unmounted() {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    const f = () => {
      if (!this.$root.isLoggedIn() || !this.$root.isGlobalManager) {
        this.$router.push({name: "home"})
        return
      }
      this.loadNextPosts()
      window.addEventListener('scroll', this.handleScroll)
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
