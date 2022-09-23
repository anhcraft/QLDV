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
    <swiper-slide v-for="user in featuredUsers">
      <div class="flex flex-row gap-5 my-10">
        <div>
          <div class="bg-slate-300 w-32 h-32"></div>
        </div>
        <div>
          <p class="text-xl">{{ user.name }}</p>
          <p class="italic text-sm">- Lớp {{ user.class }}</p>
          <ul class="list-disc list-inside mt-2">
            <li v-for="val in user.achievements">{{ val.title }} ({{ val.year }})</li>
          </ul>
        </div>
      </div>
    </swiper-slide>
  </swiper>
</template>

<script>
import {Swiper, SwiperSlide} from "swiper/vue";
import {Autoplay, Pagination} from "swiper";
import "swiper/css";
import "swiper/css/pagination";
import UserAPI from "../../api/user-api";
import {ServerError} from "../../api/server-error";

export default {
  name: "KeyMemberSlideshow",
  components: {
    Swiper,
    SwiperSlide,
  },
  data() {
    return {
      featuredUsers: []
    }
  },
  setup() {
    return {
      modules: [Pagination, Autoplay]
    }
  },
  methods: {
    loadFeaturedMembers() {
      UserAPI.listFeaturedUsers().then(data => {
        if (data instanceof ServerError) {
          this.$root.popupError(data)
          return
        }
        this.featuredUsers = data
      })
    }
  },
  mounted() {
    const f = () => {
      this.loadFeaturedMembers()
    }
    this.$root.pushQueue(f.bind(this))
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
