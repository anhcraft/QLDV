<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16 p-5 md:px-10">
    <Breadcrumb text="Vào thi" :link="'/c/' + event.id" class="mb-10"></Breadcrumb>
    <LoadingState ref="loadingState">
      <header class="border-b-2 border-b-slate-300 pb-3 text-xl flex flex-col md:flex-row gap-2">
        <div class="grow break-words">{{ event.title }}</div>
        <button class="btn-info mr-5" :class="{'opacity-50' : !event.contest.acceptingAnswers || this.contestSessions.length >= this.event.contest.limitSessions}" @click="joinContest" v-if="this.activeContestSession === undefined">Vào thi</button>
      </header>

      <div class="grid grid-cols-1 md:grid-cols-3">
        <div class="md:border-l-2 md:border-l-slate-300 pt-5 md:p-5 md:order-last" v-if="activeContestSession === undefined">
          <div class="border-l-4 border-l-gray-400 bg-gray-300 px-4 py-2">Thông tin cuộc thi</div>
          <div class="mt-3">
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
              <p>Số lượt làm: {{ option.series[0].data.reduce((s, a) => s + a, 0) }}</p>
            </div>
            <div class="break-words prose max-w-max mt-5" v-html="event.contest.info"></div>
            <div class="w-full h-48 mt-5">
              <v-chart class="chart" :option="option" />
            </div>
          </div>
        </div>

        <div class="pt-5 pr-5 break-words" :class="activeContestSession === undefined ? 'col-span-2' : 'col-span-3'">
          <div class="mt-10 md:mt-20" v-if="activeContestSession === undefined && currentContestSession === undefined">
            <img src="../assets/behind-leaves.svg" class="m-auto max-w-[20rem] md:max-w-xs" alt=""/>
            <div class="text-2xl text-center font-light mt-5">Hãy vào thi ngay nào!</div>
          </div>

          <div v-else>
            <div class="mb-5" v-if="activeContestSession === undefined && currentContestSession !== undefined">
              <div class="border-l-4 border-l-green-300 bg-green-200 px-4 py-2 flex flex-row">
                <div class="grow">Thông tin bài làm #{{ selectedContestSession + 1 }} ({{ new Intl.DateTimeFormat("vi-VN" , {timeStyle: "medium", dateStyle: "short"}).format(new Date(currentContestSession.startTime)) }})</div>
                <div class="flex flex-row">
                  <div class="p-1 cursor-pointer hover:bg-green-100" :class="{'opacity-20' : selectedContestSession === 0}"
                       @click="selectedContestSession = Math.max(0, selectedContestSession-1)">
                    <ChevronLeftIcon class="w-4 text-gray-600"></ChevronLeftIcon>
                  </div>
                  <div class="p-1 cursor-pointer hover:bg-green-100" :class="{'opacity-20' : selectedContestSession === contestSessions.length-1}"
                       @click="selectedContestSession = Math.min(contestSessions.length-1, selectedContestSession+1)">
                    <ChevronRightIcon class="w-4 text-gray-600"></ChevronRightIcon>
                  </div>
                </div>
              </div>
              <div class="flex flex-row place-items-center mt-3">
                <div class="grow text-center">
                  <span class="text-5xl">{{ currentContestSession.answerAccuracy.filter((v) => v === true).length }}</span><span class="text-4xl"> / {{ currentContestSession.answerSheet.length }}</span>
                </div>
                <div class="text-sm">
                  <p>Điểm: {{ currentContestSession.score.toFixed(1) }}</p>
                  <p>Chính xác: {{ Math.floor(currentContestSession.score*10) }}%</p>
                  <p>Thời gian: {{ this.stringifyTime(currentContestSession.lastAnswerSubmittedTime - currentContestSession.startTime) }}</p>
                </div>
              </div>
            </div>

            <div ref="questionContainer">
              <div class="mb-5 md:mb-10" v-for="(q, i) in currentContestSession.questionSheet">
                <div :class="{'text-lg' : activeContestSession !== undefined}">{{ q.question }}</div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-x-2 gap-y-1 mt-1">
                  <div class="flex flex-row gap-2 place-items-center" v-for="(c, j) in q.choices" :class="getResultBackground(currentContestSession, i, j)">
                    <input type="radio" class="scale-125" :value="j" v-model="currentContestSession.answerSheet[i]" :disabled="endingContest || activeContestSession === undefined">
                    <div :class="{'text-sm' : activeContestSession === undefined}">{{ c }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="fixed top-20 md:top-1/2 right-0 md:right-10 md:-translate-y-1/2	p-5 md:p-8 bg-gray-200 shadow-lg shadow-slate-400" v-if="activeContestSession !== undefined">
        <div class="text-2xl md:text-3xl text-center font-light">
          {{ stringifyTime(timeLeft) }}
        </div>
        <div class="grid grid-cols-4 gap-1 place-items-center mt-10">
          <svg height="16" width="16" class="scale-75 cursor-pointer" v-for="(_, i) in activeContestSession.questionSheet" @click="scrollToQuestion(i)">
            <circle cx="8" cy="8" r="8" :fill="activeContestSession.answerSheet[i] === -1 ? '#aaa' : '#3b73c2'" />
          </svg>
        </div>
        <button class="btn-danger block m-auto mt-10" @click="submitContest" v-if="!endingContest">Nộp bài</button>
        <div class="mt-5" v-if="savingContest">
          <LoadingState text="Đang lưu bài"></LoadingState>
        </div>
        <div class="mt-5" v-else-if="endingContest">
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
import lookupErrorCode from "../api/errorCode";
import utils from "../api/utils"
import {ChevronLeftIcon, ChevronRightIcon} from "@heroicons/vue/outline";
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import {LineChart} from "echarts/charts";
import {GridComponent, TitleComponent} from "echarts/components";
import VChart from "vue-echarts";

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  GridComponent
]);

export default {
  "name": "Contest",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb, Prompt, VChart,
    QuestionMarkCircleIcon, UsersIcon, ClockIcon, ChevronLeftIcon, ChevronRightIcon
  },
  data() {
    return {
      event: {},
      contestSessions: [],
      activeContestSession: undefined,
      selectedContestSession: 0,
      countdown: undefined,
      autoSaver: undefined,
      timeLeft: 0,
      savingContest: false,
      endingContest: false,
      option: {
        title: {
          text: "Thống kê kết quả",
          left: "center"
        },
        xAxis: {
          type: 'category',
          data: []
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            data: [],
            type: 'line'
          }
        ]
      }
    }
  },
  computed: {
    currentContestSession() {
      if (this.activeContestSession === undefined) {
        return this.contestSessions.length === 0 ? undefined : this.contestSessions[this.selectedContestSession];
      } else {
        return this.activeContestSession
      }
    },
    eventId() {
      let s = this.$route.params.id.split(".")
      s = s[s.length - 1]
      s = s.replace(/\D/i, s)
      return parseInt(s)
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
      if(!this.event.contest.acceptingAnswers){
        this.$notify({
          title: "Tải cuộc thi thất bại",
          text: lookupErrorCode("ERR_CONTEST_CLOSED"),
          type: "error"
        });
        return
      }
      if(this.contestSessions.length >= this.event.contest.limitSessions){
        this.$notify({
          title: "Tải cuộc thi thất bại",
          text: lookupErrorCode("ERR_CONTEST_ATTENDED_MAX"),
          type: "error"
        });
        return
      }
      this.$refs.loadingState.activate()
      server.joinContestSession(this.eventId, auth.getToken()).then(s => {
        if(s.hasOwnProperty("error")) {
          this.$notify({
            title: "Tải cuộc thi thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        } else {
          window.location.reload()
        }
      }, (e) => {
        this.$notify({
          title: "Kết nối máy chủ thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    submitContest() {
      if(this.endingContest) return
      if(this.activeContestSession.answerSheet.every((v) => v !== -1)){
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
        server.submitContestSession(this.activeContestSession.id, this.activeContestSession.answerSheet, false, auth.getToken()).then((s) => {
          if(s.hasOwnProperty("error") || (s.hasOwnProperty("success") && !s["success"])) {
            this.$notify({
              title: "Nộp bài thất bại",
              text: lookupErrorCode(s["error"]),
              type: "error"
            });
          } else {
            setTimeout(() => window.location.reload(), 2000)
          }
        }, (e) => {
          this.$notify({
            title: "Nộp bài thất bại",
            text: e.message,
            type: "error"
          });
        })
      }
      f();
    },
    scrollToQuestion(i) {
      this.$refs.questionContainer.children.item(i).scrollIntoView({ behavior: 'smooth' });
    },
    getResultBackground(c, i, j) {
      if(c.hasOwnProperty("answerAccuracy")){
        if(c.answerSheet[i] === -1) {
          return ""
        } else if (c.answerSheet[i] === j) {
          if (c.answerAccuracy[i] === true) {
            return "bg-emerald-200"
          } else {
            return "bg-rose-200"
          }
        }
      }
      return ""
    },
    loadContestStats() {
      server.getContestStats(this.eventId, auth.getToken()).then(s => {
        if (s.hasOwnProperty("error")) {
          this.$notify({
            title: "Tải thống kê thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        } else {
          s.forEach((v) => {
            this.option.xAxis.data = this.option.xAxis.data.concat(v.rank)
            this.option.series[0].data = this.option.series[0].data.concat(v.count)
          })
        }
      }, (e) => {
        this.$notify({
          title: "Tải thống kê thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    loadContestSessions() {
      server.loadContestSessions(this.eventId, 30, 0, this.$root.profile.email, false, [], auth.getToken()).then(s => {
        if(s.hasOwnProperty("error")) {
          this.$notify({
            title: "Tải bài thi thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        } else {
          this.contestSessions = s.contestSessions.map((v) => {
            v = utils.parseContestSession(v)
            if (!v.finished && this.activeContestSession === undefined) {
              this.activeContestSession = v
            }
            return v
          })
          this.$refs.loadingState.deactivate()
          if(this.activeContestSession === undefined) {
            this.loadContestStats()
            return
          }
          this.countdown = setInterval(() => {
            if(this.endingContest) return
            this.timeLeft = Math.max(0, this.activeContestSession.endTime - new Date().getTime())
            if(this.timeLeft === 0) {
              this.forceSubmitContest()
            }
          }, 1000);
          this.autoSaver = setInterval(() => {
            if(this.endingContest) return
            this.savingContest = true
            server.submitContestSession(this.activeContestSession.id, this.activeContestSession.answerSheet, true, auth.getToken()).then((s) => {
              if(s.hasOwnProperty("error") || (s.hasOwnProperty("success") && !s["success"])) {
                this.$notify({
                  title: "Lưu bài thất bại",
                  text: s["error"],
                  type: "error"
                });
              } else {
                this.savingContest = false
              }
            }, (e) => {
              this.$notify({
                title: "Lưu bài thất bại",
                text: e.message,
                type: "error"
              });
            })
          }, 10000);
        }
      }, (e) => {
        this.$notify({
          title: "Tải bài thi thất bại",
          text: e.message,
          type: "error"
        });
      })
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
      server.loadEvent(this.eventId, auth.getToken()).then(q => {
        if(!q.hasOwnProperty("error")) {
          if (q.hasOwnProperty("contest")) {
            q.contest.dataSheet = JSON.parse(q.contest.dataSheet)
            this.event = q;
            this.loadContestSessions()
          } else {
            this.$notify({
              title: "Lỗi hệ thống",
              text: "Hãy báo cáo với quản trị viên!",
              type: "error"
            });
          }
        } else {
          this.$notify({
            title: "Nộp bài thất bại",
            text: lookupErrorCode(q["error"]),
            type: "error"
          });
        }
      }, (e) => {
        this.$notify({
          title: "Tải sự kiện thất bại",
          text: e.message,
          type: "error"
        });
      })
    } else {
      this.$notify({
        title: "Lỗi hệ thống",
        text: "Hãy báo cáo với quản trị viên!",
        type: "error"
      });
    }
  }
}
</script>
