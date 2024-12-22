<template>
  <div class="image-container" @click="toggleActions">
    <v-chip v-if="showNewChip" size="x-large" class="newchip">🔥 NEW</v-chip>
    <div class="errorContainer" v-if="error != null" @click="toggleActions">
      <div class="error-heading">Hm, we could not find any pictures!</div>
      <img v-if="error != null" src="@/assets/404.png" class="mid-size-image" />
      <div class="error-footer">The server reported: {{ error }}</div>
    </div>
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
import { useSlideshowStore } from "@/stores/slideshow";
import SideMenu from "@/components/sideMenu.vue";
import backend from "@/api/backend";
import router from "@/router";

const menuItems = ref([
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
]);

const store = useSlideshowStore();

let interval = null;
let showOverlay = null;

let url = import.meta.env.VITE_BASE_URL;
if (url === undefined) {
  url = "";
}

const showNewChip = computed(() => {
  if (!slide.value) return false;
  if (store.showNewChip === false) return false;
  const date = new Date(slide.value.CreatedAt);
  const now = new Date();
  const diff = now - date;
  //console.log("Difference:", diff);
  return diff < 86400000; // 24 hours in milliseconds
});

const img1 = ref(null);
const img2 = ref(null);

let intervalId;
const imageUrl1 = ref("");
const imageUrl2 = ref("");
const slide = ref(null);
const error = ref(null);
const currentImage = ref(1);

const showActions = ref(false);

function toggleActions() {
  if (!showActions.value) {
    console.log("Opening menu");
    updateMenu();
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
  if (!slide.value) return;
  console.log("Pausing image");
  let item = menuItems.value.find((item) => item.event === "pause");
  if (item.text === "Pause Image") {
    item.text = "Resume Image";
    item.icon = "mdi-eye-off";
    backend
      .pauseImage(slide.value.Filename)
      .then((response) => {
        console.log("Image paused");
      })
      .catch((error) => {
        console.error("Error pausing image:", error);
      });
  } else {
    item.text = "Pause Image";
    item.icon = "mdi-eye";
    backend
      .resumeImage(slide.value.Filename)
      .then((response) => {
        console.log("Image resumed");
      })
      .catch((error) => {
        console.error("Error resuming image:", error);
      });
  }
}

function favorite(evt) {
  if (!slide.value) return;
  console.log("Marking as favorite");
  let item = menuItems.value.find((item) => item.event === "favorite");
  if (item.text === "Mark as Favorite") {
    item.text = "Unmark as Favorite";
    item.icon = "mdi-heart";
    item.color = "red";
    backend
      .like(slide.value.Filename)
      .then((response) => {
        console.log("Image liked");
      })
      .catch((error) => {
        console.error("Error liking image:", error);
      });
  } else {
    item.text = "Mark as Favorite";
    item.icon = "mdi-star";
    item.color = "white";
    backend
      .unlike(slide.value.Filename)
      .then((response) => {
        console.log("Image unliked");
      })
      .catch((error) => {
        console.error("Error unliking image:", error);
      });
  }
}

function deleteImage(evt) {
  if (!slide.value) return;
  console.log("Deleting image");
  backend
    .delete(slide.value.Filename)
    .then((response) => {
      console.log("Image deleted");
      toggleActions();
    })
    .catch((error) => {
      console.error("Error deleting image:", error);
    });
}

function manage(evt) {
  console.log("Managing images");
  router.push("/manage");
}

function settings(evt) {
  console.log("Slideshow settings");
  router.push("/config").catch((err) => {
    console.error("Error navigating to settings:", err);
  });
}

function info(evt) {
  console.log("Info");
}

const formattedCreatedAt = computed(() => {
  if (!slide.value) return "";
  const date = new Date(slide.value.CreatedAt);
  return date.toLocaleString();
});

function updateMenu() {
  if (!slide.value) return;
  let pausedItem = menuItems.value.find((item) => item.event === "pause");
  if (slide.value.Enabled) {
    pausedItem.text = "Pause Image";
    pausedItem.icon = "mdi-eye";
  } else {
    pausedItem.text = "Resume Image";
    pausedItem.icon = "mdi-eye-off";
  }

  let favoriteItem = menuItems.value.find((item) => item.event === "favorite");
  if (slide.value.Favorite) {
    favoriteItem.text = "Unmark as Favorite";
    favoriteItem.icon = "mdi-heart";
  } else {
    favoriteItem.text = "Mark as Favorite";
    favoriteItem.icon = "mdi-star";
  }
}

function fetchNextSlide() {
  backend
    .nextSlide()
    .then((response) => {
      error.value = null;
      slide.value = response.data;
      const newImageUrl = `${url}/${slide.value.ImageURL}`;
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
    })
    .catch((err) => {
      console.error("Error fetching next slide:", err);
      img1.value.style.opacity = 0.0;
      img2.value.style.opacity = 0.0;
      slide.value = null;
      console.log(err);
      error.value =
        err.response?.data || err.message || "Unknown error occurred"; // Update error message
    });
}

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
.full-size-image {
  width: 100vw;
  height: 100vh;
  object-fit: contain;
  position: absolute;
  top: 0;
  left: 0;
  transition: opacity 4s ease-in-out;
}

.mid-size-image {
  width: 50vw;
  height: 50vh;
  object-fit: contain;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
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

.newchip {
  position: absolute;
  top: 15px;
  right: 15px;
  z-index: 100;
}

.error-heading {
  position: absolute;
  font-size: 2em;
  text-align: center;
  top: 15%;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}

.error-footer {
  position: absolute;
  font-size: 1.1em;
  text-align: center;
  bottom: 20%;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}
</style>
