<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb text="Quản lý sự kiện" link="/em"></Breadcrumb>
    <div class="mt-10">
      <button class="btn-success" @click="edit(undefined)">Tạo sự kiện</button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="w-max md:w-full">
        <thead class="text-left">
          <tr>
            <th>Tên sự kiện</th>
            <th>Trạng thái</th>
            <th>Thời gian</th>
            <th>Thao tác</th>
            <th>Ngày đăng</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in events">
            <td class="max-w-xs break-words">{{ event.title }}</td>
            <td>{{ event.status === 'ongoing' ? "Đang diễn ra" : (event.status === 'finished' ? "Đã kết thúc" : "Chưa bắt đầu") }}</td>
            <td>
              {{
                new Intl.DateTimeFormat("vi-VN", {
                  timeStyle: "medium",
                  dateStyle: "short"
                }).format(new Date(event.beginDate))
              }} -
              {{
                new Intl.DateTimeFormat("vi-VN", {
                  timeStyle: "medium",
                  dateStyle: "short"
                }).format(new Date(event.endDate))
              }}
            </td>
            <td class="flex flex-row gap-1">
              <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="edit(event.id)"></PencilIcon>
              <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="remove(event.id, event.title)"></TrashIcon>
              <PuzzlePieceIcon class="w-6 cursor-pointer text-gray-500" @click="manageContest(event.id)"></PuzzlePieceIcon>
            </td>
            <td class="text-gray-500">
              {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.date)) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="mt-10">
      <LoadingState ref="loadingState">
        <div v-if="!pagination.available">Đã tải hết sự kiện.</div>
      </LoadingState>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt @callback="removeEventCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa sự kiện này?</p><br> {{ eventRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, PuzzlePieceIcon, TrashIcon} from '@heroicons/vue/24/solid';
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

export default {
  name: "EventManage",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    PencilIcon, TrashIcon, PuzzlePieceIcon, Prompt
  },
  data() {
    return {
      events: [],
      pagination: {
        belowId: 0,
        available: true
      },
      eventRemoveId: '',
      eventRemoveTitle: ''
    }
  },
  methods: {
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.loadingState.loading && this.pagination.available) {
          this.loadNextEvents()
        }
      }
    },
    loadNextEvents(){
      this.$refs.loadingState.activate()
      server.loadEvents(10, this.pagination.belowId, 0, 0, auth.getToken()).then(s => {
        if(s.events.length < 10) {
          this.pagination.available = false
        }
        if(s.events.length > 0) {
          this.pagination.belowId = s.events[s.events.length - 1].id
          this.events = this.events.concat(s.events)
        }
        this.$refs.loadingState.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    edit(id) {
      this.$router.push(`/ee/` + (id === undefined ? '' : id))
    },
    manageContest(id) {
      this.$router.push(`/mc/` + (id === undefined ? '' : id))
    },
    remove(id, name) {
      this.eventRemoveId = id
      this.eventRemoveTitle = name
      this.$refs.removePrompt.toggle()
    },
    removeEventCallback(b) {
      if(b) {
        server.removeEvent(this.eventRemoveId, auth.getToken()).then(s => {
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.events = this.events.filter(p => p.id !== this.eventRemoveId)
            this.eventRemoveId = ""
            this.eventRemoveTitle = ""
          } else {
            this.$notify({
              title: "Xóa sự kiện thất bại",
              text: lookupErrorCode(s["error"]),
              type: "error"
            });
          }
        }, (e) => {
          this.$notify({
            title: "Xóa sự kiện thất bại",
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
    this.loadNextEvents()
    window.addEventListener('scroll', this.handleScroll)
  }
}
</script>
