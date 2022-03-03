<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Quản lý sự kiện" link="/em"></Breadcrumb>
    <div class="mt-10">
      <button class="bg-white hover:bg-pink-300 cursor-pointer border-2 border-pink-300 px-3 py-1 text-center text-sm" @click="edit(undefined)">Tạo sự kiện</button>
    </div>
    <table class="w-full mt-10">
      <tbody>
        <tr v-for="event in events">
          <td>{{ event.title }}</td>
          <td class="text-gray-500">
            {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.startDate)) }} -
            {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.endDate)) }}
          </td>
          <td class="float-right flex flex-row gap-3">
            <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="edit(event.id)"></PencilIcon>
            <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="remove(event.id, event.title)"></TrashIcon>
            <PuzzleIcon class="w-6 cursor-pointer text-gray-500" @click="manageContest(event.id)"></PuzzleIcon>
            <p class="text-gray-500">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.date)) }}</p>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="mt-10">
      <LoadingState ref="loadingState">
        <div v-if="eventAvailable">
          <button class="rounded-md bg-blue-500 hover:bg-blue-600 cursor-pointer px-3 py-2 text-white text-center text-xs m-auto block" @click="loadNextEvents()">Xem thêm...</button>
        </div>
        <div v-else>Đã tải hết sự kiện.</div>
      </LoadingState>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt @callback="removeEventCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa sự kiện này?</p><br> {{ eventRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, PuzzleIcon, TrashIcon} from '@heroicons/vue/solid'
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";

export default {
  name: "EventManage",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    PencilIcon, TrashIcon, PuzzleIcon, Prompt
  },
  data() {
    return {
      eventAvailable: true,
      events: [],
      eventRemoveId: '',
      eventRemoveTitle: ''
    }
  },
  methods: {
    loadNextEvents(){
      this.$refs.loadingState.activate()
      const older = this.events.length === 0 ? new Date().getTime() : this.events[this.events.length - 1].date
      server.loadEvents(20, older, 0, 0, auth.getToken()).then(s => {
        if(s.events.length === 0) {
          this.eventAvailable = false
        }
        this.events = this.events.concat(s.events)
        this.$refs.loadingState.deactivate()
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
            alert(`Lỗi xóa sự kiện: ${s["error"]}`)
          }
        })
      }
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    this.loadNextEvents()
  }
}
</script>
