<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <div>
    <v-toolbar
      density="comfortable"
      :color="state.theme == 'light' ? 'white' : '#2C303A'"
      title="Dashboard"
    >
      <template #append>
        <div class="d-flex ga-3 mx-4 align-center">
          <v-select
            label="Filter Waktu"
            density="compact"
            :items="filterItems"
            style="min-width: 180px"
            hide-details
            v-model="selectFilter"
          ></v-select>
          <v-btn
            :prepend-icon="autoRefresh.icon.value"
            variant="outlined"
            :color="autoRefresh.isEnabled.value ? 'success' : 'primary'"
            @click="autoRefresh.toggle"
            density="default"
          >
            {{ autoRefresh.text.value }}
          </v-btn>
        </div>
      </template>
    </v-toolbar>
    <br />
    <v-row>
      <v-col cols="12">
        <DiskUsageChart
          :usage="widgetState.diskUsage.disk_usage_byte"
          :total="widgetState.diskUsage.disk_total_byte"
        />
      </v-col>
      <v-col cols="12">
        <WidgetLineChart
          :data="widgetState.filteredMetrics"
          title="CPU Usage"
          data-key="avg_cpu_usage"
          series-name="CPU (%)"
          unit="%"
          :filter-range="selectFilter"
        />
      </v-col>
      <v-col cols="12">
        <WidgetLineChart
          :data="widgetState.filteredMetrics"
          title="Memory Usage"
          data-key="avg_mem_usage"
          series-name="Memory (MiB)"
          unit="MiB"
          :filter-range="selectFilter"
        />
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";
import { useAppStore } from "@/stores/app";
import { useWidgetStore } from "@/stores/widget";
import WidgetLineChart from "@/components/WidgetLineChart.vue";
import DiskUsageChart from "@/components/DiskUsageChart.vue";

// --- State untuk UI ---
const selectFilter = ref("1h");
const filterItems = ref([
  { title: "Last 1 hour", value: "1h" },
  { title: "Last 6 hour", value: "6h" },
  { title: "Last 12 hour", value: "12h" },
  { title: "Last 1 day", value: "1d" },
  { title: "Last 7 day", value: "7d" },
  { title: "Last 30 day", value: "30d" },
]);

// --- Pinia Stores ---
const state = useAppStore();
const widgetState = useWidgetStore();

// --- Fungsi untuk mengambil data ---
async function fetchData(filter: string) {
  // Tambahkan loading state jika perlu
  await widgetState.getFilteredMetrics(filter);
  await widgetState.getDiskUsage();
}

// --- Logika Auto Refresh ---
const useAutoRefresh = (callback: () => void) => {
  const isEnabled = ref(false);
  let intervalId: number | null = null;

  const toggle = () => {
    isEnabled.value = !isEnabled.value;
    if (isEnabled.value) {
      // Jalankan callback langsung saat diaktifkan
      callback();
      // Set interval untuk refresh setiap 30 detik (sesuaikan)
      intervalId = window.setInterval(callback, 30000);
    } else {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = null;
      }
    }
  };

  const text = computed(() =>
    isEnabled.value ? "Auto Refresh: On" : "Enable Auto Refresh"
  );
  const icon = computed(() => (isEnabled.value ? "mdi-sync" : "mdi-sync-off"));

  // Pastikan interval dibersihkan saat komponen di-unmount untuk mencegah memory leak
  onUnmounted(() => {
    if (intervalId) {
      clearInterval(intervalId);
    }
  });

  return { isEnabled, toggle, text, icon };
};

// Gunakan auto-refresh dengan fungsi fetchData
const autoRefresh = useAutoRefresh(() => fetchData(selectFilter.value));

// --- Reaktivitas & Lifecycle Hooks ---
// Gunakan `watch` untuk memantau perubahan filter. Ini lebih deklaratif.
watch(selectFilter, (newFilter) => {
  fetchData(newFilter);
});

// Ambil data awal saat komponen di-mount
onMounted(() => {
  fetchData(selectFilter.value);
});
</script>
