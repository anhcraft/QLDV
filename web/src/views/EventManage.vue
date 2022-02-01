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
          <td class="flex flex-row gap-3">
            <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="edit(event.id)"></PencilIcon>
            <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="remove(event.id, event.title)"></TrashIcon>
            <p class="text-gray-500">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.date)) }}</p>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="loadingEvents">
      <svg class="animate-spin h-6 w-6 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>
    <div class="mt-10" v-else-if="eventAvailable">
      <button class="rounded-md bg-blue-500 hover:bg-blue-600 cursor-pointer px-3 py-2 text-white text-center text-xs m-auto block" @click="loadNextEvents">Xem thêm...</button>
    </div>
    <div class="mt-10" v-else>Đã tải hết sự kiện.</div>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt :content="'<p class=font-bold>Bạn có muốn xóa sự kiện này?</p><br>' + eventRemoveTitle" @callback="removeEventCallback" ref="removePrompt"></Prompt>
</template>

<script>
import {PencilIcon, TrashIcon} from '@heroicons/vue/solid'
import server from "../api/server";
import Prompt from "../components/Prompt.vue";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
  name: "EventManage",
  components: {
    Header, FloatingMenu, Breadcrumb,
    PencilIcon, TrashIcon, Prompt
  },
  data() {
    return {
      loadingEvents: false,
      eventAvailable: true,
      events: [],
      eventRemoveId: '',
      eventRemoveTitle: ''
    }
  },
  methods: {
    loadNextEvents(){
      this.loadingEvents = true
      const older = this.events.length === 0 ? new Date().getTime() : this.events[this.events.length - 1].date
      server.loadEvents(20, older, 0, 0, auth.getToken()).then(s => {
        if(s.events.length === 0) {
          this.eventAvailable = false
        }
        this.events = this.events.concat(s.events)
        this.loadingEvents = false
      })
    },
    edit(id) {
      this.$router.push(`/ee/` + (id === undefined ? '' : id))
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
    if(!this.$root.isLoggedIn) {
      this.$router.push(`/`)
    }
    this.loadNextEvents()
  }
}
</script>
