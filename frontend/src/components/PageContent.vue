<template>
  <v-toolbar
    density="comfortable"
    :color="state.theme == 'light' ? 'white' : '#2C303A'"
    :title="title"
  >
    <template v-slot:append>
      <slot name="append" />
    </template>
  </v-toolbar>
  <br />
  <v-sheet width="100%">
    <slot />
  </v-sheet>
</template>
<script setup lang="ts">
import { useAppStore } from "@/stores/app";

defineProps<{
  title: string; // Optional string prop
}>();
const state = useAppStore();
</script>
