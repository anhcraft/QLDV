<template>
  <Header></Header>
  <div class="py-16">
    <div class="max-w-[1024px] m-auto grid grid-cols-6 gap-24">
      <div class="col-span-4">
        <div class="border-l-4 border-l-emerald-400 font-light text-xl px-4">BAN CHẤP HÀNH ĐOÀN</div>
        <div class="flex flex-col gap-5 mt-5">
          <div v-for="i in committees">
            <div class="float-left w-32 h-32 bg-cover bg-center mr-5" :style="`background-image: url(${i.avatar})`"></div>
            <div>
              <div><span class="font-bold">Họ và tên:</span> <router-link :to="'/u/' + i.user">{{ i.name }}</router-link></div>
              <div><span class="font-bold">Chức vụ:</span> {{ i.role }}</div>
              <div><span class="font-bold">Email:</span> {{ i.user }}@dian.sgdbinhduong.edu.vn</div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-span-2">
        <Sidebar></Sidebar>
      </div>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Sidebar from "../components/Sidebar.vue";

export default {
  name: "CommitteePage",
  components: {
    Header, FloatingMenu, Sidebar
  },
  data() {
    return {
      committees: [
        {
          name: "Vũ Thị Huyền Trang",
          avatar: "http://thptdian.edu.vn/upload/104090/20210222/ae6aa2e712f4e1aab8e5_f26b3e6acd.jpg",
          role: "Bí thư Đoàn",
          user: "trangvth"
        },
        {
          name: "Phan Thị Trang",
          avatar: "http://thptdian.edu.vn/upload/104090/20210223/_DSC8628_b74db5f3ea.jpg",
          role: "Phó Bí thư Đoàn",
          user: "trangpt"
        },
        {
          name: "Huỳnh Thị Ngọc Huyền",
          avatar: "http://thptdian.edu.vn/upload/104090/20210223/_DSC8840_97ab0d3527.jpg",
          role: "Phó Bí thư Đoàn",
          user: "huyenhtn"
        }
      ]
    }
  },
  mounted() {
  }
}
</script>
