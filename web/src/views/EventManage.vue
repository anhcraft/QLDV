<template>
  <Header></Header>
  <section class="page-section px-10 py-8 lg:py-16">
    <div>
      <button class="btn-success" @click="edit(undefined)">Tạo sự kiện</button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="w-max md:w-full">
        <thead class="text-left">
          <tr>
            <th>Tên sự kiện</th>
            <th>Trạng thái</th>
            <th>Thời gian</th>
            <th>Ngày cập nhật</th>
            <th>Ngày tạo</th>
            <th>Thao tác</th>
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
            <td class="max-w-xs break-words">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.updateDate)) }}</td>
            <td class="max-w-xs break-words">{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(event.createDate)) }}</td>
            <td class="flex flex-row gap-1">
              <PencilIcon class="w-6 cursor-pointer text-gray-500" @click="edit(event.id)"></PencilIcon>
              <TrashIcon class="w-6 cursor-pointer text-gray-500" @click="remove(event.id, event.title)"></TrashIcon>
              <PuzzlePieceIcon class="w-6 cursor-pointer text-gray-500" @click="manageContest(event.id)"></PuzzlePieceIcon>
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
  </section>
  <Footer></Footer>
  <Prompt @callback="removeEventCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa sự kiện này?</p><br> {{ eventRemoveTitle }}
  </Prompt>
</template>

<script>
import {PencilIcon, PuzzlePieceIcon, TrashIcon} from '@heroicons/vue/24/solid';
import Prompt from "../components/Prompt.vue";
import Header from "../components/Header.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";
import {ServerError} from "../api/server-error";
import EventAPI from "../api/event-api";
import Footer from "../components/Footer.vue";

export default {
  name: "EventManage",
  components: {
    LoadingState, Header, Footer, Breadcrumb,
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
      eventRemoveTitle: '',
      deletingEvent: false
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
      const limit = 15
      EventAPI.listEvents({
        limit: limit,
        "below-id": this.pagination.belowId,
        "begin-date": 0,
        "end-date": 0
      }).then((res) => {
        if(res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        if(res.length < limit) {
          this.pagination.available = false
        }
        if(res.length > 0) {
          this.pagination.belowId = res[res.length - 1].id
          this.events = this.events.concat(res)
        }
        this.$refs.loadingState.deactivate()
      })
    },
    edit(id) {
      if(id === undefined) this.$router.push({name: "createEvent"})
      else this.$router.push({name: "updateEvent", params: { id: id }})
    },
    manageContest(id) {
      //
    },
    remove(id, name) {
      this.eventRemoveId = id
      this.eventRemoveTitle = name
      this.$refs.removePrompt.toggle()
    },
    removeEventCallback(b) {
      if(!b || this.deletingEvent) return
      this.deletingEvent = true
      EventAPI.deleteEvent(this.eventRemoveId).then(s => {
        this.deletingEvent = false
        if(s instanceof ServerError) {
          this.$root.popupError(s)
          return
        }
        this.events = this.events.filter(p => p.id !== this.eventRemoveId)
        this.eventRemoveId = ""
        this.eventRemoveTitle = ""
      })
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    const f = () => {
      if(!this.$root.isLoggedIn() || !this.$root.isGlobalManager) {
        this.$router.push({name: "home"})
        return
      }
      this.loadNextEvents()
      window.addEventListener('scroll', this.handleScroll)
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
