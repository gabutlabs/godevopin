<template>
  <v-snackbar
    v-model="snackbar"
    :color="color"
    :timeout="timeout"
    location="top"
    multi-line
  >
    <div class="d-flex align-center">
      <v-icon v-if="icon" class="mr-2">{{ icon }}</v-icon>
      <span>{{ message }}</span>
    </div>
    <template v-slot:actions>
      <v-btn variant="text" @click="snackbar = false"> Tutup </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { useNotify } from "@/composables/useNotify";

const { snackbar, message, color, timeout, icon } = useNotify();
</script>
