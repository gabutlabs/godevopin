<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.networks" :headers="headers"> </v-data-table>
  </v-sheet>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";

const state = useDockerStore();
const headers = [
  { title: "ID", value: "id" },
  { title: "Name", value: "name" },
  { title: "Driver", value: "driver" },
  { title: "Scope", value: "scope" },
];

onMounted(async () => {
  await state.fetchAllNetworks();
});
</script>
