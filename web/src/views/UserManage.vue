<template>
  <Header></Header>
  <div class="pb-16 max-w-[1024px] m-auto p-5 md:px-10">
    <Breadcrumb text="Quản lý thành viên" link="/um"></Breadcrumb>
    <div class="w-full h-48 my-10 border-slate-400 border-2" v-if="$root.profile.admin">
      <v-chart class="chart" :option="option" />
    </div>
    <div class="overflow-auto">
      <table class="w-max">
        <thead>
          <tr>
            <th>Tên</th>
            <th>Email</th>
            <th class="w-32">Lớp</th>
            <th>N.Sinh</th>
            <th>Giới</th>
            <th>ĐV</th>
            <th>BT</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full" v-model="filter.name"></td>
            <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full" v-model="filter.email"></td>
            <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full" v-model="filter.class" v-if="$root.profile.admin"></td>
            <td></td>
            <td></td>
            <td>
              <select v-model="filter.certified" class="bg-white">
                <option v-for="option in filter.certified_options" v-bind:value="option.value">
                  {{ option.text }}
                </option>
              </select>
            </td>
            <td></td>
          </tr>
          <tr class="border-b-2 border-b-slate-400">
            <td colspan="5" class="text-sm italic">Đang hiện {{ this.users.length }} thành viên, trong đó có {{ this.users.filter(u => u.gender).length }} nữ. Tổng cộng có {{ this.users.filter(u => u.certified).length }} đoàn viên.</td>
            <td><button class="btn-info" @click="search">Tìm & lọc</button></td>
            <td><button class="btn-success" @click="saveChanges" :class="{'opacity-20' : sumChanges === 0}">Lưu ({{ sumChanges }})</button></td>
          </tr>
          <tr v-for="user in users" class="text-sm hover:bg-blue-200" :class="selectedUser === user.email ? 'border-2 border-gray-400' : (user.certified ? '' : 'bg-red-200')">
            <td @click="selectUser(user)" class="flex flex-row cursor-pointer text-base hover:underline" :class="user.admin ? 'font-bold text-red-500' : (user['mod'] ? 'text-emerald-500' : '')">{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.class }}</td>
            <td>{{ new Intl.DateTimeFormat("vi-VN" , {dateStyle: "short"}).format(new Date(user.birth)) }}</td>
            <td class="text-center">{{ user.gender ? "Nữ" : "Nam" }}</td>
            <td>
              <BadgeCheckIcon class="w-6 m-auto" :class="user.certified ? 'text-sky-400' : 'text-gray-400'" @click="toggleCertified(user)"></BadgeCheckIcon>
            </td>
            <td>
              <StarIcon class="w-6 cursor-pointer" :class="user.mod ? 'text-emerald-500' : 'text-white'" @click="toggleMod(user)" v-if="$root.profile.admin && !user.admin"></StarIcon>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <LoadingState ref="loadingStateForUserList">
      <div class="mt-5" v-if="!userAvailable">Đã tải hết thành viên.</div>
    </LoadingState>
  </div>

  <LoadingState ref="loadingStateForUserProgression" hidden>
    <div v-if="selectedUser !== undefined">
      <div class="bg-black opacity-75 fixed top-0 left-0 w-screen h-screen" @click="selectUser(undefined)"></div>
      <div class="fixed right-0 top-0 z-10 bg-white w-full md:w-auto h-screen overflow-auto border-l-2 border-l-slate-300 p-10">
        <ChevronDoubleRightIcon class="w-8 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-1" @click="selectUser(undefined)"></ChevronDoubleRightIcon>
        <p class="my-5 font-bold">{{ selectedUser }}</p>
        <router-link target="_blank" :to="'/u/' + selectedUser.substring(0, selectedUser.search('@'))">
          <button class="btn-info">Xem trang cá nhân</button>
        </router-link>
        <div class="border-t-2 border-t-slate-300 mt-10">
          <section class="mt-5">
            <p class="text-xl">Xếp hạng</p>
            <ul class="list-disc list-inside">
              <li v-for="(value, name) in this.userProgression.rates">
                <select v-model="this.userProgression.rates[name]" class="bg-white">
                  <option v-for="option in this.rateOptions" v-bind:value="option.value">
                    {{ option.text }}
                  </option>
                </select>
                ({{ name }} - {{ parseInt(name) + 1 }})
              </li>
            </ul>
          </section>
          <section class="mt-5" v-if="this.userProgression.achievements.length > 0">
            <div class="text-xl flex flex-row gap-1">
              <p>Thành tích</p>
              <PlusCircleIcon class="w-6 cursor-pointer text-slate-500" @click="addAchievementSlot"></PlusCircleIcon>
            </div>
            <ul class="list-disc list-inside">
              <li v-for="value in this.userProgression.achievements">
                <input type="text" v-model="value.title"> (
                <select v-model="value.year" class="bg-white">
                  <option v-for="option in this.achievementOption" v-bind:value="option">
                    {{ option }}
                  </option>
                </select>)
              </li>
            </ul>
          </section>
          <button class="btn-success mt-5" @click="saveProgressionChanges">Lưu lại</button>
        </div>
      </div>
    </div>
  </LoadingState>

  <FloatingMenu></FloatingMenu>
</template>

<script>
import {
  BadgeCheckIcon,
  ChevronDoubleRightIcon,
  PlusCircleIcon,
  StarIcon
} from '@heroicons/vue/solid'
import server from "../api/server";
import auth from "../api/auth";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent
} from "echarts/components";
import {use} from "echarts/core";
import VChart from "vue-echarts";
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import LoadingState from "../components/LoadingState.vue";
import lookupErrorCode from "../api/errorCode";

use([
  CanvasRenderer,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent
])

export default {
  name: "UserManage",
  components: {
    LoadingState, Header, FloatingMenu, Breadcrumb,
    BadgeCheckIcon, StarIcon, VChart, ChevronDoubleRightIcon, PlusCircleIcon
  },
  data() {
    return {
      userAvailable: true,
      users: [],
      userProgression: {},
      selectedUser: undefined,
      dataOffset: 0,
      certChanges: {},
      modChanges: {},
      savingUserChanges: false,
      filter: {
        name: "",
        email: "",
        class: "",
        certified: 0,
        certified_options: [
          { text: 'Tất cả', value: 0 },
          { text: 'Đoàn viên', value: 1 },
          { text: 'Thanh niên', value: 2 }
        ]
      },
      rateOptions: [
        {text: '#', value: 0},
        {text: 'Tốt', value: 1},
        {text: 'Khá', value: 2}
      ],
      achievementOption: [],
      option: {
        title: {
          text: "Thống kê học sinh",
          left: "center"
        },
        tooltip: {
          trigger: "item",
          formatter: "{a} <br/>{b} : {c} ({d}%)"
        },
        series: [
          {
            name: "Theo khối",
            type: "pie",
            radius: "55%",
            center: ["25%", "50%"],
            data: [],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)"
              }
            }
          },
          {
            name: "Theo giới tính",
            type: "pie",
            radius: "55%",
            center: ["50%", "50%"],
            data: [],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)"
              }
            }
          },
          {
            name: "Theo ĐV/TN",
            type: "pie",
            radius: "55%",
            center: ["75%", "50%"],
            data: [],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)"
              }
            }
          }
        ]
      }
    }
  },
  methods: {
    loadNextUsers(){
      this.$refs.loadingStateForUserProgression.deactivate()
      this.$refs.loadingStateForUserList.activate()
      server.loadUsers(50, this.dataOffset, this.filter, auth.getToken()).then(s => {
        if(s.users.length === 0) {
          this.userAvailable = false
        } else {
          this.dataOffset += s.users.length
        }
        this.users = this.users.concat(s.users)
        this.$refs.loadingStateForUserList.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải thành viên thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    selectUser(user){
      this.userProgression = {}
      this.selectedUser = undefined
      if(user === undefined || this.$refs.loadingStateForUserProgression.loading) return;
      this.$refs.loadingStateForUserProgression.activate()
      server.loadProgression(auth.getToken(), user.email).then(s => {
        if (s.hasOwnProperty("error")) {
          this.$notify({
            title: "Tải thông tin thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
          this.userProgression = {}
          return
        }
        this.achievementOption = [user.entry, user.entry + 1, user.entry + 2]
        const map = {}
        map[user.entry] = 0
        map[user.entry + 1] = 0
        map[user.entry + 2] = 0
        this.userProgression = {
          rates: Object.assign(map, s.rates.reduce(function(map, obj) {
            map[obj["year"]] = obj["level"];
            return map;
          }, {})),
          achievements: s.achievements.concat({
            "title": "",
            "year": user.entry
          })
        }
        this.selectedUser = user.email
        this.$refs.loadingStateForUserProgression.deactivate()
      }, (e) => {
        this.$notify({
          title: "Tải thông tin thất bại",
          text: e.message,
          type: "error"
        });
      });
    },
    addAchievementSlot() {
      this.userProgression.achievements = this.userProgression.achievements.concat({
        "title": "",
        "year": this.userProgression.achievements[this.userProgression.achievements.length - 1].year
      })
    },
    toggleCertified(user) {
      if(this.savingUserChanges) return
      user.certified = !user.certified
      if(this.certChanges.hasOwnProperty(user.email)) {
        delete this.certChanges[user.email]
      } else {
        this.certChanges[user.email] = user.certified
      }
    },
    toggleMod(user) {
      if(this.savingUserChanges) return
      user.mod = !user['mod']
      if(this.modChanges.hasOwnProperty(user.email)) {
        delete this.modChanges[user.email]
      } else {
        this.modChanges[user.email] = user['mod']
      }
    },
    saveChanges() {
      if(this.savingUserChanges || this.sumChanges === 0) return
      this.savingUserChanges = true
      server.saveUserChanges({
        certified: this.certChanges,
        mod: this.modChanges
      }, auth.getToken()).then(s => {
        this.savingUserChanges = false
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.certChanges = {};
          this.modChanges = {};
          this.$notify({
            title: "Đã lưu thay đổi",
            text: "",
            type: "success"
          });
        } else {
          this.$notify({
            title: "Lưu thay đổi thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        }
      }, (e) => {
        this.$notify({
          title: "Lưu thay đổi thất bại",
          text: e.message,
          type: "error"
        });
      });
    },
    saveProgressionChanges() {
      this.$refs.loadingStateForUserProgression.activate()
      server.saveProgressionChanges(this.userProgression, this.selectedUser, auth.getToken()).then(s => {
        this.$refs.loadingStateForUserProgression.deactivate()
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$notify({
            title: "Đã lưu thay đổi",
            text: "",
            type: "success"
          });
          this.selectUser(undefined)
        } else {
          this.$notify({
            title: "Lưu thay đổi thất bại",
            text: lookupErrorCode(s["error"]),
            type: "error"
          });
        }
      }, (e) => {
        this.$notify({
          title: "Lưu thay đổi thất bại",
          text: e.message,
          type: "error"
        });
      })
    },
    search() {
      this.userAvailable = true
      this.users = []
      this.$refs.loadingStateForUserList.deactivate()
      this.dataOffset = 0
      this.certChanges = {}
      this.selectUser(undefined)
      this.loadNextUsers()
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(!this.$refs.loadingStateForUserList.loading && this.userAvailable) {
          this.loadNextUsers()
        }
      }
    }
  },
  computed: {
    sumChanges() {
      return Object.keys(this.certChanges).length + Object.keys(this.modChanges).length
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
    this.loadNextUsers()
    window.addEventListener('scroll', this.handleScroll)
    server.getUserStats(auth.getToken()).then(s => {
      if (s.hasOwnProperty("error")) {
        this.$notify({
          title: "Tải dữ liệu thống kê thất bại",
          text: lookupErrorCode(s["error"]),
          type: "error"
        });
        return
      }
      this.option.series[0].data.push({
        value: s["class10"],
        name: "10"
      })
      this.option.series[0].data.push({
        value: s["class11"],
        name: "11"
      })
      this.option.series[0].data.push({
        value: s["class12"],
        name: "12"
      })
      this.option.series[1].data.push({
        value: s["women"],
        name: "Nữ"
      })
      this.option.series[1].data.push({
        value: s["class10"] + s["class11"] + s["class12"] - s["women"],
        name: "Nam"
      })
      this.option.series[2].data.push({
        value: s["certified"],
        name: "Đoàn viên"
      })
      this.option.series[2].data.push({
        value: s["class10"] + s["class11"] + s["class12"] - s["certified"],
        name: "Thanh niên"
      })
    }, (e) => {
      this.$notify({
        title: "Tải dữ liệu thống kê thất bại",
        text: e.message,
        type: "error"
      });
    });
  }
}
</script>
