<template>
  <div class="menu-container">
    <SideMenuItem
      v-for="(item, index) in itemsInternal"
      :icon="item.icon"
      :text="item.text"
      :show="item.show"
      :event="item.event"
      class="item"
      @trigger="eventTrigger"
    ></SideMenuItem>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";
import SideMenuItem from "./sideMenuItem.vue";

// Define emits
const emit = defineEmits(["trigger"]);

// Define props
const props = defineProps({
  open: {
    type: Boolean,
    required: true,
  },
  items: {
    type: Array,
    required: true,
  },
});

function eventTrigger(evt) {
  emit("trigger", evt);
}

const itemsInternal = ref([]);

onMounted(() => {
  itemsInternal.value = props.items.map((item) => {
    return {
      ...item,
      show: false,
    };
  });
});

// Watch the open prop
watch(
  () => props.open,
  (newVal, oldVal) => {
    if (newVal) {
      console.log("Menu opened");
      itemsInternal.value.forEach((item, index) => {
        setTimeout(() => {
          item.show = true;
        }, index * 100);
      });

      // Perform actions when the menu is opened
    } else {
      console.log("Menu closed");
      // Perform actions when the menu is closed
      itemsInternal.value.forEach((item, index) => {
        item.show = false;
      });
    }
  }
);
</script>

<style scoped>
.menu-container {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 0; /* Adjust as needed */
}
.item:first-child {
  border-radius: 0 50px 0 0;
}

.item:last-child {
  border-radius: 0 0 50px;
}
</style>
