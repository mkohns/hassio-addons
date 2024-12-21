<template>
  <v-text-field
    v-model="datestring"
    readonly
    :disabled="disabled"
    :label="label"
    @click="showDialog"
    append-inner-icon="mdi-close"
    @click:append-inner="clearDate"
  ></v-text-field>
  <v-dialog v-model="dialog" persistent max-width="390px">
    <v-card>
      <v-card-title class="headline">Select a {{ label }}</v-card-title>
      <v-card-text>
        <v-date-picker v-model="date" @input="save"></v-date-picker>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="dialog = false">Cancel</v-btn>
        <v-btn :disabled="!date" color="blue darken-1" text @click="save"
          >OK</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, computed, watch } from "vue";

// Define emits
const emit = defineEmits(["update:modelValue"]);

const props = defineProps({
  modelValue: {
    type: [String, null],
    required: true,
  },
  label: {
    type: String,
    default: "Date",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  nullValue: {
    type: String,
  },
});

function clearDate(evt) {
  evt.stopPropagation();
  if (props.nullValue) {
    datestring.value = props.nullValue;
  } else {
    datestring.value = null;
  }
  emit("update:modelValue", null);
}

function showDialog() {
  date.value = null;
  dialog.value = true;
}

function save() {
  console.log("Save date");
  datestring.value = date.value.toLocaleDateString();
  emit("update:modelValue", datestring.value);
  dialog.value = false;
}

const dialog = ref(false);
const datestring = ref(null);
const date = ref(null);

// Watch for changes in modelValue
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue === null) {
      datestring.value = props.nullValue;
    } else {
      datestring.value = newValue;
    }
  },
  { immediate: true }
);
</script>
