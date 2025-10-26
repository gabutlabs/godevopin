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
import { useNotify } from "@/composables/useNotify";
import { useAlarmStore, type ActiveAlarm } from "@/stores/alarm";
import { useAuthStore } from "@/stores/auth";
import type { User } from "@/stores/user";
const props = withDefaults(defineProps<{ filterStatus: string }>(), {
  filterStatus: "FIRING",
});
const state = useAlarmStore();
const authState = useAuthStore();
const notify = useNotify();
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
    title: "Acknowledge",
    icon: "mdi-check-circle-outline",
    onClick: async (item: ActiveAlarm) => {
      const user = authState.user! as any as User;
      await state.acknowledgeAlarm(
        { alarmName: item.alarm_name, target: item.target },
        { acknowledged_by: String(user.id) }
      );
      if (state.action_result.is_success) {
        await state.fetchActiveAlarmByStatus(props.filterStatus);
      } else {
        notify.error(
          `Failed to acknowledge alarm: ${
            state.action_result?.message ?? "Unknown error"
          }`
        );
      }
    },
  },
];
watch(
  () => props.filterStatus,
  async (newStatus) => {
    await state.fetchActiveAlarmByStatus(newStatus);
  },
  { immediate: true }
);
onMounted(async () => {
  await state.fetchActiveAlarmByStatus(props.filterStatus);
});
</script>
