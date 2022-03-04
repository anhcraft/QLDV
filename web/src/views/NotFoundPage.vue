<template>
  <Header></Header>
  <div class="py-32">
    <img src="../assets/tree-drawing.svg" class="m-auto max-w-[10rem] md:max-w-xs"  alt=""/>
    <div class="text-2xl text-center font-light font-serif mt-5">
      <p>Trang này không tồn tại</p>
      <p>{{ $route.fullPath }}</p>
    </div>
  </div>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";

export default {
  name: "NotFoundPage",
  components: {
    Header, FloatingMenu
  }
}
</script>
