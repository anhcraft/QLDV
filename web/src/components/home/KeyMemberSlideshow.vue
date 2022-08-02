<template>
  <swiper :rewind="true"
          :autoplay="{
            delay: 2500
          }"
          :pagination="{
            clickable: true,
          }"
          :lazy="true"
          :modules="modules">
    <swiper-slide v-for="mem in members">
      <div class="flex flex-row gap-5 my-10">
        <div>
          <div class="bg-slate-300 w-32 h-32"></div>
        </div>
        <div>
          <p class="text-2xl">{{ mem.name }}</p>
          <p class="italic text-sm">- Lớp {{ mem.className }}</p>
          <p class="mt-2">{{ mem.comment }}</p>
        </div>
      </div>
    </swiper-slide>
  </swiper>
</template>

<script>
import { Swiper, SwiperSlide } from "swiper/vue";
import { Pagination, Autoplay } from "swiper";
import "swiper/css";
import "swiper/css/pagination";

export default {
  name: "KeyMemberSlideshow",
  components: {
    Swiper,
    SwiperSlide,
  },
  data(){
    return {
      members: [
        {
          name: "Nguyễn Văn A",
          className: "12AB1",
          comment: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras aliquam velit vel ipsum elementum efficitur. Nulla porttitor id leo sit amet rhoncus. Curabitur pulvinar elit non ex commodo elementum. Duis sit amet quam aliquet, mollis metus sed, iaculis tortor."
        },
        {
          name: "Nguyễn Văn B",
          className: "12AB2",
          comment: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras aliquam velit vel ipsum elementum efficitur. Nulla porttitor id leo sit amet rhoncus. Curabitur pulvinar elit non ex commodo elementum. Duis sit amet quam aliquet, mollis metus sed, iaculis tortor."
        },
        {
          name: "Nguyễn Văn C",
          className: "12A1",
          comment: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras aliquam velit vel ipsum elementum efficitur. Nulla porttitor id leo sit amet rhoncus. Curabitur pulvinar elit non ex commodo elementum. Duis sit amet quam aliquet, mollis metus sed, iaculis tortor."
        }
      ]
    }
  },
  setup() {
    return {
      modules: [Pagination,Autoplay]
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
</style>

<style>
.swiper-pagination-bullet {
  background-color: #000;
  opacity: 0.5;
}
.swiper-pagination-bullet-active {
  background-color: #007AFFFF !important;
  opacity: 1;
}
</style>
