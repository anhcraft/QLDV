<template>
  <Header></Header>
  <section class="page-section px-10 py-8 lg:py-16">
    <div class="centered-horizontal mb-5 gap-3">
      <UserGroupIcon class="w-8 h-8 text-rose-500"></UserGroupIcon>
      <p class="text-3xl font-heading">Ban chấp hành</p>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-7 md:gap-16 mt-5">
      <div class="col-span-5">
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
    </div>
  </section>
  <Footer></Footer>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
import {UserGroupIcon} from "@heroicons/vue/24/solid";
import Footer from "../components/Footer.vue";

export default {
  name: "CommitteePage",
  components: {
    Header, Footer, FloatingMenu, Breadcrumb, UserGroupIcon
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
