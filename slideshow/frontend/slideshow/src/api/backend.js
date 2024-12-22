import axios from "axios";
import { useSlideshowStore } from "@/stores/slideshow";

class BackendClient {
  instance = axios.create({
    baseURL: import.meta.env.VITE_BASE_URL,
    timeout: 10000,
  });

  constructor() {
    console.log("BackendClient constructor");
  }

  handleError(error) {
    let err;
    console.log(error);
    if (error.response && error.response.data) {
      err = error.response.data;
      if (err.message === "Error validating claims. Not in app user group") {
        err.message =
          "You do not have all needed permissions to use the Consumer Application Self Service. Please contact api-mgmt-support@schaeffler.com to get needed group membership.";
      }
    } else {
      err = {
        message: "network error occurred. please retry.",
        requestId: "not available",
      };
    }
    // This feature enables slot buttons to do something which needs more time
    err.loading = false;
    return err;
  }

  pauseImage(attachmendId) {
    let body = {
      Enabled: false,
    };
    return this.instance.patch(`/slides/${attachmendId}`, body);
  }
  resumeImage(attachmendId) {
    let body = {
      Enabled: true,
    };
    return this.instance.patch(`/slides/${attachmendId}`, body);
  }
  like(attachmendId) {
    let body = {
      Favorite: true,
    };
    return this.instance.patch(`/slides/${attachmendId}`, body);
  }
  unlike(attachmendId) {
    let body = {
      Favorite: false,
    };
    return this.instance.patch(`/slides/${attachmendId}`, body);
  }
  delete(attachmendId) {
    return this.instance.delete(`/slides/${attachmendId}`);
  }
  getInfo() {
    return this.instance.get(`/info`);
  }
  updateOrientation(portraitMode) {
    let body = {
      PortraitMode: portraitMode,
    };
    return this.instance.patch(`/portraitMode`, body);
  }
  nextSlide() {
    const store = useSlideshowStore();

    let params = {
      showOnlyFavorites: store.showOnlyFavorites,
      showOnlyActive: store.showOnlyActive,
      showOnlyInTimeFrame: store.showOnlyInTimeFrame,
      showNewImagesWithPriority: store.prioNewImages,
      portraitMode: store.portraitMode,
      modeRandom: store.modeRandom,
      modeChronological: store.modeChronological,
      modeReverseChronological: store.modeReverseChronological,
      startDate: store.startDate,
      endDate: store.endDate,
    };

    return this.instance.get(`/nextslide`, { params: params });
  }
}

const backend = new BackendClient();

export default backend;
