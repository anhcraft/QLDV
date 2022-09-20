<template>
  <Header></Header>
  <div class="py-32">
    <img src="../assets/tree-drawing.svg" class="m-auto max-w-[10rem] md:max-w-xs"  alt=""/>
    <div class="text-center font-light font-serif mt-10">
      <p class="text-3xl mb-3">Trang này không tồn tại!</p>
      <p class="text-2xl text-cyan-600 hover:underline">
        <router-link to="/">Quay về trang chủ</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import Header from "../components/Header.vue";

export default {
  name: "NotFoundPage",
  components: {
    Header
  }
}
</script>
