<template>
  <div @click="onClick" class="menu-item">
    <v-btn color="primary" size="x-large" :icon="icon"></v-btn
    ><span class="menu-text ml-4">{{ text }}</span>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";

function onClick(evt) {
  evt.stopPropagation();
  emit("trigger", props.event);
}

// Define emits
const emit = defineEmits(["trigger"]);

// Define props
const props = defineProps({
  icon: {
    type: String,
    required: true,
  },
  text: {
    type: String,
    required: true,
  },
  event: {
    type: String,
    required: true,
  },
});
</script>

<style scoped>
.menu-item {
  background-color: rgb(0 0 0 / 73%);
  padding: 10px;
}
.menu-text {
  margin-left: 10px;
  font-size: 1.2rem;
  font-weight: 500;
  color: white;
}
:deep() .v-btn .mdi-heart {
  color: red;
}
</style>
