<template>
  <v-main class="d-flex align-center justify-center pa-4">
    <v-sheet
      width="100%"
      max-width="450"
      class="mx-auto"
      rounded
      :elevation="3"
    >
      <router-view />
    </v-sheet>
  </v-main>
</template>
