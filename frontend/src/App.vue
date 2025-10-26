<template>
  <v-app :theme="state.theme">
    <router-view />
    <Notify />
  </v-app>
</template>

<script lang="ts" setup>
import { useAppStore } from "./stores/app";
import Notify from "./components/Notify.vue";
const state = useAppStore();
</script>
