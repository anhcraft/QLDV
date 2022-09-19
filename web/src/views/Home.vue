<template>
  <Header></Header>
  <section class="page-section px-10 py-16">
    <PostSection></PostSection>
  </section>
  <section class="bg-pattern-hideout">
    <div class="page-section px-10 py-16">
      <ActivitySection></ActivitySection>
    </div>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import PostSection from "../components/home/PostSection.vue";
import ActivitySection from "../components/home/ActivitySection.vue";
import Footer from "../components/Footer.vue";

export default {
  name: "Home",
  components: {
    Header, Footer, PostSection, ActivitySection
  }
}
</script>
