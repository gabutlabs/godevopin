<template>
  <v-sheet class="pa-5">
    <v-data-table :items="state.activeAlarms" :headers="activeHeaders">
      <template #item.actions="{ item }">
        <ActionTable :action-items="activeActionMenuItems" :item="item" />
      </template>
      <template #item.status="{ item }">
        <v-chip
          :color="item.status == 'FIRING' ? 'error' : 'success'"
          size="small"
          >{{ item.status }}</v-chip
        >
      </template>
      <template #item.started_at="{ item }">
        {{ new Date(item.started_at).toLocaleString() }}
      </template>
      <template #item.last_notified_at="{ item }">
        {{ new Date(item.last_notified_at).toLocaleString() }}
      </template>
    </v-data-table>
  </v-sheet>
</template>
<script lang="ts" setup>
import ActionTable from "@/components/ActionTable.vue";
import { useAlarmStore, type ActiveAlarm } from "@/stores/alarm";
const state = useAlarmStore();
const activeHeaders = [
  { title: "Alarm Name", value: "alarm_name" },
  { title: "Target", value: "target" },
  { title: "Status", value: "status" },
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
];
onMounted(async () => {
  await state.fetchActiveAlarms();
});
</script>
