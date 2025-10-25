<template>
  <page-content title="Alarm">
    <v-tabs color="primary" hide-slider v-model="tab">
      <v-tab value="active">Active</v-tab>
      <v-tab value="history">History</v-tab>
    </v-tabs>

    <v-divider></v-divider>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="active">
        <ActiveAlarmComponent />
      </v-tabs-window-item>
      <v-tabs-window-item value="history">
        <v-sheet class="pa-5">
          <HistoryAlarm />
        </v-sheet>
      </v-tabs-window-item>
    </v-tabs-window>
  </page-content>
</template>
<script lang="ts" setup>
import PageContent from "@/components/PageContent.vue";
import ActiveAlarmComponent from "@/components/features/alarm/ActiveAlarm.vue";
const tab = ref("active");
</script>
