import axios from "axios";

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

  pauseImage(imageId) {
    let body = {
      Enabled: false,
    };
    return this.instance.patch(`/slides/${imageId}`, body);
  }
  resumeImage(imageId) {
    let body = {
      Enabled: true,
    };
    return this.instance.patch(`/slides/${imageId}`, body);
  }
  like(imageId) {
    let body = {
      Favorite: true,
    };
    return this.instance.patch(`/slides/${imageId}`, body);
  }
  unlike(imageId) {
    let body = {
      Favorite: false,
    };
    return this.instance.patch(`/slides/${imageId}`, body);
  }
  delete(imageId) {
    return this.instance.delete(`/slides/${imageId}`);
  }
}

const backend = new BackendClient();

export default backend;