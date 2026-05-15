<template>
  <page-content title="Docker Management">
    <template #append>
      <v-btn
        variant="flat"
        color="primary"
        icon="mdi-refresh"
        size="small"
        @click="refreshData"
      ></v-btn>
    </template>
    <v-card>
      <v-tabs color="primary" v-model="tab">
        <v-tab value="container">Container</v-tab>
        <v-tab value="images">Images</v-tab>
        <v-tab value="network">Network</v-tab>
        <v-tab value="volume">Volume</v-tab>
      </v-tabs>

      <v-divider></v-divider>

      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="container">
          <v-lazy>
            <ContainerList />
          </v-lazy>
        </v-tabs-window-item>
        <v-tabs-window-item value="images">
          <v-lazy>
            <ImageList />
          </v-lazy>
        </v-tabs-window-item>
        <v-tabs-window-item value="network">
          <v-lazy>
            <NetworkList />
          </v-lazy>
        </v-tabs-window-item>
        <v-tabs-window-item value="volume">
          <v-lazy>
            <VolumeList />
          </v-lazy>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </page-content>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";

const tab = ref("container");
const state = useDockerStore();
async function refreshData() {
  await state.fetchAll();
}
</script>
