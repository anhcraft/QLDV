<template>
  <Header></Header>
  <section class="page-section px-10 py-16">
    <PostSection></PostSection>
  </section>
  <section class="bg-neutral-100 px-10 relative">
    <div class="curved-divider-top">
      <svg data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
        <path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z" class="curved-shape-fill"></path>
      </svg>
    </div>
    <div class="page-section py-32">
      <ActivitySection></ActivitySection>
    </div>
    <div class="curved-divider-bottom">
      <svg data-name="curved-divider" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
        <path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z" class="curved-shape-fill"></path>
      </svg>
    </div>
  </section>
  <Footer></Footer>
  <FloatingMenu></FloatingMenu>
</template>

<script>
import Header from "../components/Header.vue";
import FloatingMenu from "../components/FloatingMenu.vue";
import PostSection from "../components/home/PostSection.vue";
import ActivitySection from "../components/home/ActivitySection.vue";
import Footer from "../components/Footer.vue";

export default {
  name: "Home",
  components: {
    Header, Footer, FloatingMenu, PostSection, ActivitySection
  }
}
</script>
