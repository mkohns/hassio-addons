import { defineStore } from "pinia";
import { set } from "vue-demi";

export const useSlideshowStore = defineStore("slideshow", {
  state: () => ({
    interval: 5,
    showOverlay: true,
    showNewChip: false,
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
  },
});
