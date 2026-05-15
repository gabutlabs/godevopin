<template>
  <v-card title="Storage Disk">
    <div v-if="props.total > 0">
      <VueApexCharts
        type="donut"
        height="350"
        :options="chartOptions"
        :series="series"
      />
    </div>
    <div v-else class="loading-state">Data penyimpanan tidak tersedia.</div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import VueApexCharts from "vue3-apexcharts";
import type { ApexOptions } from "apexcharts";
import { useAppStore } from "@/stores/app";

// --- Props untuk menerima data dari luar ---
const appStore = useAppStore();
const props = defineProps({
  usage: {
    type: Number,
    required: true,
    default: 0,
  },
  total: {
    type: Number,
    required: true,
    default: 0, // Default 0, akan menampilkan pesan loading
  },
});

// --- Helper Function untuk format Bytes ---
function formatBytes(bytes: number, decimals = 2): string {
  if (!+bytes) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

// --- Computed Properties untuk reaktivitas ---

// Menghitung sisa ruang disk
const freeSpace = computed(() =>
  props.total > props.usage ? props.total - props.usage : 0
);

// Menghitung persentase penggunaan
const usagePercentage = computed(() => {
  if (props.total === 0) return 0;
  return (props.usage / props.total) * 100;
});

// Data series untuk pie chart
const series = computed(() => [props.usage, freeSpace.value]);
// Konfigurasi lengkap untuk chart
const chartOptions = computed(
  (): ApexOptions => ({
    theme: {
      mode: appStore.theme as "light" | "dark",
    },
    chart: {
      type: "donut",
      background: "transparent",
    },
    // Label untuk setiap irisan
    labels: ["Ruang Digunakan", "Ruang Kosong"],
    // Warna untuk setiap irisan
    colors: ["#008FFB", "#00E396"],
    // Pengaturan spesifik untuk Donut Chart
    plotOptions: {
      pie: {
        donut: {
          size: "65%",
          labels: {
            show: true,
            total: {
              show: true,
              showAlways: true,
              label: "Digunakan",
              fontSize: "18px",
              fontWeight: 600,
              // Menampilkan persentase penggunaan di tengah
              formatter: () => usagePercentage.value.toFixed(1) + "%",
            },
          },
        },
      },
    },
    // Mengatur label data di setiap irisan (menampilkan persentase)
    dataLabels: {
      enabled: true,
      formatter: (val, opts) => {
        console.log(opts.w.globals);
        return (
          opts.w.globals.seriesPercent[opts.seriesIndex][0].toFixed(1) + "%"
        );
      },
    },
    // Mengatur tooltip saat di-hover (menampilkan ukuran dalam GB/MB/KB)
    tooltip: {
      y: {
        formatter: (val: number) => formatBytes(val),
      },
    },
    legend: {
      position: "bottom",
    },
  })
);
</script>

<style scoped>
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 350px;
  color: #888;
}
</style>
