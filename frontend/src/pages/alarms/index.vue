<template>
  <page-content title="Alarm">
    <v-tabs color="primary" bg-color="white" hide-slider v-model="tab">
      <v-tab value="active">Active</v-tab>
      <v-tab value="history">History</v-tab>
    </v-tabs>

    <v-divider></v-divider>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="active">
        <v-sheet class="pa-5">
          <v-data-table :items="state.activeAlarms" :headers="activeHeaders">
            <template #item.actions="{ item }">
              <ActionTable :action-items="activeActionMenuItems" :item="item" />
            </template>
            <template #item.started_at="{ item }">
              {{ new Date(item.started_at).toLocaleString() }}
            </template>
            <template #item.last_notified_at="{ item }">
              {{ new Date(item.last_notified_at).toLocaleString() }}
            </template>
          </v-data-table>
        </v-sheet>
      </v-tabs-window-item>
      <v-tabs-window-item value="history">
        <v-sheet class="pa-5">
          <v-data-table :items="state.historyAlarms.data"></v-data-table>
        </v-sheet>
      </v-tabs-window-item>
    </v-tabs-window>
  </page-content>
</template>
<script lang="ts" setup>
import PageContent from "@/components/PageContent.vue";
import ActionTable from "@/components/ActionTable.vue";
import { useAlarmStore, type ActiveAlarm } from "@/stores/alarm";
const tab = ref("active");
const state = useAlarmStore();
const page = ref(1);
const limit = ref(10);
const activeHeaders = [
  { title: "Alarm Name", value: "alarm_name" },
  { title: "Target", value: "target" },
  { title: "Started At", value: "started_at" },
  { title: "Last Notified At", value: "last_notified_at" },
  { title: "Message", value: "message" },
  { title: "Actions", value: "actions", sortable: false },
];
const activeActionMenuItems = [
  {
    title: "Edit",
    icon: "mdi-pencil",
    onClick: (item: ActiveAlarm) => {},
  },
  {
    title: "Delete",
    style: "color: red",
    icon: "mdi-delete",
    onClick: async (item: ActiveAlarm) => {},
  },
];
onMounted(async () => {
  await state.fetchActiveAlarms();
  await state.fetchHistoryAlarms(page.value, limit.value);
});
</script>
