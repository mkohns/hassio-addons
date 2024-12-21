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
        <v-switch
          hide-details
          color="primary"
          label="Show Subtitle Overlay"
          v-model="overlay"
        >
        </v-switch>
        <v-switch
          hide-details
          color="primary"
          label="Show New Indicator"
          v-model="newchip"
        >
        </v-switch>
        <div class="subtitle">
          <v-icon class="mr-2" color="primary">mdi-filter</v-icon>Slide Filter
        </div>
        <v-switch
          hide-details
          color="primary"
          label="Show Only Favorites"
          v-model="showOnlyFavorites"
        >
        </v-switch>
        <v-switch
          hide-details
          color="primary"
          label="Show Only Active Images"
          v-model="showOnlyActive"
        >
        </v-switch>
        <v-row no-gutters>
          <v-col cols="6">
            <v-switch
              hide-details
              color="primary"
              label="Show Only Images in timeframe"
              v-model="showOnlyInTimeFrame"
            >
            </v-switch>
          </v-col>
          <v-col class="pl-1 pr-1">
            <DateSelect
              :disabled="!showOnlyInTimeFrame"
              label="Start Date"
              v-model="startDate"
              @update:modelValue="startDateUpdate"
            ></DateSelect>
          </v-col>
          <v-col class="pl-1 pr-1">
            <DateSelect
              nullValue="Today"
              :disabled="!showOnlyInTimeFrame"
              label="End Date"
              v-model="endDate"
              @update:modelValue="endDateUpdate"
            ></DateSelect>
          </v-col>
        </v-row>

        <div class="subtitle">
          <v-icon class="mr-2" color="primary">mdi-play-box-outline</v-icon
          >Slide Modes
        </div>
        <v-switch
          hide-details
          color="primary"
          label="Random (default)"
          v-model="modeRandom"
          @update:modelValue="
            modeChronological = false;
            modeReverseChronological = false;
            checkMode();
          "
        >
        </v-switch>
        <v-switch
          hide-details
          color="primary"
          label="Chronological"
          v-model="modeChronological"
          @update:modelValue="
            modeRandom = false;
            modeReverseChronological = false;
            checkMode();
          "
        >
        </v-switch>
        <v-switch
          hide-details
          color="primary"
          label="Reverse Chronological"
          v-model="modeReverseChronological"
          @update:modelValue="
            modeRandom = false;
            modeChronological = false;
            checkMode();
          "
        >
        </v-switch>
        <div class="mt-4">
          <v-btn color="primary" @click="saveConfig">Save</v-btn>
          <v-btn class="ml-3" color="primary" @click="goBack">Back</v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useSlideshowStore } from "@/stores/slideshow";
import router from "@/router";
import DateSelect from "@/components/dateSelect.vue";

const store = useSlideshowStore();

const interval = ref(store.interval);
const overlay = ref(store.showOverlay);
const newchip = ref(store.showNewChip);
const showOnlyFavorites = ref(store.showOnlyFavorites);
const showOnlyActive = ref(store.showOnlyActive);
const showOnlyInTimeFrame = ref(store.showOnlyInTimeFrame);
const modeRandom = ref(store.modeRandom);
const modeChronological = ref(store.modeChronological);
const modeReverseChronological = ref(store.modeReverseChronological);
const startDate = ref(store.startDate);
const endDate = ref(store.endDate);

function startDateUpdate(value) {
  console.log("Start Date Update:", value);
}

function endDateUpdate(value) {
  console.log("End Date Update:", value);
}

const saveConfig = () => {
  store.setInterval(interval.value);
  store.setShowOverlay(overlay.value);
  store.setShowNewChip(newchip.value);
  store.setShowOnlyFavorites(showOnlyFavorites.value);
  store.setShowOnlyActive(showOnlyActive.value);
  store.setModeRandom(modeRandom.value);
  store.setModeChronological(modeChronological.value);
  store.setModeReverseChronological(modeReverseChronological.value);

  if (showOnlyInTimeFrame.value && startDate.value == null) {
    showOnlyInTimeFrame.value = false;
  }
  if (startDate.value != null && endDate.value != null) {
    if (startDate.value > endDate.value) {
      endDate.value = null;
    }
  }

  store.setShowOnlyInTimeFrame(showOnlyInTimeFrame.value);
  store.setStartDate(startDate.value);
  store.setEndDate(endDate.value);
  //store.setStartDate(null);
  //store.setEndDate(null);
};

const goBack = () => {
  router.push("/");
};

function checkMode() {
  if (
    !modeRandom.value &&
    !modeChronological.value &&
    !modeReverseChronological.value
  ) {
    modeRandom.value = true;
  }
}

onMounted(() => {});

onUnmounted(() => {});
</script>

<style scoped>
.subtitle {
  font-size: 1rem;
  font-weight: 500;
  margin-top: 1rem;
}
</style>
