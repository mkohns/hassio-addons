<template>
  <v-container class="d-flex align-center justify-center" style="height: 100vh">
    <v-card width="100%">
      <v-card-text>
        <v-row>
          <v-col>
            <div class="subtitle">
              <v-icon class="mr-2" color="primary">mdi-cog</v-icon>General
            </div>
          </v-col>
          <v-col>
            <div class="subtitle">
              <v-icon class="mr-2" color="primary">mdi-filter</v-icon>Filter
            </div>
          </v-col>
          <v-col>
            <div class="subtitle">
              <v-icon class="mr-2" color="primary">mdi-play-box-outline</v-icon
              >Slide Modes
            </div>
          </v-col>
        </v-row>
        <v-row class="mt-7">
          <v-col>
            <div class="d-flex">
              <v-text-field
                class="seconds"
                v-model="interval"
                density="compact"
                label="Seconds"
                readonly
                @click="openIntervalDialog"
                variant="outlined"
                hide-details
              ></v-text-field>
              <v-btn
                class="icon-size ml-5"
                @click="timeUp"
                icon="mdi-plus-circle"
              >
              </v-btn>
              <v-btn
                class="icon-size ml-1"
                @click="timeDown"
                icon="mdi-minus-circle"
              >
              </v-btn>
            </div>

            <v-switch
              hide-details
              color="primary"
              :label="portraitMode ? 'Portrait Mode' : 'Landscape Mode'"
              v-model="portraitMode"
            >
            </v-switch>
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
          </v-col>
          <v-col>
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
            <v-switch
              hide-details
              color="primary"
              label="Prioritize New Images"
              v-model="prioNewImages"
            ></v-switch>
            <v-switch
              hide-details
              color="primary"
              label="Show Only Images in timeframe"
              v-model="showOnlyInTimeFrame"
            >
            </v-switch>
            <DateSelect
              v-if="showOnlyInTimeFrame"
              label="Start Date"
              v-model="startDate"
              @update:modelValue="startDateUpdate"
            ></DateSelect>
            <DateSelect
              nullValue="Today"
              v-if="showOnlyInTimeFrame"
              :disabled="!showOnlyInTimeFrame"
              label="End Date"
              v-model="endDate"
              @update:modelValue="endDateUpdate"
            ></DateSelect>
          </v-col>
          <v-col>
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
          </v-col>
        </v-row>

        <div class="d-flex mt-4">
          <v-spacer></v-spacer>
          <v-btn color="primary" @click="saveConfig">Save</v-btn>
          <v-btn class="ml-3" color="primary" @click="goBack">Close</v-btn>
        </div>
      </v-card-text>
    </v-card>
    <v-dialog v-model="intervalDialog" width="70%">
      <v-card>
        <div class="dialog-seconds">{{ intervalTmp }} seconds</div>
        <v-slider
          v-model="intervalTmp"
          min="5"
          max="500"
          step="10"
          hide-details="true"
          color="primary"
          class="px-10 pb-10"
        >
        </v-slider>
        <div class="d-flex justify-center mb-6">
          <v-btn
            class="mr-5"
            color="primary"
            text
            @click="intervalDialog = false"
            >Close</v-btn
          >
          <v-btn color="primary" text @click="setInterval">Apply</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useSlideshowStore } from "@/stores/slideshow";
import router from "@/router";
import DateSelect from "@/components/dateSelect.vue";

const store = useSlideshowStore();

const interval = ref(store.interval);
const intervalTmp = ref(null);
const overlay = ref(store.showOverlay);
const newchip = ref(store.showNewChip);
const showOnlyFavorites = ref(store.showOnlyFavorites);
const showOnlyActive = ref(store.showOnlyActive);
const showOnlyInTimeFrame = ref(store.showOnlyInTimeFrame);
const prioNewImages = ref(store.prioNewImages);
const modeRandom = ref(store.modeRandom);
const modeChronological = ref(store.modeChronological);
const modeReverseChronological = ref(store.modeReverseChronological);
const startDate = ref(store.startDate);
const endDate = ref(store.endDate);
const portraitMode = ref(store.portraitMode);

const intervalDialog = ref(false);

function setInterval() {
  interval.value = intervalTmp.value;
  intervalDialog.value = false;
}

function openIntervalDialog() {
  intervalTmp.value = interval.value;
  intervalDialog.value = true;
}

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
  store.setPrioNewImages(prioNewImages.value);
  store.setModeRandom(modeRandom.value);
  store.setModeChronological(modeChronological.value);
  store.setModeReverseChronological(modeReverseChronological.value);
  store.setPortraitMode(portraitMode.value);

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

  goBack();
};

const goBack = () => {
  router.push("/");
};

function timeUp() {
  interval.value++;
}

function timeDown() {
  if (interval.value > 5) {
    interval.value--;
  }
}

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
  font-size: 1.3rem;
  font-weight: 500;
  margin-top: 1rem;
}
.seconds {
  max-width: 100px;
}
.icon-size {
  font-size: 1.8rem;
}
.dialog-seconds {
  text-align: center;
  font-size: 3rem;
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
