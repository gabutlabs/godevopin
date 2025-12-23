<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.containers" :headers="headers">
      <template #item.actions="{ item }">
        <action-table :action-items="actionMenuItems" :item="item" />
      </template>
      <template #item.state="{ item }">
        <v-chip :color="containerStateColor(item.state)">{{
          item.state
        }}</v-chip>
      </template>
    </v-data-table>
  </v-sheet>
</template>
<script setup lang="ts">
import type ActionTableVue from "@/components/ActionTable.vue";
import router from "@/router";
import { useDockerStore } from "@/stores/docker";
import type { ContainerInfo, ContainerState } from "@/types/docker.type";
import { containerStateColor } from "@/utils";

const state = useDockerStore();
const headers = [
  { title: "ID", value: "id" },
  { title: "Name", value: "name" },
  { title: "Image", value: "image" },
  { title: "status", value: "status" },
  { title: "State", value: "state" },
  { title: "Actions", value: "actions", sortable: false },
];
const actionMenuItems = [
  {
    title: "Detail",
    icon: "mdi-eye",
    onClick: (item: ContainerInfo) => {
      router.push(`/docker/container/${item.id}`);
    },
  },
  // {
  //   title: "Delete",
  //   style: "color: red",
  //   icon: "mdi-delete",
  //   onClick: async (item: User) => {
  //     await state.deleteUser(item.id);
  //     if (state.action_result.is_success) {
  //       await state.fetchAllUsers();
  //       notify.success(state.action_result.message);
  //     } else {
  //       notify.error(state.action_result.message);
  //     }
  //   },
  // },
];

onMounted(async () => {
  await state.fetchAllContainers();
});
</script>
