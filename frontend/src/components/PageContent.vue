<template>
  <div class="d-flex align-center justify-space-between mb-6">
    <h1 class="text-h4 font-weight-bold">{{ title }}</h1>
    <div>
      <slot name="append" />
    </div>
  </div>
  <v-sheet width="100%" class="bg-transparent">
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
