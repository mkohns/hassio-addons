<template>
  <v-container class="d-flex align-center justify-center" style="height: 100vh">
    <v-card width="100%" max-height="95%">
      <div class="d-flex">
        <v-card-title
          ><v-icon color="primary" class="mr-2">mdi-cog</v-icon>Slideshow
          Manage</v-card-title
        >
        <v-spacer></v-spacer>
        <v-btn
          variant="outlined"
          @click="timefilter = true"
          class="mt-3 mr-4"
          color="primary"
          >Time Filter</v-btn
        >
      </div>
      <v-card-text class="mt-5">
        <v-container style="height: 70vh; overflow-y: auto">
          <v-data-table
            :headers="headers"
            :items="filteredslides"
            item-value="AttachmentID"
            hide-default-footer
            class="elevation-1"
          >
            <template #item.TumbnailURL="{ item }">
              <v-img :src="item.TumbnailURL" width="100" height="100"></v-img>
            </template>
            <template #item.CreatedAt="{ item }">
              {{ new Date(item.CreatedAt).toLocaleString() }}
            </template>
            <template #item.Enabled="{ item }">
              <v-icon @click="disableImage(item)" v-if="item.Enabled"
                >mdi-eye</v-icon
              >
              <v-icon @click="enableImage(item)" v-else>mdi-eye-off</v-icon>
            </template>
            <template #item.Favorite="{ item }">
              <v-icon
                @click="unlikeImage(item)"
                v-if="item.Favorite"
                color="red"
                >mdi-heart</v-icon
              >
              <v-icon @click="likeImage(item)" v-else>mdi-heart-outline</v-icon>
            </template>
            <template #item.Delete="{ item }">
              <v-icon @click="deleteImage(item)">mdi-trash-can</v-icon>
            </template>
          </v-data-table>
        </v-container>
        <div class="d-flex">
          <v-spacer></v-spacer>
          <v-btn class="ml-4 mt-8" color="primary" @click="goBack">Back</v-btn>
        </div>
      </v-card-text>
    </v-card>
    <v-dialog v-model="timefilter" width="unset">
      <v-card>
        <v-date-picker
          v-model="timefilterSelect"
          color="primary"
          range
        ></v-date-picker>
        <div class="d-flex justify-center">
          <v-btn color="primary" text @click="timefilter = false">Close</v-btn>
          <v-btn class="mx-3 mb-3" color="primary" text @click="filterReset"
            >reset</v-btn
          >
          <v-btn color="primary" text @click="filterAction">Apply</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from "vue";
import router from "@/router";
import backend from "@/api/backend";

const timefilter = ref(false);
const timefilterSelect = ref(null);
const slides = ref([]);
const filteredslides = ref([]);
const headers = [
  { title: "Image", value: "TumbnailURL", sortable: true },
  { title: "Send By", value: "CreatedBy", sortable: true },
  { title: "Send At", value: "CreatedAt", sortable: true },
  { title: "Active", value: "Enabled", sortable: true },
  { title: "Favorite", value: "Favorite", sortable: true },
  { title: "Delete", value: "Delete", sortable: true },
];

const goBack = () => {
  router.push("/");
};

let url = import.meta.env.VITE_BASE_URL;
if (url === undefined) {
  url = "";
}

function enableImage(item) {
  console.log("enable");
  backend
    .resumeImage(item.AttachmentID)
    .then(() => {
      loadData();
    })
    .catch((err) => {
      console.log(err);
    });
}

function disableImage(item) {
  console.log("disable");
  backend
    .pauseImage(item.AttachmentID)
    .then(() => {
      loadData();
    })
    .catch((err) => {
      console.log(err);
    });
}

function likeImage(item) {
  console.log("like");
  backend
    .like(item.AttachmentID)
    .then(() => {
      loadData();
    })
    .catch((err) => {
      console.log(err);
    });
}

function unlikeImage(item) {
  console.log("like");
  backend
    .unlike(item.AttachmentID)
    .then(() => {
      loadData();
    })
    .catch((err) => {
      console.log(err);
    });
}

function deleteImage(item) {
  console.log("delete");
  backend
    .delete(item.AttachmentID)
    .then(() => {
      loadData();
    })
    .catch((err) => {
      console.log(err);
    });
}

const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      const index = entry.target.getAttribute("data-index");
      slides.value[index].isVisible = true;
      observer.unobserve(entry.target);
    }
  });
});

function filterReset() {
  console.log(timefilterSelect.value);
  timefilter.value = false;
  timefilterSelect.value = null;
  filteredslides.value = slides.value;
}

function filterAction() {
  console.log(timefilterSelect.value);
  const nextDay = new Date(timefilterSelect.value);
  nextDay.setDate(nextDay.getDate() + 1);
  console.log(nextDay);
  if (timefilterSelect.value === null) {
    return;
  }

  filteredslides.value = slides.value.filter((slide) => {
    return (
      new Date(slide.CreatedAt) >= new Date(timefilterSelect.value) &&
      new Date(slide.CreatedAt) <= nextDay
    );
  });
}

function loadData() {
  fetch(url + "/slides")
    .then((res) => res.json())
    .then((data) => {
      slides.value = data;
      if (url !== "") {
        slides.value.forEach((slide) => {
          slide.TumbnailURL = `${url}/${slide.TumbnailURL}`;
          //console.log(slide.TumbnailURL);
        });
      }
      filteredslides.value = slides.value;
    });
}

onMounted(() => {
  loadData();
  console.log("Manage Page");
});

onUnmounted(() => {});
</script>

<style scoped></style>
