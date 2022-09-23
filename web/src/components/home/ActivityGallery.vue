<template>
  <swiper :rewind="true"
          :autoplay="{
            delay: 2500
          }"
          :pagination="{
            clickable: true,
          }"
          :effect="'fade'"
          :lazy="true"
          :modules="modules">
    <swiper-slide v-for="(url) in images">
      <img :data-src="url" class="swiper-lazy"/>
      <div class="swiper-lazy-preloader swiper-lazy-preloader-white"></div>
    </swiper-slide>
  </swiper>
</template>

<script>
import {Swiper, SwiperSlide} from "swiper/vue";
import {Autoplay, EffectFade, Lazy, Pagination} from "swiper";
import "swiper/css";
import "swiper/css/lazy";
import "swiper/css/pagination";
import "swiper/css/effect-fade";

export default {
  name: "ActivityGallery",
  components: {
    Swiper,
    SwiperSlide,
  },
  props: {
    images: Array
  },
  setup() {
    return {
      modules: [Autoplay, Lazy, Pagination, EffectFade]
    }
  }
}
</script>

<style scoped="">
.swiper {
  width: 100%;
  height: 100%;
}

.swiper-slide {
  text-align: center;
  font-size: 18px;
  background: #fff;

  /* Center slide text vertically */
  display: -webkit-box;
  display: -ms-flexbox;
  display: -webkit-flex;
  display: flex;
  -webkit-box-pack: center;
  -ms-flex-pack: center;
  -webkit-justify-content: center;
  justify-content: center;
  -webkit-box-align: center;
  -ms-flex-align: center;
  -webkit-align-items: center;
  align-items: center;
}

.swiper-slide img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>

<style>
.swiper-pagination-bullet {
  background-color: #fff;
  opacity: 0.5;
}

.swiper-pagination-bullet-active {
  background-color: #007AFFFF !important;
  opacity: 1;
}
</style>
