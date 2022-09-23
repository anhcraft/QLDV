<template>
  <Header></Header>
  <section class="page-section px-10 py-16">
    <PostSection></PostSection>
  </section>
  <section class="bg-pattern-hideout">
    <div class="page-section px-10 py-16">
      <ActivitySection :image-gallery="imageGallery"></ActivitySection>
    </div>
  </section>
  <Footer></Footer>
</template>

<script>
import Header from "../components/Header.vue";
import PostSection from "../components/home/PostSection.vue";
import ActivitySection from "../components/home/ActivitySection.vue";
import Footer from "../components/Footer.vue";
import SettingAPI from "../api/setting-api";
import {ServerError} from "../api/server-error";

export default {
  name: "Home",
  components: {
    Header, Footer, PostSection, ActivitySection
  },
  data() {
    return {
      imageGallery: []
    }
  },
  mounted() {
    const f = () => {
      SettingAPI.getSetting("homepage").then(data => {
        if (data instanceof ServerError) {
          this.$root.popupError(data)
          return
        }
        this.imageGallery = data["activitySlideshow"]
      })
    }
    this.$root.pushQueue(f.bind(this))
  }
}
</script>
