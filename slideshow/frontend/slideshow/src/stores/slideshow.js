import { defineStore } from "pinia";
import { set } from "vue-demi";
import { en } from "vuetify/locale";

export const useSlideshowStore = defineStore("slideshow", {
  state: () => ({
    interval: 5,
    showOverlay: true,
    showNewChip: false,
    showOnlyFavorites: false,
    showOnlyActive: false,
    prioNewImages: false,
    showOnlyInTimeFrame: false,
    modeRandom: true,
    modeChronological: false,
    modeReverseChronological: false,
    startDate: null,
    endDate: null,
  }),
  persist: true, // Enable persistence
  actions: {
    setInterval(newInterval) {
      this.interval = newInterval;
    },
    setShowOverlay(value) {
      this.showOverlay = value;
    },
    setShowNewChip(value) {
      this.showNewChip = value;
    },
    setShowOnlyFavorites(value) {
      this.showOnlyFavorites = value;
    },
    setShowOnlyActive(value) {
      this.showOnlyActive = value;
    },
    setShowOnlyInTimeFrame(value) {
      this.showOnlyInTimeFrame = value;
    },
    setModeRandom(value) {
      this.modeRandom = value;
    },
    setModeChronological(value) {
      this.modeChronological = value;
    },
    setModeReverseChronological(value) {
      this.modeReverseChronological = value;
    },
    setStartDate(value) {
      this.startDate = value;
    },
    setEndDate(value) {
      this.endDate = value;
    },
    setPrioNewImages(value) {
      this.prioNewImages = value;
    },
  },
});
