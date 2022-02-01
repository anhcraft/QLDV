<template>
  <Header></Header>
  <div class="pb-16 max-w-[1024px] m-auto">
    <Breadcrumb text="Trang cá nhân" link="/u"></Breadcrumb>
    <div class="grid grid-cols-3 gap-16 mt-5">
      <div class="col-span-1 shadow-lg shadow-slate-400">
        <section>
          <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÔNG TIN</div>
          <div class="px-5 py-3">
            <p>Tên: {{ $root.profile.name }}</p>
            <p>Lớp: {{ $root.profile.class }}</p>
            <p>Giới tính: {{ $root.profile.gender ? "Nữ" : "Nam" }}</p>
            <p>Đoàn viên: {{ $root.profile.certified ? "Đã kết nạp" : "Không" }}</p>
            <p>Niên khóa: {{ $root.profile.entryDate }} - {{ $root.profile.endDate }}</p>
          </div>
        </section>
        <section>
          <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">THÀNH TÍCH</div>
          <ul class="px-5 py-3 list-disc list-inside">
            <li v-for="a in $root.profile.achievements">{{ a }}</li>
          </ul>
        </section>
        <section>
          <div class="border-l-4 border-l-emerald-400 bg-emerald-100 px-4 py-2 shadow-lg shadow-slate-300">XẾP LOẠI</div>
          <ul class="px-5 py-3 list-disc list-inside">
            <li v-for="(value, name) in $root.profile.rates">
              {{ value === 1 ? "Tốt" : (value === 2 ? "Khá" : "-") }} ({{ name }} - {{ parseInt(name) + 1 }})
            </li>
          </ul>
        </section>
      </div>
      <div class="col-span-2">
        <section v-if="profileCoverUploading || $root.profile.profileCover === undefined" class="border-4 border-dashed border-black py-10">
          <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </section>
        <section v-else class="w-full inline-block relative overflow-hidden hover:opacity-80 border-4 border-dashed border-white hover:border-black">
          <div :style="{ 'background-image': 'url(' + $root.profile.profileCover + ')' }" class="w-full h-64 bg-contain bg-center bg-no-repeat" />
          <input type="file" class="absolute left-0 top-0 opacity-0 h-64 w-full cursor-pointer" @change="onProfileCoverChange" accept="image/*" />
        </section>
        <section v-if="submittingBoard" class="border-4 border-dashed border-black py-10">
          <svg class="animate-spin h-8 w-8 text-sky-400 m-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </section>
        <section v-else class="mt-10">
          <Editor
              apiKey="r7g4lphizuprqmrjv0ooj15pn5qpcesynrg101ekc40avzlg"
              :init="{
                  height: 500,
                  plugins: [
                    'advlist autolink link image lists charmap print preview hr anchor pagebreak',
                    'searchreplace wordcount visualblocks visualchars code fullscreen insertdatetime media nonbreaking',
                    'table emoticons template paste help'
                  ],
                  toolbar: 'undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | ' +
                    'bullist numlist outdent indent | link media | ' +
                    'forecolor backcolor emoticons | help',
                  menubar: false
                }"
              v-model="$root.profile.profileBoard"
          ></Editor>
          <button class="bg-pink-400 hover:bg-pink-500 cursor-pointer px-4 py-2 text-white text-center text-sm mt-5" @click="saveBoard">Sửa tường nhà</button>
        </section>
      </div>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import server from "../api/server";
import auth from "../api/auth";
import Editor from '@tinymce/tinymce-vue'

export default {
  name: "Profile",
  components: {Breadcrumb, Header, FloatingMenu, Editor},
  data() {
    return {
      profileCoverUploading: false,
      submittingBoard: false
    }
  },
  methods: {
    onProfileCoverChange(e) {
      if (e.target.files.length > 0) {
        this.profileCoverUploading = true
        server.setProfileCover(e.target.files[0], auth.getToken()).then(s => {
          if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
            this.profileCoverUploading = false
            window.location.reload()
          } else {
            alert(`Lỗi lưu ảnh bìa: ${s["error"]}`)
          }
        })
      }
    },
    saveBoard(){
      this.submittingBoard = true
      server.setProfileBoard(this.$root.profile.profileBoard, auth.getToken()).then(s => {
        if (!s.hasOwnProperty("error") && s.hasOwnProperty("success") && s["success"]) {
          this.submittingBoard = false
        } else {
          alert(`Lỗi lưu tường nhà: ${s["error"]}`)
        }
      })
    }
  },
  mounted() {
    if(!this.$root.progressionLoaded) {
      this.$root.loadProgression()
    }
  }
}
</script>
