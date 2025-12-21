<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.containers" :headers="headers">
      <template #item.state="{ item }">
        <v-chip :color="stateColor(item.state)">{{ item.state }}</v-chip>
      </template>
    </v-data-table>
  </v-sheet>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";
import type { ContainerState } from "@/types/docker.type";

const state = useDockerStore();
const headers = [
  { title: "ID", value: "id" },
  { title: "Name", value: "name" },
  { title: "Image", value: "image" },
  { title: "status", value: "status" },
  { title: "State", value: "state" },
];

function stateColor(state: ContainerState): string {
  let colorType = "";
  switch (true) {
    case state == "running":
      colorType = "green";
      break;
    case state == "exited":
      colorType = "red";
    case state == "dead":
      colorType = "red";
      break;
    case state == "paused":
      colorType = "yellow";
      break;
    default:
      colorType = "primary";
      break;
  }
  return colorType;
}

onMounted(async () => {
  await state.fetchAllContainers();
});
</script>
