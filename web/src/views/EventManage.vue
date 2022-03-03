<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb text="Quản lý sự kiện" link="/em"></Breadcrumb>
    <div class="mt-10">
      <button class="btn-success" @click="edit(undefined)">Tạo sự kiện</button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="w-max md:w-full">
        <tbody>
          <tr v-for="event in events">
            <td class="max-w-xs break-words">{{ event.title }}</td>
            <td class="text-gray-500">
              {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.startDate)) }} -
              {{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.endDate)) }}
            </td>
            <td class="float-right ml-5 flex flex-row gap-1">
              <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="edit(event.id)"></PencilIcon>
              <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="remove(event.id, event.title)"></TrashIcon>
              <PuzzleIcon class="w-6 cursor-pointer text-gray-500" @click="manageContest(event.id)"></PuzzleIcon>
              <p class="text-gray-500">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.date)) }}</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="mt-10">
      <LoadingState ref="loadingState">
        <div v-if="eventAvailable">
          <button class="btn-info m-auto block" @click="loadNextEvents()">Xem thêm...</button>
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
import lookupErrorCode from "../api/errorCode";

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
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      });
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
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    this.loadNextEvents()
  }
}
</script>
