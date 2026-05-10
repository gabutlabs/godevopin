import axios from "@/plugins/axios";
import { AxiosError } from "axios";
import { defineStore } from "pinia";

export type AppSetting = {
  id: number;
  monitoring_interval_seconds: number;
  alarm_check_interval_seconds: number;
  alarm_repeat_interval_minutes: number;
  system_cpu_critical_percent: number;
  system_disk_critical_percent: number;
  system_mem_critical_percent: number;
  worker_heartbeat_timeout_seconds: number;
  updated_at: string;
};

export type UpdateSettingRequest = Omit<AppSetting, "id" | "updated_at">;

export const useSettingStore = defineStore("setting", {
  state: () => ({
    settings: null as AppSetting | null,
    loading: false,
    actionResult: { isSuccess: false, message: "" } as {
      isSuccess: boolean;
      message: string;
    },
  }),
  getters: {
    getSettings: (state) => state.settings,
  },
  actions: {
    async fetchSettings() {
      this.loading = true;
      try {
        const response = await axios.get("/settings");
        this.settings = response.data;
      } catch (error) {
        console.error("Fetch settings failed:", error);
      } finally {
        this.loading = false;
      }
    },

    async updateSettings(payload: UpdateSettingRequest) {
      this.loading = true;
      try {
        const response = await axios.put("/settings", payload);
        this.actionResult = {
          isSuccess: true,
          message: response.data.message,
        };
        await this.fetchSettings();
      } catch (error) {
        console.error("Update settings failed:", error);
        let message = "Failed to update settings";
        if (error instanceof AxiosError) {
          message = error.response?.data.error || message;
        }
        this.actionResult = {
          isSuccess: false,
          message: message,
        };
      } finally {
        this.loading = false;
      }
    },
  },
});
