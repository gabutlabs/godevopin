<template>
  <v-sheet class="pa-5">
    <v-data-table-server
      v-model:items-per-page="pagination.pageSize"
      :headers="historyHeaders"
      :items="state.historyAlarms.data"
      :items-length="state.historyAlarms.pagination?.total || 0"
      :loading="loading"
      item-value="id"
      @update:options="loadHistoryAlarms"
    >
      <template #item.created_at="{ item }">
        {{ new Date(item.created_at).toLocaleString() }}
      </template>
    </v-data-table-server>
  </v-sheet>
</template>
<script setup lang="ts">
import { useAlarmStore } from "@/stores/alarm";
const props = withDefaults(defineProps<{ filterStatus: string }>(), {
  filterStatus: "",
});
const historyHeaders = [
  { title: "Alarm Name", value: "alarm_name" },
  { title: "Target", value: "target" },
  { title: "Status", value: "status" },
  { title: "Message", value: "message" },
  { title: "Created At", value: "created_at" },
];
const loading = ref(false);
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0,
});
const state = useAlarmStore();
function loadHistoryAlarms({ page, itemsPerPage, sortBy }: any) {
  pagination.currentPage = page;
  pagination.pageSize = itemsPerPage;
  loading.value = true;
  state
    .fetchHistoryAlarms(
      pagination.currentPage,
      pagination.pageSize,
      props.filterStatus
    )
    .finally(() => {
      loading.value = false;
    });
}
watch(
  () => props.filterStatus,
  () => {
    loadHistoryAlarms({
      page: pagination.currentPage,
      itemsPerPage: pagination.pageSize,
    });
  }
);
onMounted(() => {
  loading.value = true;
  state
    .fetchHistoryAlarms(
      pagination.currentPage,
      pagination.pageSize,
      props.filterStatus
    )
    .finally(() => {
      loading.value = false;
    });
});
</script>
