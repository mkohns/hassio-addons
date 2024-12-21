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
    <SideMenu :items="menuItems" :open="showActions" @trigger="menuTrigger" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import axios from "axios";
import { useSlideshowStore } from "@/stores/slideshow";
import SideMenu from "@/components/sideMenu.vue";

const menuItems = [
  {
    icon: "mdi-close",
    text: "Close Menu",
    event: "close",
  },
  {
    icon: "mdi-eye",
    text: "Pause Image",
    event: "pause",
  },
  {
    icon: "mdi-star",
    text: "Mark as Favorite",
    event: "favorite",
  },
  {
    icon: "mdi-delete",
    text: "Delete Image",
    event: "delete",
  },
  {
    icon: "mdi-image-multiple-outline",
    text: "Manage Images",
    event: "manage",
  },
  {
    icon: "mdi-cog",
    text: "Slideshow Settings",
    event: "settings",
  },
  {
    icon: "mdi-information-outline",
    text: "Info",
    event: "info",
  },
];

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

function menuTrigger(evt) {
  switch (evt) {
    case "close":
      close(evt);
      break;
    case "pause":
      pause(evt);
      break;
    case "favorite":
      favorite(evt);
      break;
    case "delete":
      deleteImage(evt);
      break;
    case "manage":
      manage(evt);
      break;
    case "settings":
      settings(evt);
      break;
    case "info":
      info(evt);
      break;
    default:
      console.log("Unknown event: " + evt);
  }
}

function close(evt) {
  toggleActions();
}

function pause(evt) {
  console.log("Pausing image");
}

function favorite(evt) {
  console.log("Marking as favorite");
}

function deleteImage(evt) {
  console.log("Deleting image");
}

function manage(evt) {
  console.log("Managing images");
}

function settings(evt) {
  console.log("Slideshow settings");
}

function info(evt) {
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
