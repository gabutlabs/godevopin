<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.volumes" :headers="headers"> </v-data-table>
  </v-sheet>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";

const state = useDockerStore();
const headers = [
  { title: "Name", value: "name" },
  { title: "Driver", value: "driver" },
  { title: "Mount Point", value: "mountpoint" },
];

onMounted(async () => {
  await state.fetchAllVolumes();
});
</script>
