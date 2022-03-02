<template>
  <Header></Header>
  <div class="max-w-[1024px] m-auto pb-16">
    <Breadcrumb text="Quản lý cuộc thi" link="/em" class="mb-10"></Breadcrumb>
    <LoadingState ref="loadingState">
      <header class="border-b-2 border-b-slate-300 pb-3 text-xl flex flex-row gap-2">
        <div class="grow">{{ event.title }}</div>
        <button class="bg-emerald-600 hover:bg-emerald-700 cursor-pointer px-4 py-2 text-white text-center text-sm" @click="saveChanges">Xem kết quả</button>
        <button class="bg-rose-400 hover:bg-rose-500 cursor-pointer px-4 py-2 text-white text-center text-sm" @click="this.$refs.removePrompt.toggle()">Xóa cuộc thi</button>
      </header>
      <div class="mt-5 border-b-2 border-b-slate-300 pb-3">
        <div class="flex flex-row gap-5 place-items-center">
          <p class="font-bold">Chấp nhận câu trả lời</p>
          <input type="checkbox" class="w-4 h-4" v-model="event.contest.acceptingAnswers">
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p class="font-bold">Giới hạn số câu hỏi</p>
          <input type="number" min="3" max="100" class="border border-slate-300 px-1" v-model.number="event.contest.limitQuestions">
        </div>
        <div class="flex flex-row gap-5 place-items-center">
          <p class="font-bold">Giới hạn thời gian làm</p>
          <input type="number" min="0" max="1440" class="border border-slate-300 px-1" v-model.number="event.contest.limitTime"> phút
        </div>
        <p class="italic">(Điền 0 để quy định không giới hạn về thời gian)</p>
        <div class="mt-5">
          <p class="font-bold">Thông tin cuộc thi</p>
          <Editor
            apiKey="r7g4lphizuprqmrjv0ooj15pn5qpcesynrg101ekc40avzlg"
            :init="{
                  height: 300,
                  plugins: [
                    'advlist autolink link image lists charmap print preview hr anchor pagebreak',
                    'searchreplace wordcount visualblocks visualchars code fullscreen insertdatetime media nonbreaking',
                    'table emoticons template paste help'
                  ],
                  toolbar: 'undo redo | styleselect | bold italic | forecolor backcolor emoticons link | ' +
                    'bullist numlist outdent indent | media | ' +
                    'alignleft aligncenter alignright alignjustify | help',
                  menubar: false
                }"
            v-model="event.contest.info"
          ></Editor>
        </div>
      </div>
      <div class="mt-5 text-sm border-2 border-dashed border-slate-500 p-5">
        <p class="font-bold">Tải file excel chứa câu hỏi và đáp án:</p>
        <p class="italic">- Cột đầu tiên: câu hỏi</p>
        <p class="italic">- Cột thứ hai: đáp án (a, b, c, d hoặc 1, 2, 3, 4)</p>
        <p class="italic">- 4 cột tiếp theo: các lựa chọn</p>
        <p>Người tham gia thi sẽ nhận được các câu hỏi ngẫu nhiên trong bộ câu hỏi trên.</p>
        <input @change="onUploadDataFile" class="block mt-5 px-3 py-1 text-gray-700 text-sm bg-white border border-solid border-gray-300 rounded focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none" type="file">
      </div>
      <div class="mt-5" v-if="event.contest.dataSheet.length > 0">
        <p class="font-bold">Bộ câu hỏi</p>
        <table class="table-fixed border-collapse border border-slate-300 w-full mt-5">
          <thead>
            <tr>
              <th class="border border-slate-300">Câu hỏi</th>
              <th class="border border-slate-300">Lựa chọn 1</th>
              <th class="border border-slate-300">Lựa chọn 2</th>
              <th class="border border-slate-300">Lựa chọn 3</th>
              <th class="border border-slate-300">Lựa chọn 4</th>
            </tr>
          </thead>
          <tbody class="text-sm">
            <tr v-for="item in event.contest.dataSheet">
              <td class="border border-slate-300 overflow-auto p-2">{{ item.question }}</td>
              <td v-for="(choice, i) in item.choices" class="border border-slate-300 p-2 overflow-auto" :class="{'bg-yellow-400' : i === item.answer}">{{ choice }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="mt-5">
        <button class="bg-pink-400 hover:bg-pink-500 cursor-pointer px-4 py-2 text-white text-center text-sm" @click="saveChanges">Lưu thay đổi</button>
      </div>
   </LoadingState>
  </div>
  <FloatingMenu></FloatingMenu>
  <Prompt @callback="removeContestCallback" ref="removePrompt">
    <p class=font-bold>Bạn có muốn xóa cuộc thi này?</p><br> {{ event.title }}
  </Prompt>
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
  "name": "ContestManage",
  components: {LoadingState, Header, FloatingMenu, Breadcrumb, Datepicker, Prompt, Editor},
  data() {
    return {
      event: {}
    }
  },
  methods: {
    parseAnswer(s) {
      s = String(s).trim().toLowerCase();
      if(s === "1" || s === "a") return 0;
      if(s === "2" || s === "b") return 1;
      if(s === "3" || s === "c") return 2;
      if(s === "4" || s === "d") return 3;
    },
    async onUploadDataFile(e) {
      const file = e.target.files[0];
      const data = await file.arrayBuffer();
      const workbook = XLSX.read(data);
      if(workbook.SheetNames.length > 0) {
        const sheet = workbook.Sheets[workbook.SheetNames[0]];
        const data = XLSX.utils.sheet_to_json(sheet, {header: 1});
        const preview = [];
        for(let i = 0; i < data.length; i++){
          const row = data[i];
          preview.push({
            question: row[0],
            answer: this.parseAnswer(row[1]),
            choices: [row[2], row[3], row[4], row[5]]
          });
        }
        this.event.contest.dataSheet = preview;
      }
    },
    removeContestCallback() {
      this.$refs.loadingState.activate()
      server.removeContest(this.$route.params.id, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$refs.loadingState.deactivate()
          this.$router.push('/em/')
        } else {
          alert(`Lỗi xóa cuộc thi: ${s["error"]}`)
        }
      })
    },
    saveChanges(){
      this.$refs.loadingState.activate()
      server.changeContest(this.$route.params.id, this.event.contest, auth.getToken()).then(s => {
        if(!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.$refs.loadingState.deactivate()
          this.$router.push('/em/')
        } else {
          alert(`Lỗi lưu cuộc thi: ${s["error"]}`)
        }
      })
    }
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
            s.contest.limitTime /= 60000
          } else {
            s["contest"] = {
              acceptingAnswers: false,
              limitQuestions: 10,
              limitTime: 900000,
              dataSheet: [],
              info: ""
            }
          }
          this.event = s;
          this.$refs.loadingState.deactivate()
        } else {
          alert(`Lỗi tải sự kiện: ${s["error"]}`)
        }
      });
    } else {
      this.$refs.loadingState.deactivate()
    }
  }
}
</script>
