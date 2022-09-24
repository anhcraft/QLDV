<template>
  <Header></Header>
  <section class="page-section px-3 lg:px-10 py-8 lg:py-16">
    <div class="w-full h-48 border-slate-400 border-2 overflow-auto" v-if="this.$root.isGlobalManager">
      <v-chart class="chart" :option="option"/>
    </div>
    <div class="border border-slate-400 text-sm p-5" :class="{'opacity-20' : selectedBatchUsers.length === 0}">
      <p class="font-bold">Bạn đang chọn {{ selectedBatchUsers.length }} người dùng.</p>
      <p class="font-bold">- Thực hiện đổi chức vụ cho các người dùng trên sang:</p>
      <select class="border-2 border-gray-300 px-2 py-0.5"
              :disabled="updatingBatchUsers || selectedBatchUsers.length === 0" v-model.number="batchUpdateRole">
        <option v-for="v in roleTables.filter(v => v.role > 0)" :value="v.role">{{ v.name }}</option>
      </select>
      <button class="btn-outline-sm ml-3" :class="{'opacity-50' : updatingBatchUsers}" @click="batchUpdate()">Đồng ý
      </button>
    </div>
    <div class="overflow-auto mt-10">
      <table class="min-w-[48rem] w-full">
        <thead>
        <tr>
          <th></th>
          <th>Tên</th>
          <th>Email</th>
          <th class="w-16">Chi đoàn</th>
          <th class="w-24">Ngày sinh</th>
          <th class="w-16">Giới tính</th>
          <th class="w-16">SĐT</th>
          <th class="w-44">Chức vụ</th>
          <th>Thao tác</th>
        </tr>
        </thead>
        <tbody>
        <tr>
          <td></td>
          <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full"
                     v-model.trim="pagination.name"></td>
          <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full"
                     v-model.trim="pagination.email"></td>
          <td><input placeholder="..." class="border-2 border-gray-300 px-2 py-0.5 w-full"
                     v-model.trim="pagination.class" v-if="this.$root.isGlobalManager"></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <select class="border-2 border-gray-300 px-2 py-0.5 w-full" v-model.number="pagination.role">
              <option v-for="v in roleTables" :value="v.role">{{ v.name }}</option>
            </select>
          </td>
          <td>
            <button class="btn-info" @click="search">Tìm và lọc</button>
          </td>
        </tr>
        <tr class="border-b-2 border-b-slate-400 w-full">
          <td></td>
          <td colspan="8" class="text-sm italic">Đang hiện {{ this.users.length }} tài khoản, trong đó có
            {{ this.users.filter(u => u.gender === "female").length }} nữ. Tổng cộng có
            {{ this.users.filter(u => isMember(u.role)).length }} đoàn viên.
          </td>
        </tr>
        <tr v-for="(user, i) in users" class="text-sm hover:bg-blue-200 w-full"
            :class="{'bg-blue-200' : selectedBatchUsers.includes(i)}">
          <td>
            <CheckCircleIcon v-if="user.id !== $root.user.profile.id" class="w-6 cursor-pointer text-gray-300"
                             :class="{'text-blue-500' : selectedBatchUsers.includes(i)}"
                             @click="toggleSelectBatchUser(i)"></CheckCircleIcon>
          </td>
          <td class="text-base" :class="getRoleStyle(user.role)">{{ user.name }}</td>
          <td><input v-model.trim="user.email" class="w-full" readOnly></td>
          <td>{{ user.class }}</td>
          <td>{{ new Intl.DateTimeFormat("vi-VN", {dateStyle: "short"}).format(new Date(user.birthday)) }}</td>
          <td class="text-center">{{ user.gender === "female" ? "Nữ" : "Nam" }}</td>
          <td>{{ user.phone }}</td>
          <td>{{ getRoleName(user.role) }}</td>
          <td class="centered-horizontal gap-2">
            <PencilSquareIcon class="w-6 cursor-pointer text-gray-500" @click="selectUser(i)"></PencilSquareIcon>
            <StarIcon class="w-6 cursor-pointer text-gray-500"
                      :class="updatingUser ? 'opacity-50' : (user.featured ? 'text-amber-500' : '')"
                      @click="toggleFeatured(i)"></StarIcon>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <LoadingState ref="loadingStateForUserList">
      <div class="mt-5" v-if="!pagination.available">Đã tải hết thành viên.</div>
    </LoadingState>
  </section>

  <LoadingState ref="loadingStateForUserProgression" hidden>
    <div v-if="selectedUser >= 0">
      <div class="bg-black opacity-75 fixed z-[100] top-0 left-0 w-screen h-screen" @click="selectUser(-1)"></div>
      <div
          class="fixed z-[200] w-full md:w-[500px] h-screen right-0 top-0 overflow-auto bg-white border-l-2 border-l-slate-300 p-10">
        <ChevronDoubleRightIcon class="w-8 cursor-pointer border-slate-400 border-2 rounded-full text-slate-500 p-1"
                                @click="selectUser(-1)"></ChevronDoubleRightIcon>
        <p class="my-5 font-bold">{{ users[selectedUser].name }}</p>
        <router-link target="_blank"
                     :to="'/u/' + users[selectedUser].email.substring(0, users[selectedUser].email.search('@'))">
          <button class="btn-info">Xem trang cá nhân</button>
        </router-link>
        <div class="border-t-2 border-t-slate-300 mt-10">
          <section class="mt-5">
            <p class="text-xl">Xếp hạng</p>
            <ul class="list-disc list-inside">
              <li v-for="v in users[selectedUser].annualRanks">
                <select v-model.number="v.level" class="bg-white">
                  <option v-for="option in this.rateOptions" :value="option.value">{{ option.text }}</option>
                </select>
                ({{ v.year }} - {{ parseInt(v.year) + 1 }})
              </li>
            </ul>
          </section>
          <section class="mt-5">
            <div class="text-xl flex flex-row gap-1">
              <p>Thành tích</p>
              <PlusCircleIcon class="w-6 cursor-pointer text-slate-500" @click="addAchievementSlot()"></PlusCircleIcon>
            </div>
            <ul class="list-disc list-inside" v-if="users[selectedUser].hasOwnProperty('achievements')">
              <li class="flex flex-row" v-for="value in users[selectedUser].achievements">
                <input type="text" class="grow border-b border-b-slate-400" v-model="value.title"> (
                <select v-model.number="value.year" class="bg-white">
                  <option v-for="option in [0, 1, 2, 3]" :value="users[selectedUser].entryYear + option">
                    {{ users[selectedUser].entryYear + option }}
                  </option>
                </select>)
              </li>
            </ul>
          </section>
          <button class="btn-success mt-5" :class="{'opacity-50' : updatingUser}" @click="saveProgressionChanges">Lưu
            lại
          </button>
        </div>
      </div>
    </div>
  </LoadingState>
  <Footer></Footer>
</template>

<script>
import {
  CheckCircleIcon,
  ChevronDoubleRightIcon,
  PencilSquareIcon,
  PlusCircleIcon,
  StarIcon
} from '@heroicons/vue/24/solid'
import {CanvasRenderer} from "echarts/renderers";
import {PieChart} from "echarts/charts";
import {LegendComponent, TitleComponent, TooltipComponent} from "echarts/components";
import {use} from "echarts/core";
import VChart from "vue-echarts";
import Header from "../components/Header.vue";
import LoadingState from "../components/LoadingState.vue";
import {
  GetRoleGroup,
  GetRoleName,
  GetRoleTable,
  IsMember,
  RoleGroupClassManager,
  RoleGroupGlobalManager
} from "../auth/roles";
import UserAPI from "../api/user-api";
import {ServerError} from "../api/server-error";
import Footer from "../components/Footer.vue";

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
    LoadingState,
    Header,
    Footer,
    VChart,
    ChevronDoubleRightIcon,
    PlusCircleIcon,
    PencilSquareIcon,
    CheckCircleIcon,
    StarIcon
  },
  data() {
    return {
      users: [],
      pagination: {
        belowId: 0,
        filterName: "",
        filterEmail: "",
        filterClass: "",
        filterRole: 0,
        available: true
      },
      //
      selectedUser: -1,
      updatingUser: false,
      //
      selectedBatchUsers: [],
      batchUpdateRole: 0,
      processedBatchedUsers: 0,
      updatingBatchUsers: false,
      //
      rateOptions: [
        {text: '#', value: 0},
        {text: 'Xuất sắc', value: 1},
        {text: 'Khá', value: 2},
        {text: 'Trung bình', value: 3}
      ],
      option: {
        textStyle: {
          fontFamily: "sans-serif"
        },
        title: {
          text: "Thống kê",
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
            name: "Theo ĐV",
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
          }
        ]
      }
    }
  },
  methods: {
    toggleSelectBatchUser(u) {
      if (this.updatingBatchUsers || this.users[u].id === this.$root.user.profile.id) return
      if (this.selectedBatchUsers.includes(u)) {
        this.selectedBatchUsers = this.selectedBatchUsers.filter(v => v !== u)
      } else {
        this.selectedBatchUsers = this.selectedBatchUsers.concat(u)
      }
    },
    getRoleName(r) {
      return GetRoleName(r)
    },
    getRoleStyle(role) {
      if (GetRoleGroup(role) === RoleGroupClassManager) {
        return 'text-emerald-500'
      } else if (GetRoleGroup(role) >= RoleGroupGlobalManager) {
        return 'font-bold text-red-500'
      } else {
        return ''
      }
    },
    isMember(r) {
      return IsMember(r)
    },
    loadNextUsers(callback) {
      this.$refs.loadingStateForUserList.activate()
      const limit = 15
      UserAPI.listUsers({
        limit: limit,
        "below-id": this.pagination.belowId,
        "filter-name": this.pagination.filterName,
        "filter-class": this.pagination.filterClass,
        "filter-role": this.pagination.filterRole,
        "filter-email": this.pagination.filterEmail,
      }).then(res => {
        this.$refs.loadingStateForUserList.deactivate()
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          if (res.length < limit) {
            this.pagination.available = false
          }
          if (res.length > 0) {
            this.pagination.belowId = res[res.length - 1].id
            this.users = this.users.concat(res)
          }
        }
        callback.call(null)
      })
    },
    addAchievementSlot() {
      let a = this.users[this.selectedUser].achievements
      if (a === undefined) a = []
      this.users[this.selectedUser].achievements = a.concat({
        "title": "",
        "year": a.length === 0 ? this.users[this.selectedUser].entryYear : a[a.length - 1].year
      })
    },
    search() {
      if (this.$refs.loadingStateForUserList.loading || this.$refs.loadingStateForUserProgression.loading || this.updatingBatchUsers || this.updatingUser) return
      this.users = []
      this.pagination.belowId = 0
      this.pagination.available = false
      this.selectedUser = -1
      this.loadNextUsers()
    },
    handleScroll() {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if (!this.$refs.loadingStateForUserList.loading && this.pagination.available) {
          this.loadNextUsers()
        }
      }
    },
    loadStats() {
      UserAPI.getUserStats().then(res => {
        if (res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        this.option.series[0].data.push({name: "10", value: res["user-count-by-grade"]["grade-10"]})
        this.option.series[0].data.push({name: "11", value: res["user-count-by-grade"]["grade-11"]})
        this.option.series[0].data.push({name: "12", value: res["user-count-by-grade"]["grade-12"]})
        this.option.series[1].data.push({name: "Chưa kết nạp", value: res["user-count-by-role"]["regular-member"]})
        this.option.series[1].data.push({
          name: "Đã kết nạp",
          value: res["user-count-by-role"]["certified-member"] +
              res["user-count-by-role"]["class-deputy-secretary"] +
              res["user-count-by-role"]["class-secretary"]
        })
      })
    },
    selectUser(index) {
      if (this.$refs.loadingStateForUserList.loading || this.$refs.loadingStateForUserProgression.loading || this.updatingBatchUsers || this.updatingUser) return
      if (index < 0) {
        this.selectedUser = index
        return
      }
      this.$refs.loadingStateForUserProgression.activate()
      UserAPI.getUser(this.users[index].id, {
        "profile": false,
        "achievements": true,
        "annual-ranks": true
      }).then(res => {
        if (res instanceof ServerError) {
          this.$root.popupError(res)
        } else {
          this.users[index]["achievements"] = res.achievements
          const ar = {}
          ar[this.users[index].entryYear] = 0
          ar[this.users[index].entryYear + 1] = 0
          ar[this.users[index].entryYear + 2] = 0
          if (res.hasOwnProperty("annualRanks")) {
            res["annualRanks"].forEach(v => {
              if (ar.hasOwnProperty(v.year)) {
                ar[v.year] = v.level
              }
            })
          }
          this.users[index]["annualRanks"] = Object.keys(ar).map(v => {
            return {
              level: ar[v],
              year: v
            }
          })
          this.selectedUser = index
        }
        this.$refs.loadingStateForUserProgression.deactivate()
      })
    },
    batchUpdate() {
      if (this.$refs.loadingStateForUserList.loading || this.updatingUser || this.updatingBatchUsers) return
      this.updatingBatchUsers = true
      this.processedBatchedUsers = 0
      const role = this.batchUpdateRole
      const users = [...this.selectedBatchUsers] // clone
      users.map(v => {
        UserAPI.updateUser(this.users[v].id, {
          profile: {
            role: role
          },
          achievements: undefined,
          annualRanks: undefined
        }).then(res => {
          if (res instanceof ServerError) {
            this.$root.popupError(res)
          } else {
            this.selectedBatchUsers = this.selectedBatchUsers.filter(q => q !== v)
            this.users[v].role = role
          }
          if (++this.processedBatchedUsers === users.length) {
            this.updatingBatchUsers = false
          }
        })
      })
    },
    toggleFeatured(index) {
      if (this.$refs.loadingStateForUserProgression.loading || this.updatingUser || this.updatingBatchUsers) return
      this.updatingUser = true
      UserAPI.updateUser(this.users[index].id, {
        profile: {
          featured: !this.users[index].featured
        },
        achievements: undefined,
        annualRanks: undefined
      }).then(res => {
        this.updatingUser = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        this.users[index].featured = !this.users[index].featured
      })
    },
    saveProgressionChanges() {
      if (this.$refs.loadingStateForUserProgression.loading || this.updatingUser || this.updatingBatchUsers) return
      this.updatingUser = true
      UserAPI.updateUser(this.users[this.selectedUser].id, {
        profile: undefined,
        achievements: this.users[this.selectedUser].achievements,
        annualRanks: this.users[this.selectedUser].annualRanks
      }).then(res => {
        this.updatingUser = false
        if (res instanceof ServerError) {
          this.$root.popupError(res)
          return
        }
        this.$notify({
          title: "Đã lưu thay đổi",
          text: "",
          type: "success"
        });
      })
    }
  },
  computed: {
    roleTables() {
      return GetRoleTable().filter(v => v.role <= this.$root.user.profile.role)
    }
  },
  unmounted() {
    window.removeEventListener('scroll', this.handleScroll);
  },
  mounted() {
    const f = () => {
      if (!this.$root.isLoggedIn() || !this.$root.isManager) {
        this.$router.push({name: "home"})
        return
      }
      window.addEventListener('scroll', this.handleScroll)
      this.$refs.loadingStateForUserProgression.deactivate()
      this.loadNextUsers(() => {
        if (this.$root.isGlobalManager) {
          this.loadStats()
        }
      })
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
