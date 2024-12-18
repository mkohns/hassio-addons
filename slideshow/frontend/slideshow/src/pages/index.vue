<template>
  <div class="image-container" @click="press4times">
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
    <div v-if="showOverlay && filename" class="overlay">
      <p v-if="message">Message: {{ message }}</p>
      <p>Send By: {{ createdBy }}</p>
      <p>Send At: {{ formattedCreatedAt }}</p>
    </div>
  </div>
  <v-dialog v-model="showMenu" persistent max-width="290">
    <v-card>
      <v-card-title>🚀 Menu</v-card-title>
      <v-card-text>
        <v-btn width="100%" color="primary" @click="$router.push('/manage')"
          >Manage Pictures</v-btn
        >
        <v-btn
          width="100%"
          class="mt-3"
          @click="$router.push('/config')"
          color="primary"
          >Configuration</v-btn
        >
        <v-btn
          width="100%"
          class="mt-3"
          @click="showMenu = false"
          color="primary"
          >close</v-btn
        >
      </v-card-text>
    </v-card>
  </v-dialog>
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
const filename = ref("");
const message = ref("");
const createdBy = ref("");
const createdAt = ref("");
const currentImage = ref(1);

let pressCount = 0;
const showMenu = ref(false);

function press4times() {
  if (pressCount === 0) {
    setTimeout(() => {
      console.log("Resetting press count");
      pressCount = 0;
    }, 2000);
  }
  pressCount++;
  if (pressCount === 4) {
    console.log("Yeah, you got it");
    showMenu.value = true;
  }
}

const formattedCreatedAt = computed(() => {
  if (!createdAt.value) return "";
  const date = new Date(createdAt.value);
  return date.toLocaleString();
});

const fetchNextSlide = async () => {
  try {
    const response = await axios.get(url + "/nextslide"); // Fetch the next slide
    const data = response.data;
    filename.value = data.Filename;
    message.value = data.Message;
    createdBy.value = data.CreatedBy;
    createdAt.value = data.CreatedAt;
    const newImageUrl = `${url}/${data.ImageURL}`; // Construct the image URL

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
  interval = store.interval * 1000;
  showOverlay = store.showOverlay;

  fetchNextSlide(); // Fetch the next slide immediately
  console.log("Setting up interval", interval);
  intervalId = setInterval(fetchNextSlide, interval); // Set up the interval
});

onUnmounted(() => {
  clearInterval(intervalId);
});
</script>

<style scoped>
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
