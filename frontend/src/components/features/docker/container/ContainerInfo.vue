<template>
  <div v-if="!state.loading">
    <section id="info-section" class="mb-5">
      <div class="w-full border rounded-lg d-flex flex-column ga-2 pa-3">
        <div class="d-flex justify-space-between">
          <p>ID</p>
          <p>{{ state.container.container_info.id }}</p>
        </div>
        <hr />
        <div class="d-flex justify-space-between">
          <p>Name</p>
          <p>{{ state.container.container_info.name }}</p>
        </div>
        <hr />
        <div class="d-flex justify-space-between">
          <p>Image</p>
          <p>{{ state.container.container_info.image }}</p>
        </div>
        <hr />
        <div class="d-flex justify-space-between">
          <p>Status</p>
          <p>{{ state.container.container_info.status }}</p>
        </div>
        <hr />
        <div class="d-flex justify-space-between">
          <p>State</p>
          <v-chip
            :color="containerStateColor(state.container.container_info.state)"
            >{{ state.container.container_info.state }}</v-chip
          >
        </div>
      </div>
    </section>
    <section id="info-section" class="mb-5">
      <h5>Port Forwards</h5>
      <v-table density="compact" class="border rounded-lg">
        <thead>
          <tr>
            <th class="text-left">Host Port</th>
            <th class="text-left">Container Port</th>
            <th class="text-left">Protocol</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in state.container.ports" :key="index">
            <td>{{ item.PublicPort }}</td>
            <td>{{ item.PrivatePort }}</td>
            <td>{{ item.Type.toUpperCase() }}</td>
          </tr>
        </tbody>
      </v-table>
    </section>
    <section id="info-section" class="mb-5">
      <h5>Mounts</h5>
      <v-table density="compact" class="border rounded-lg">
        <thead>
          <tr>
            <th class="text-left">Source</th>
            <th class="text-left">Destination</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in state.container.mounts" :key="index">
            <td>{{ item.Source }}</td>
            <td>{{ item.Destination }}</td>
          </tr>
        </tbody>
      </v-table>
    </section>
    <section id="label-section">
      <h5>Labels</h5>
      <v-table density="compact" class="border rounded-lg">
        <thead>
          <tr>
            <th class="text-left">Key</th>
            <th class="text-left">Value</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="([key, value], index) in Object.entries(
              state.container.container_info.labels
            )"
            :key="index"
          >
            <td>{{ key }}</td>
            <td>{{ value }}</td>
          </tr>
        </tbody>
      </v-table>
    </section>
  </div>
</template>
<script setup lang="ts">
import { useDockerStore } from "@/stores/docker";
import { containerStateColor } from "@/utils";

const state = useDockerStore();
const route = useRoute();
const params = route.params as { id: string };
onMounted(async () => {
  console.log(params);
  await state.fetchContainerInfo(params.id);
});
</script>
