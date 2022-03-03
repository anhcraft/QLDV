<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Vào thi" :link="'/c/' + event.id" class="mb-10"></Breadcrumb>
    <LoadingState ref="loadingState">
      <header class="border-b-2 border-b-slate-300 pb-3 text-xl flex flex-row gap-2">
        <div class="grow">{{ event.title }}</div>
        <button class="bg-sky-500 hover:bg-sky-600 cursor-pointer px-4 py-2 text-white text-center text-sm" @click="joinContest" v-if="contestSession === undefined && event.contest.acceptingAnswers">VÀO THI</button>
      </header>
      <div class="grid grid-cols-3">
        <div class="pt-5 pr-5 break-all" :class="contestSession === undefined || contestSession.finished ? 'col-span-2' : 'col-span-3'">
          <div v-if="contestSession === undefined" class="break-words prose max-w-max" v-html="event.contest.info"></div>
          <div v-else ref="questionContainer">
            <div class="mb-10" v-for="(q, i) in contestSession.questionSheet">
              <div class="text-lg">{{ q.question }}</div>
              <div class="grid grid-cols-2">
                <div class="flex flex-row gap-2 place-items-center" v-for="(c, j) in q.choices" :class="getResultBackground(i, j)">
                  <input type="radio" :value="j" v-model="contestSession.answerSheet[i]" :disabled="endingContest || contestSession.finished">
                  <div>{{ c }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="border-l-2 border-l-slate-300 p-5" v-if="contestSession === undefined || contestSession.finished">
          <div class="mb-10 flex flex-row place-items-center" v-if="contestSession !== undefined && contestSession.finished">
            <div class="grow text-center">
              <span class="text-5xl">{{ rightAnswers }}</span><span class="text-4xl"> / {{ contestSession.answerSheet.length }}</span>
            </div>
            <div>
              <p>Điểm: {{ (rightAnswers/contestSession.answerSheet.length*10).toFixed(1) }}</p>
              <p>Chính xác: {{ Math.floor(rightAnswers/contestSession.answerSheet.length*100) }}%</p>
            </div>
          </div>
          <div class="px-5">
            <div class="flex flex-row gap-2">
              <QuestionMarkCircleIcon class="w-4"></QuestionMarkCircleIcon>
              <p>{{ event.contest.limitQuestions }} câu hỏi</p>
            </div>
            <div class="flex flex-row gap-2">
              <ClockIcon class="w-4"></ClockIcon>
              <p>Thời gian: {{ this.stringifyTime(event.contest.limitTime) }}</p>
            </div>
            <div class="flex flex-row gap-2">
              <UsersIcon class="w-4"></UsersIcon>
              <p>Số lượt làm: 0</p>
            </div>
          </div>
        </div>
      </div>
      <div class="fixed top-1/2 right-10 -translate-y-1/2	p-10 pb-5 bg-gray-200 shadow-lg shadow-slate-400" v-if="contestSession !== undefined && !contestSession.finished">
        <div class="text-3xl text-center font-light">
          {{ stringifyTime(timeLeft) }}
        </div>
        <div class="grid grid-cols-5 gap-3 place-items-center mt-10">
          <svg height="16" width="16" class="cursor-pointer" v-for="(_, i) in contestSession.questionSheet" @click="scrollToQuestion(i)">
            <circle cx="8" cy="8" r="8" :fill="contestSession.answerSheet[i] === -1 ? '#aaa' : '#3b73c2'" />
          </svg>
        </div>
        <button class="bg-rose-400 hover:bg-rose-500 cursor-pointer px-4 py-2 text-white text-center text-sm block m-auto mt-10" @click="submitContest" v-if="!endingContest">NỘP BÀI</button>
        <div class="mt-10" v-if="savingContest">
          <LoadingState text="Đang lưu bài"></LoadingState>
        </div>
        <div class="mt-10" v-else-if="endingContest">
          <LoadingState text="Đang nộp bài"></LoadingState>
        </div>
      </div>
   </LoadingState>
  </div>
  <FloatingMenu></FloatingMenu>
  <div v-if="!endingContest">
    <Prompt @callback="submitContestPromptCallback" ref="submitPrompt">
      <p class=font-bold>Bài làm chưa hoàn thành, bạn có chắc chắn muốn nộp?</p>
    </Prompt>
  </div>
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
import Prompt from "../components/Prompt.vue";

export default {
  "name": "Contest",
  components: {LoadingState, Header, FloatingMenu, Breadcrumb, Prompt, QuestionMarkCircleIcon, UsersIcon, ClockIcon},
  data() {
    return {
      event: {},
      contestSession: undefined,
      countdown: undefined,
      autoSaver: undefined,
      timeLeft: 0,
      savingContest: false,
      endingContest: false
    }
  },
  computed: {
    rightAnswers() {
      let q = 0;
      for (let i = 0; i < this.contestSession.answerSheet.length; i++){
        if(this.contestSession.expectedAnswerSheet[i] === this.contestSession.answerSheet[i]){
          q++;
        }
      }
      return q
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
          window.location.reload()
        }
      })
    },
    submitContest() {
      if(this.endingContest) return
      if(this.contestSession.answerSheet.every((v) => v !== -1)){
        this.forceSubmitContest();
      } else {
        this.$refs.submitPrompt.toggle()
      }
    },
    submitContestPromptCallback(b) {
      if(b) this.forceSubmitContest()
    },
    forceSubmitContest() {
      if(this.endingContest) return
      this.endingContest = true
      const f = () => {
        // the saving process should have finished before the submitting part happens
        if(this.savingContest) {
          setTimeout(f, 1000)
          return
        }
        server.submitContestSession(this.$route.params.id, this.contestSession.answerSheet, false, auth.getToken()).then((s) => {
          if(s.hasOwnProperty("error") || (s.hasOwnProperty("success") && !s["success"])) {
            alert(`Lỗi nộp bài: ${s["error"]}`)
          } else {
            setTimeout(() => window.location.reload(), 2000)
          }
        })
      }
      f();
    },
    scrollToQuestion(i) {
      this.$refs.questionContainer.children.item(i).scrollIntoView({ behavior: 'smooth' });
    },
    getResultBackground(i, j) {
      if(this.contestSession.hasOwnProperty("expectedAnswerSheet")){
        if(this.contestSession.answerSheet[i] === -1) {
          return ""
        } else if (this.contestSession.answerSheet[i] === j) {
          if (this.contestSession.expectedAnswerSheet[i] === j) {
            return "bg-emerald-300"
          } else {
            return "bg-rose-200"
          }
        }
      }
      return ""
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
                if(s.hasOwnProperty("expectedAnswerSheet")) {
                  s.expectedAnswerSheet = JSON.parse(s.expectedAnswerSheet)
                }
                this.contestSession = s
                this.$refs.loadingState.deactivate()
                if(s.finished) return
                this.countdown = setInterval(() => {
                  if(this.endingContest) return
                  this.timeLeft = Math.max(0, s.endTime - new Date().getTime())
                  if(this.timeLeft === 0) {
                    this.forceSubmitContest()
                  }
                }, 1000);
                this.autoSaver = setInterval(() => {
                  if(this.endingContest) return
                  this.savingContest = true
                  server.submitContestSession(q.id, this.contestSession.answerSheet, true, auth.getToken()).then((s) => {
                    if(s.hasOwnProperty("error") || (s.hasOwnProperty("success") && !s["success"])) {
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
