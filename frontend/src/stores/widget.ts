import axios from "@/plugins/axios";
import { defineStore } from "pinia";

export const useWidgetStore = defineStore("widget", {
  state: () => ({
    filteredMetrics: [] as Array<{
      avg_cpu_usage: number;
      avg_mem_usage: number;
      created_at: string;
    }>,
    diskUsage: {
      disk_usage_byte: 0,
      disk_total_byte: 0,
    } as { disk_usage_byte: number; disk_total_byte: number },
  }),
  actions: {
    async getFilteredMetrics(filter: string) {
      try {
        const response = await axios.get(
          `/system-metrics/filtered?filter=${filter}`
        );
        this.filteredMetrics = response.data.data;
      } catch (error) {
        this.filteredMetrics = [];
        console.error("Error fetching filtered metrics:", error);
      }
      return this.filteredMetrics;
    },
    async getDiskUsage() {
      try {
        const response = await axios.get(
          `/system-metrics/disk-usage`
        );
        this.diskUsage = response.data.data;
      } catch (error) {
        console.error("Error fetching filtered metrics:", error);
      }
      return this.diskUsage;
    },
  },
});
