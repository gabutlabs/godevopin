<template>
  <PageContent title="Container Info">
    <v-sheet elevation="4">
      <v-tabs color="primary" v-model="tab">
        <v-tab value="info">Info</v-tab>
        <v-tab value="logs">Logs</v-tab>
      </v-tabs>

      <v-divider></v-divider>

      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="info">
          <v-sheet class="pa-5 w-full"><ContainerInfo /></v-sheet>
        </v-tabs-window-item>
        <v-tabs-window-item value="logs">
          <v-sheet class="pa-5 w-full"><ContainerLogs /></v-sheet>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-sheet>
  </PageContent>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import ContainerInfo from '@/components/features/docker/container/ContainerInfo.vue';
import ContainerLogs from '@/components/features/docker/container/ContainerLogs.vue';

const tab = ref("info");
</script>
