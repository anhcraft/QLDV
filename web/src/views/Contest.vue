<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Vào thi" :link="'/c/' + event.id" class="mb-10"></Breadcrumb>
    <LoadingState ref="loadingState">
      <header class="border-b-2 border-b-slate-300 pb-3 text-xl flex flex-row gap-2">
        <div class="grow">{{ event.title }}</div>
        <button class="bg-sky-500 hover:bg-sky-600 cursor-pointer px-4 py-2 text-white text-center text-sm" @click="joinContest" v-if="contestSession === undefined">VÀO THI</button>
      </header>
      <div class="grid grid-cols-3" v-if="contestSession === undefined">
        <div class="col-span-2 pt-5">
          <div v-html="event.contest.info"></div>
        </div>
        <div class="border-l-2 border-l-slate-300 pl-5 pt-5">
          <div class="flex flex-row gap-1">
            <QuestionMarkCircleIcon class="w-4"></QuestionMarkCircleIcon>
            <p>{{ event.contest.limitQuestions }} câu hỏi</p>
          </div>
          <div class="flex flex-row gap-1">
            <ClockIcon class="w-4"></ClockIcon>
            <p>Thời gian: {{ this.stringifyTime(event.contest.limitTime) }}</p>
          </div>
          <div class="flex flex-row gap-1">
            <UsersIcon class="w-4"></UsersIcon>
            <p>Số lượt làm: 0</p>
          </div>
        </div>
      </div>
      <div v-else>
        <div class="pt-5" ref="questionContainer">
          <div class="mb-10" v-for="(q, i) in contestSession.questionSheet">
            <div class="text-lg">{{ q.question }}</div>
            <div class="grid grid-cols-2">
              <div class="flex flex-row gap-2 place-items-center" v-for="(c, j) in q.choices">
                <input type="radio" :value="j" v-model="contestSession.answerSheet[i]">
                <div>{{ c }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="fixed top-1/2 right-10 -translate-y-1/2	p-10 pb-5 bg-gray-200 shadow-lg shadow-slate-400">
          <div class="text-3xl text-center font-light">
            {{ stringifyTime(timeLeft) }}
          </div>
          <div class="grid grid-cols-5 gap-3 place-items-center mt-10">
            <svg height="16" width="16" class="cursor-pointer" v-for="(_, i) in contestSession.questionSheet" @click="scrollToQuestion(i)">
              <circle cx="8" cy="8" r="8" :fill="contestSession.answerSheet[i] === -1 ? '#aaa' : '#3b73c2'" />
            </svg>
          </div>
          <button class="bg-rose-400 hover:bg-rose-500 cursor-pointer px-4 py-2 text-white text-center text-sm block m-auto mt-10" @click="submitContest">NỘP BÀI</button>
          <div class="mt-10" v-if="savingContest">
            <LoadingState text="Đang lưu bài" ref="contestSaveLoadingState"></LoadingState>
          </div>
        </div>
      </div>
   </LoadingState>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import server from "../api/server";
import auth from "../api/auth";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import 'vue3-date-time-picker/dist/main.css'
import LoadingState from "../components/LoadingState.vue";
import {ClockIcon, QuestionMarkCircleIcon, UsersIcon} from "@heroicons/vue/solid";

export default {
  "name": "Contest",
  components: {LoadingState, Header, FloatingMenu, Breadcrumb, QuestionMarkCircleIcon, UsersIcon, ClockIcon},
  data() {
    return {
      event: {},
      contestSession: undefined,
      countdown: undefined,
      autoSaver: undefined,
      timeLeft: 0,
      savingContest: false
    }
  },
  methods: {
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
    joinContest() {
      this.$refs.loadingState.activate()
      server.joinContestSession(this.$route.params.id, auth.getToken()).then(s => {
        if(s.hasOwnProperty("error")) {
          alert(`Lỗi tải cuộc thi: ${s["error"]}`)
        } else {
          this.contestSession = s
          this.$refs.loadingState.deactivate()
        }
      })
    },
    scrollToQuestion(i) {
      this.$refs.questionContainer.children.item(i).scrollIntoView({ behavior: 'smooth' });
    }
  },
  beforeDestroy() {
    if(this.countdown !== undefined) {
      clearInterval(this.countdown);
    }
    if(this.autoSaver !== undefined) {
      clearInterval(this.autoSaver);
    }
  },
  mounted() {
    if(!this.$root.isLoggedIn()) {
      this.$router.push(`/`)
      return
    }
    if(this.$route.params.id !== undefined) {
      server.loadEvent(this.$route.params.id, auth.getToken()).then(q => {
        if(!q.hasOwnProperty("error")) {
          if (q.hasOwnProperty("contest")) {
            q.contest.dataSheet = JSON.parse(q.contest.dataSheet)
            this.event = q;
            server.loadContestSession(q.id, auth.getToken()).then(s => {
              if(s.hasOwnProperty("error")) {
                if(s["error"] === "ERR_UNKNOWN_CONTEST_SESSION") {
                  this.contestSession = undefined
                  this.$refs.loadingState.deactivate()
                } else {
                  alert(`Lỗi tải cuộc thi: ${s["error"]}`)
                }
              } else {
                s.questionSheet = JSON.parse(s.questionSheet)
                s.answerSheet = JSON.parse(s.answerSheet)
                this.contestSession = s
                this.$refs.loadingState.deactivate()
                this.countdown = setInterval(() => {
                  this.timeLeft = Math.max(0, s.endTime - new Date().getTime())
                  if(this.timeLeft === 0) {

                  }
                }, 1000);
                this.autoSaver = setInterval(() => {
                  this.savingContest = true
                  server.submitContestSession(q.id, this.contestSession.answerSheet, true, auth.getToken()).then((s) => {
                    if(s.hasOwnProperty("error")) {
                     // alert(`Lỗi lưu dữ liệu: ${s["error"]}`)
                    } else {
                      this.savingContest = false
                    }
                  })
                }, 10000);
              }
            })
          } else {
            alert('Lỗi sự kiện invalid')
          }
        } else {
          alert(`Lỗi tải sự kiện: ${q["error"]}`)
        }
      });
    } else {
      alert(`Lỗi tải sự kiện invalid`)
    }
  }
}
</script>
