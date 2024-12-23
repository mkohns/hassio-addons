<template>
  <v-container
    @click="goBack"
    class="d-flex align-center justify-center"
    style="height: 100vh"
  >
    <v-card width="50%">
      <v-card-title
        ><v-icon color="primary" class="mr-2">mdi-information</v-icon
        >Information</v-card-title
      >
      <v-card-text v-if="info" class="mt-5 mb-5">
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2"
                >mdi-image-multiple-outline</v-icon
              >Total pictures
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">{{ info.slidesCount }}</div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-heart-outline</v-icon>Total
              favorite pictures
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">{{ info.favoriteCount }}</div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-eye-check-outline</v-icon
              >Total active pictures
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">{{ info.activeCount }}</div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-database</v-icon>Total
              picture size
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">
              {{ formatBytes(info.slidesSize) }}
            </div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-database</v-icon>Total
              thumbnail size
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">
              {{ formatBytes(info.thumbnailSize) }}
            </div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2"
                >mdi-alpha-v-circle-outline</v-icon
              >Version
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">{{ info.version }}</div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col cols="6">
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-git</v-icon>Git Commit
            </div>
          </v-col>
          <v-col cols="6">
            <div class="list-item-value">{{ info.gitCommit }}</div>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col>
            <div class="list-item">
              <v-icon size="large" class="mr-2">mdi-ip-network-outline</v-icon
              >Your IP address
            </div>
          </v-col>
          <v-col>
            <div class="list-item-value">{{ info.remoteIP }}</div>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-text v-if="error" class="mt-5"> {{ error }} </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import router from "@/router";
import backend from "@/api/backend";

onMounted(() => {
  console.log("Mounted");
  loadData();
});

const info = ref(null);
const error = ref(null);

function loadData() {
  console.log("Loading data");
  backend
    .getInfo()
    .then((response) => {
      console.log(response);
      info.value = response.data;
    })
    .catch((err) => {
      console.log(err);
      error.value = err.response.data || err.message || "Unknown error";
    });
}

function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return (
    parseFloat((bytes / Math.pow(k, i)).toFixed(decimals)) + " " + sizes[i]
  );
}

const goBack = () => {
  router.push("/");
};
</script>

<style scoped>
.list-item {
  font-size: 1.2em;
  padding: 5px;
}
.list-item-value {
  background-color: rgb(70, 70, 70);
  font-size: 1.2em;
  padding: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
