<template>
  <div class="image-container" @click="toggleActions">
    <img
      ref="img1"
      style="opacity: 0"
      :src="imageUrl1"
      class="full-size-image"
    />
    <img
      ref="img2"
      style="opacity: 0"
      :src="imageUrl2"
      class="full-size-image"
    />
    <div v-if="showOverlay && slide" class="overlay">
      <p v-if="slide.Message">Message: {{ slide.Message }}</p>
      <p>Send By: {{ slide.CreatedBy }}</p>
      <p>Send At: {{ formattedCreatedAt }}</p>
    </div>
    <div class="menu-container">
      <transition
        enter-active-class="animated animate__slideInLeft"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="close" class="menu-item-top" v-if="showActions">
          <v-btn size="x-large" icon="mdi-close"></v-btn
          ><span class="menu-text">Close Menu</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item1"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="pause" class="menu-item" v-if="showActions">
          <v-btn size="x-large" icon="mdi-eye"></v-btn
          ><span class="menu-text">Pause Image</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item3"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="favorite" class="menu-item" v-if="showActions">
          <v-btn size="x-large" icon="mdi-star"></v-btn
          ><span class="menu-text">Mark as Favorite</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item4"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="deleteImage" class="menu-item" v-if="showActions">
          <v-btn size="x-large" icon="mdi-delete"></v-btn
          ><span class="menu-text">Delete Image</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item5"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="manage" class="menu-item" v-if="showActions">
          <v-btn size="x-large" icon="mdi-image-multiple-outline"></v-btn
          ><span class="menu-text">Manage Images</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item6"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="settings" class="menu-item" v-if="showActions">
          <v-btn size="x-large" icon="mdi-cog"></v-btn
          ><span class="menu-text">Slideshow Settings</span>
        </div>
      </transition>

      <transition
        enter-active-class="animated animate__slideInLeft item7"
        leave-active-class="animated animate__slideOutLeft"
      >
        <div @click="info" class="menu-item-bottom" v-if="showActions">
          <v-btn size="x-large" icon="mdi-information-outline"></v-btn
          ><span class="menu-text">Info</span>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import axios from "axios";
import { useSlideshowStore } from "@/stores/slideshow";

const store = useSlideshowStore();

let interval = null;
let showOverlay = null;

let url = import.meta.env.VITE_BASE_URL;
if (url === undefined) {
  url = "";
}

const img1 = ref(null);
const img2 = ref(null);

let intervalId;
const imageUrl1 = ref("");
const imageUrl2 = ref("");
const slide = ref(null);
const currentImage = ref(1);

const showActions = ref(false);

function toggleActions() {
  if (!showActions.value) {
    console.log("Opening menu");
    showActions.value = true;
    stopSlideshow();
  } else {
    console.log("Closing menu");
    showActions.value = false;
    startSlideshow();
  }
}

function close(evt) {
  evt.stopPropagation();
  toggleActions();
}

function pause(evt) {
  evt.stopPropagation();
  console.log("Pausing image");
}

function favorite(evt) {
  evt.stopPropagation();
  console.log("Marking as favorite");
}

function deleteImage(evt) {
  evt.stopPropagation();
  console.log("Deleting image");
}

function manage(evt) {
  evt.stopPropagation();
  console.log("Managing images");
}

function settings(evt) {
  evt.stopPropagation();
  console.log("Slideshow settings");
}

function info(evt) {
  evt.stopPropagation();
  console.log("Info");
}

const formattedCreatedAt = computed(() => {
  if (!slide.value) return "";
  const date = new Date(slide.value.CreatedAt);
  return date.toLocaleString();
});

const fetchNextSlide = async () => {
  try {
    const response = await axios.get(url + "/nextslide"); // Fetch the next slide
    slide.value = response.data;
    const newImageUrl = `${url}/${slide.value.ImageURL}`; // Construct the image URL

    if (currentImage.value === 1) {
      imageUrl2.value = newImageUrl;
      img1.value.style.opacity = 0.0;
      img2.value.style.opacity = 1.0;
      currentImage.value = 2;
    } else {
      imageUrl1.value = newImageUrl;
      img1.value.style.opacity = 1.0;
      img2.value.style.opacity = 0.0;
      currentImage.value = 1;
    }
  } catch (error) {
    console.error("Error fetching next slide:", error);
  }
};

onMounted(() => {
  startSlideshow();
});

function startSlideshow() {
  fetchNextSlide();
  interval = store.interval * 1000;
  showOverlay = store.showOverlay;
  console.log("Starting slideshow");
  console.log("Setting up interval", interval);
  intervalId = setInterval(fetchNextSlide, interval);
}

function stopSlideshow() {
  console.log("Stopping slideshow");
  clearInterval(intervalId);
}

onUnmounted(() => {
  stopSlideshow();
});
</script>

<style scoped>
.menu-container {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 0; /* Adjust as needed */
}
.animated {
  -webkit-animation-duration: 1s;
  animation-duration: 1s;
  -webkit-animation-fill-mode: both;
  animation-fill-mode: both;
}
.item1 {
  -webkit-animation-delay: 0.2s;
  animation-delay: 0.2s;
}
.item2 {
  -webkit-animation-delay: 0.3s;
  animation-delay: 0.3s;
}
.item3 {
  -webkit-animation-delay: 0.4s;
  animation-delay: 0.4s;
}
.item4 {
  -webkit-animation-delay: 0.5s;
  animation-delay: 0.5s;
}
.item5 {
  -webkit-animation-delay: 0.6s;
  animation-delay: 0.6s;
}
.item6 {
  -webkit-animation-delay: 0.7s;
  animation-delay: 0.7s;
}
.item7 {
  -webkit-animation-delay: 0.8s;
  animation-delay: 0.8s;
}

.menu-text {
  margin-left: 10px;
  font-size: 1.3rem;
  font-weight: bold;
  color: black;
}

.menu-item {
  background-color: rgba(255, 255, 255, 0.346);
  padding: 10px;
}

.menu-item-top {
  border-radius: 0 50px 0 0;
  background-color: rgba(255, 255, 255, 0.346);
  padding: 10px;
}

.menu-item-bottom {
  border-radius: 0 0 50px;
  background-color: rgba(255, 255, 255, 0.346);
  padding: 10px;
}

.full-size-image {
  width: 100vw;
  height: 100vh;
  object-fit: contain;
  position: absolute;
  top: 0;
  left: 0;
  transition: opacity 4s ease-in-out;
}

.overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 10px;
  box-sizing: border-box;
}
</style>
