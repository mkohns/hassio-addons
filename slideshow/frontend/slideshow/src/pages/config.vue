<template>
  <v-container class="d-flex align-center justify-center" style="height: 100vh">
    <v-card width="80%">
      <v-card-title
        ><v-icon color="primary" class="mr-2">mdi-cog</v-icon>Slideshow
        Settings</v-card-title
      >
      <v-card-text class="mt-5">
        <v-slider
          v-model="interval"
          label="Slide Interval"
          min="1"
          max="100"
          step="1"
          hide-details="true"
          color="primary"
        >
          <template v-slot:append>
            <v-text-field
              v-model="interval"
              label="Seconds"
              density="compact"
              style="width: 80px"
              type="number"
              variant="outlined"
              hide-details
            ></v-text-field>
          </template>
        </v-slider>
        <v-switch color="primary" label="Show Overlay" v-model="overlay">
        </v-switch>
        <v-btn color="primary" @click="saveConfig">Save</v-btn>
        <v-btn class="ml-3" color="primary" @click="goBack">Back</v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useSlideshowStore } from "@/stores/slideshow";
import router from "@/router";

const store = useSlideshowStore();

let interval = ref(store.interval);
let overlay = ref(store.showOverlay);

const saveConfig = () => {
  store.setInterval(interval.value);
  store.setShowOverlay(overlay.value);
  console.log("Config saved");
};

const goBack = () => {
  router.push("/");
};

onMounted(() => {});

onUnmounted(() => {});
</script>

<style scoped></style>
