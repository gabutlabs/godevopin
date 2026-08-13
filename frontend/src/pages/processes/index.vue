<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <PageContent title="Processes">
    <template #append>
      <div class="d-flex ga-3 mx-4 align-center">
        <v-text-field
          v-model="search"
          label="Search process"
          prepend-inner-icon="mdi-magnify"
          density="compact"
          hide-details
          clearable
          variant="outlined"
          style="min-width: 220px"
        />
        <v-select
          v-model="refreshSeconds"
          label="Refresh"
          :items="refreshItems"
          density="compact"
          hide-details
          variant="outlined"
          style="width: 130px"
        />
        <v-btn
          :prepend-icon="isRefreshing ? 'mdi-pause' : 'mdi-play'"
          variant="outlined"
          @click="toggleRefresh"
        >
          {{ isRefreshing ? "Pause" : "Resume" }}
        </v-btn>
      </div>
    </template>

    <v-row>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">VISIBLE PROCESSES</div>
            <div class="text-h4">{{ processStore.processes.length }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">TOTAL CPU</div>
            <div class="text-h4">{{ formatPercent(totalCPU) }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">TOTAL RESIDENT MEMORY</div>
            <div class="text-h4">{{ formatBytes(totalMemory) }}</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12">
        <v-card variant="outlined">
          <v-data-table
            :headers="headers"
            :items="processStore.processes"
            :loading="processStore.loading"
            item-value="process_key"
            :items-per-page="25"
          >
            <template #item.pid="{ item }">
              <span class="font-weight-medium">{{ item.pid }}</span>
            </template>
            <template #item.name="{ item }">
              <div class="text-truncate" style="max-width: 280px">
                {{ item.name }}
              </div>
              <div class="text-caption text-medium-emphasis">
                {{ item.username || "unknown user" }}
              </div>
            </template>
            <template #item.status="{ item }">
              <v-chip size="small" :color="statusColor(item.status)">
                {{ item.status || "unknown" }}
              </v-chip>
            </template>
            <template #item.cpu_percent="{ item }">
              <span :class="usageClass(item.cpu_percent)">
                {{ formatPercent(item.cpu_percent) }}
              </span>
            </template>
            <template #item.memory_percent="{ item }">
              <div>{{ formatPercent(item.memory_percent) }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ formatBytes(item.memory_bytes) }}
              </div>
            </template>
            <template #item.thread_count="{ item }">
              {{ item.thread_count || "-" }}
            </template>
            <template #item.actions="{ item }">
              <v-tooltip text="View process history" location="top">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    icon="mdi-chart-line"
                    size="small"
                    variant="text"
                    :aria-label="`View history for ${item.name}`"
                    @click.stop.prevent="openHistory(normalizeProcessItem(item))"
                  />
                </template>
              </v-tooltip>
            </template>
          </v-data-table>
        </v-card>
      </v-col>

    </v-row>

    <v-dialog v-model="historyDialog" max-width="1200">
      <v-card v-if="selectedProcess">
        <v-card-title class="d-flex align-center">
          <span>History: {{ selectedProcess.name }} (PID {{ selectedProcess.pid }})</span>
          <v-spacer />
          <v-select
            v-model="historyFilter"
            :items="historyItems"
            density="compact"
            hide-details
            variant="outlined"
            style="width: 150px"
          />
          <v-btn icon="mdi-close" variant="text" @click="closeHistory" />
        </v-card-title>
        <v-card-subtitle class="pb-2">
          Process identity: {{ selectedProcess.process_key }}
        </v-card-subtitle>
        <v-card-text>
          <div v-if="processStore.loadingHistory" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" />
            <div class="mt-3">Loading process history...</div>
          </div>
          <template v-else-if="processStore.history.length > 0">
            <v-row>
              <v-col cols="12" lg="6">
                <WidgetLineChart
                  :data="processStore.history"
                  title="CPU Usage"
                  data-key="avg_cpu_percent"
                  series-name="CPU (%)"
                  unit="%"
                  :filter-range="historyFilter"
                />
              </v-col>
              <v-col cols="12" lg="6">
                <WidgetLineChart
                  :data="processStore.history"
                  title="Memory Usage"
                  data-key="avg_memory_percent"
                  series-name="Memory (%)"
                  unit="%"
                  :filter-range="historyFilter"
                />
              </v-col>
            </v-row>
          </template>
          <div v-else-if="processStore.historyError" class="text-center text-error py-8">
            {{ processStore.historyError }}
          </div>
          <div v-else class="text-center text-medium-emphasis py-8">
            No historical samples are available for this process yet.
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>
  </PageContent>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import PageContent from "@/components/PageContent.vue";
import WidgetLineChart from "@/components/WidgetLineChart.vue";
import { useProcessStore } from "@/stores/process";
import type { ProcessMetric } from "@/types/process.type";

const processStore = useProcessStore();
const search = ref("");
const refreshSeconds = ref(5);
const isRefreshing = ref(true);
const selectedProcess = ref<ProcessMetric | null>(null);
const historyDialog = ref(false);
const historyFilter = ref("1h");
let refreshTimer: number | null = null;

const refreshItems = [
  { title: "5 seconds", value: 5 },
  { title: "10 seconds", value: 10 },
  { title: "30 seconds", value: 30 },
];
const historyItems = [
  { title: "Last 1 hour", value: "1h" },
  { title: "Last 6 hours", value: "6h" },
  { title: "Last 12 hours", value: "12h" },
  { title: "Last 1 day", value: "1d" },
  { title: "Last 7 days", value: "7d" },
  { title: "Last 30 days", value: "30d" },
];
const headers = [
  { title: "PID", key: "pid", width: 90 },
  { title: "Process", key: "name" },
  { title: "Status", key: "status" },
  { title: "CPU", key: "cpu_percent", align: "end" as const },
  { title: "Memory", key: "memory_percent", align: "end" as const },
  { title: "Threads", key: "thread_count", align: "end" as const },
  { title: "", key: "actions", sortable: false, align: "end" as const },
];

const totalCPU = computed(() =>
  processStore.processes.reduce((total, process) => total + process.cpu_percent, 0)
);
const totalMemory = computed(() =>
  processStore.processes.reduce((total, process) => total + process.memory_bytes, 0)
);

function formatPercent(value: number) {
  return `${Number(value || 0).toFixed(2)}%`;
}

function formatBytes(value: number) {
  if (!value) return "0 B";
  const units = ["B", "KiB", "MiB", "GiB", "TiB"];
  let size = value;
  let unit = 0;
  while (size >= 1024 && unit < units.length - 1) {
    size /= 1024;
    unit++;
  }
  return `${size.toFixed(1)} ${units[unit]}`;
}

function usageClass(value: number) {
  if (value >= 80) return "text-error font-weight-bold";
  if (value >= 50) return "text-warning font-weight-medium";
  return "";
}

function normalizeProcessItem(process: ProcessMetric) {
  const tableItem = process as ProcessMetric & { raw?: ProcessMetric };
  return tableItem.raw ?? process;
}

function statusColor(status: string) {
  if (status === "running") return "success";
  if (status === "zombie" || status === "dead") return "error";
  return "default";
}

async function fetchProcesses() {
  await processStore.fetchLiveProcesses(search.value);
}

function stopRefreshTimer() {
  if (refreshTimer !== null) {
    window.clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

function startRefreshTimer() {
  stopRefreshTimer();
  if (!isRefreshing.value) return;
  refreshTimer = window.setInterval(fetchProcesses, refreshSeconds.value * 1000);
}

function toggleRefresh() {
  isRefreshing.value = !isRefreshing.value;
  startRefreshTimer();
}

async function openHistory(process: ProcessMetric) {
  selectedProcess.value = process;
  historyDialog.value = true;
  await processStore.fetchProcessHistory(process.process_key, historyFilter.value);
}

function closeHistory() {
  historyDialog.value = false;
  selectedProcess.value = null;
}

watch(search, fetchProcesses);
watch(refreshSeconds, startRefreshTimer);
watch(historyFilter, async (filter) => {
  if (selectedProcess.value) {
    await processStore.fetchProcessHistory(selectedProcess.value.process_key, filter);
  }
});

onMounted(async () => {
  await fetchProcesses();
  startRefreshTimer();
});

onUnmounted(stopRefreshTimer);
</script>
