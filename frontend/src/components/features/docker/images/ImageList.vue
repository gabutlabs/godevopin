<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.images" :headers="imageHeaders">
      <template #item.tags="{ item }">
        <template v-if="(item as any).tags.length > 0">
          <v-chip v-for="tag in (item as any).tags">
            {{ tag }}
          </v-chip>
        </template>
      </template>
    </v-data-table>
  </v-sheet>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";

const state = useDockerStore();

const imageHeaders = [
  { title: "ID", value: "id" },
  { title: "Tags", value: "tags" },
  { title: "Size Megabyte(MB)", value: "size_mb" },
  { title: "Created", value: "created" },
];

onMounted(async () => {
  await state.fetchAllImages();
});
</script>
