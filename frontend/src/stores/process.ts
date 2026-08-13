import axios from "@/plugins/axios";
import { defineStore } from "pinia";
import type { ProcessHistoryPoint, ProcessMetric } from "@/types/process.type";

export const useProcessStore = defineStore("process", {
  state: () => ({
    processes: [] as ProcessMetric[],
    history: [] as ProcessHistoryPoint[],
    historyError: "",
    loading: false,
    loadingHistory: false,
  }),
  actions: {
    async fetchLiveProcesses(search = "", sort = "cpu") {
      this.loading = true;
      try {
        const response = await axios.get("/processes/live", {
          params: { search, sort, limit: 500 },
        });
        this.processes = response.data.data ?? [];
      } catch (error) {
        this.processes = [];
        console.error("Fetch live processes failed:", error);
      } finally {
        this.loading = false;
      }
      return this.processes;
    },
    async fetchProcessHistory(processKey: string, filter: string) {
      this.loadingHistory = true;
      this.historyError = "";
      try {
        const response = await axios.get("/processes/history", {
          params: { process_key: processKey, filter },
        });
        this.history = response.data.data ?? [];
      } catch (error) {
        this.history = [];
        this.historyError = "Unable to load process history.";
        console.error("Fetch process history failed:", error);
      } finally {
        this.loadingHistory = false;
      }
      return this.history;
    },
  },
});
