<template>
  <v-app :theme="state.theme">
    <router-view />
    <Notivue v-slot="item">
      <Notification :item="item" />
    </Notivue>
  </v-app>
</template>

<script lang="ts" setup>
import { useAppStore } from "./stores/app";
import { Notivue, Notification } from "notivue";
const state = useAppStore();
</script>
