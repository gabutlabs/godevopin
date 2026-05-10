<template>
  <v-card :title="title">
    <div v-if="isChartReady">
      <VueApexCharts
        type="line"
        height="320"
        :options="chartOptions"
        :series="series"
      />
    </div>
    <div v-else class="loading-state">Memuat data grafik...</div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import VueApexCharts from "vue3-apexcharts";
import type { ApexOptions } from "apexcharts";
import { useAppStore } from "@/stores/app";

const props = defineProps({
  // ... (props Anda yang lain tidak berubah)
  data: { type: Array as () => any[], required: true },
  title: { type: String, required: true },
  dataKey: { type: String, required: true },
  seriesName: { type: String, required: true },
  unit: { type: String, default: "" },

  // --- TAMBAHKAN PROP BARU INI ---
  // Menerima nilai filter seperti '1h', '6h', '7d', '30d'
  filterRange: {
    type: String,
    required: true,
  },
});

const appStore = useAppStore();
const series = ref<any[]>([]);
const chartOptions = ref<ApexOptions>({});
const isChartReady = ref(false);

// ... (fungsi formatBytes tidak berubah)

async function updateChartData(newData: any[]) {
  isChartReady.value = false;

  if (!newData || newData.length === 0) {
    series.value = [];
    await nextTick();
    isChartReady.value = true;
    return;
  }

  const chartData = newData.map((item: any) => ({
    x: new Date(item.time_interval).getTime(),
    y: parseFloat(
      props.dataKey == "avg_mem_usage"
        ? item[props.dataKey] / 1024 / 1024
        : item[props.dataKey]
    ).toFixed(2),
  }));

  series.value = [{ name: props.seriesName, data: chartData }];

  // --- LOGIKA UTAMA ADA DI SINI ---
  // Tentukan format label berdasarkan prop `filterRange`
  let xAxisFormat: string;
  let xAxisTitle: string;

  if (props.filterRange.includes("h") || props.filterRange === "1d") {
    // Jika filter mengandung 'h' (hour), gunakan format Jam:Menit
    xAxisFormat = "HH:mm";
    xAxisTitle = `Waktu (Hari ini)`;
  } else {
    // Jika filter mengandung 'd' (day), gunakan format Tanggal Bulan
    xAxisFormat = "dd MMM";
    xAxisTitle = "Tanggal";
  }

  chartOptions.value = {
    theme: {
      mode: appStore.theme as "light" | "dark",
    },
    // ... (opsi chart lainnya)
    xaxis: {
      type: "datetime",
      labels: {
        datetimeUTC: false,
        format: xAxisFormat, // <-- Gunakan format dinamis
      },
      title: {
        text: xAxisTitle, // <-- Gunakan judul dinamis
      },
    },
    tooltip: {
      x: {
        format: "dd MMM yyyy - HH:mm", // Tooltip tetap detail
      },
      // ... (sisa konfigurasi tooltip)
    },
    // ... (sisa konfigurasi chart Anda)
  };

  await nextTick();
  isChartReady.value = true;
}

watch(
  () => props.data,
  (newData) => {
    updateChartData(newData);
  },
  { immediate: true, deep: true }
);

watch(
  () => appStore.theme,
  () => {
    updateChartData(props.data);
  }
);
</script>

<style scoped>
/* ... (style Anda tidak berubah) ... */
</style>
