import { defineStore } from "pinia";

export const useSlideshowStore = defineStore("slideshow", {
  state: () => ({
    interval: 5,
    showOverlay: true,
  }),
  persist: true, // Enable persistence
  actions: {
    setInterval(newInterval) {
      this.interval = newInterval;
    },
    setShowOverlay(value) {
      this.showOverlay = value;
    },
  },
});
