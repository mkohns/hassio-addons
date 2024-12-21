<template>
  <div @click="onClick" :class="classes">
    <v-btn size="x-large" :icon="icon"></v-btn
    ><span class="menu-text ml-4">{{ text }}</span>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";

function onClick(evt) {
  evt.stopPropagation();
  emit("trigger", props.event);
}

const classes = computed(() => {
  if (props.show === undefined) {
    return {
      "menu-item": true,
      "menu-hide": true,
    };
  }
  return {
    "menu-item": true,
    animate__animated: true,
    animate__slideInLeft: props.show ? true : false,
    animate__slideOutLeft: props.show ? false : true,
  };
});

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
  show: {
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
  background-color: rgba(255, 255, 255, 0.346);
  padding: 10px;
}
.menu-hide {
  transform: translateX(-100%);
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
