<template>
  <transition name="slide-x">
    <div v-if="props.open" class="menu-container">
      <SideMenuItem
        v-for="item in itemsInternal"
        :key="item.text"
        :icon="item.icon"
        :text="item.text"
        :event="item.event"
        class="item"
        @trigger="eventTrigger"
      ></SideMenuItem>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";
import SideMenuItem from "./sideMenuItem.vue";
import { VSlideXTransition } from "vuetify/components";

onMounted(() => {
  itemsInternal.value = props.items.map((item) => {
    return {
      ...item,
    };
  });
});

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

onMounted(() => {});

// Watch the open prop
watch(
  () => props.open,
  (newVal, oldVal) => {
    if (newVal) {
      console.log("Menu opened");
    } else {
      console.log("Menu closed");
    }
  }
);

watch(
  () => props.items,
  (newVal, oldVal) => {
    console.log("Items changed");
    for (let i = 0; i < props.items.length; i++) {
      itemsInternal.value[i].icon = props.items[i].icon;
      itemsInternal.value[i].text = props.items[i].text;
    }
  },
  { deep: true }
);
</script>

<style scoped>
.menu-container {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.item:first-child {
  border-radius: 0 50px 0 0;
}
.item:last-child {
  border-radius: 0 0 50px;
}
/* Custom slide transition */
.slide-x-enter-active,
.slide-x-leave-active {
  transition: transform 1s ease;
}
.slide-x-enter-from,
.slide-x-leave-to {
  transform: translateX(-100%);
}
</style>
