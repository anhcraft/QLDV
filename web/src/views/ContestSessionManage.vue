<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Quản lý cuộc thi" :link="'/mc/' + $route.params.id" class="mb-10"></Breadcrumb>
    <header class="border-b-2 border-b-slate-300 pb-3 text-xl flex flex-row gap-2">
      <div class="grow">{{ event.title }}</div>
    </header>
    <div class="mt-5">

      <table class="w-full">
        <thead>
          <tr>
            <th>Email</th>
            <th>Thời gian bắt đầu</th>
            <th>Thời gian còn lại</th>
            <th>Thời gian lưu bài gần nhất</th>
            <th>Tình trạng</th>
            <th>Điểm</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cs in contestSessions" class="text-sm hover:bg-blue-200 text-center">
            <td class="text-left">{{ cs.userId }}</td>
            <td>{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(cs.startTime)) }}</td>
            <td>{{ stringifyTime(cs.finished ? 0 : Math.max(0, cs.endTime - new Date())) }}</td>
            <td>{{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(cs.lastAnswerSubmittedTime)) }}</td>
            <td>{{ cs.finished ? "Đã hoàn thành" : "Đang làm bài" }}</td>
            <td>{{ cs.hasOwnProperty("expectedAnswerSheet") ? `${(getCorrectAnswers(cs)/cs.answerSheet.length*10).toFixed(1)}  (${getCorrectAnswers(cs)}/${cs.answerSheet.length})` : "" }}</td>
          </tr>
        </tbody>
      </table>
      <LoadingState ref="sessionLoadingState">
        <div class="mt-5" v-if="!sessionAvailable">Đã tải hết.</div>
      </LoadingState>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import server from "../api/server";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import Datepicker from 'vue3-date-time-picker';
import 'vue3-date-time-picker/dist/main.css'
import LoadingState from "../components/LoadingState.vue";
import * as XLSX from "xlsx";
import Prompt from "../components/Prompt.vue";
import Editor from '@tinymce/tinymce-vue'

export default {
  "name": "ContestSessionManage",
  components: {LoadingState, Header, FloatingMenu, Breadcrumb},
  data() {
    return {
      event: {},
      sessionAvailable: true,
      contestSessions: []
    }
  },
  methods: {
    getCorrectAnswers(a) {
      let q = 0;
      for (let i = 0; i < a.answerSheet.length; i++){
        if(a.expectedAnswerSheet[i] === a.answerSheet[i]){
          q++;
        }
      }
      return q
    },
    stringifyTime(num) {
      num /= 1000
      let hours = Math.floor(num / 3600);
      let minutes = Math.floor((num - (hours * 3600)) / 60);
      let seconds = Math.floor(num - (hours * 3600) - (minutes * 60));
      if (hours < 10) {
        hours = "0" + hours;
      }
      if (minutes < 10) {
        minutes = "0" + minutes;
      }
      if (seconds < 10) {
        seconds = "0" + seconds;
      }
      return hours + ':' + minutes + ':' + seconds;
    },
    loadNextSessions(){
      this.$refs.sessionLoadingState.activate()
      const older = this.contestSessions.length === 0 ? new Date().getTime() : this.contestSessions[this.contestSessions.length - 1].lastAnswerSubmittedTime
      server.loadContestSessions(this.$route.params.id, 50, older, auth.getToken()).then(s => {
        if(s.contestSessions.length === 0) {
          this.sessionAvailable = false
        }
        this.contestSessions = this.contestSessions.concat(s.contestSessions.map((v) => {
          v.questionSheet = JSON.parse(v.questionSheet)
          v.answerSheet = JSON.parse(v.answerSheet)
          if(v.hasOwnProperty("expectedAnswerSheet")) {
            v.expectedAnswerSheet = JSON.parse(v.expectedAnswerSheet)
          }
          return v
        }))
        this.$refs.sessionLoadingState.deactivate()
      })
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.sessionLoadingState.loading && this.sessionAvailable) {
          this.loadNextSessions()
        }
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
    if(this.$route.params.id !== undefined) {
      server.loadEvent(this.$route.params.id, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error")) {
          if (s.hasOwnProperty("contest")) {
            s.contest.dataSheet = JSON.parse(s.contest.dataSheet)
            this.event = s;
            this.loadNextSessions()
            window.addEventListener('scroll', this.handleScroll)
          } else {
            alert('Cuộc thi chưa được tạo!')
          }
        } else {
          alert(`Lỗi tải sự kiện: ${s["error"]}`)
        }
      });
    } else {
      alert(`Lỗi tải sự kiện: ${s["error"]}`)
    }
  }
}
</script>
