<template>
  <Header></Header>
  <div class="pb-16 max-w-[1024px] m-auto grid grid-cols-3 gap-16">
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
      <section>
        <img src="../assets/profile-cover.jpg" />
      </section>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";

export default {
  name: "Profile",
  components: {Header, FloatingMenu},
  data() {
    return {

    }
  },
  mounted() {
    if(!this.$root.progressionLoaded) {
      this.$root.loadProgression()
    }
  }
}
</script>
