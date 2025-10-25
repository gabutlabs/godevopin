import axios from "@/plugins/axios";
import { defineStore } from "pinia";
export type ActiveAlarm = {
  alarm_name: string;
  target: string;
  status: "FIRING" | "ACKNOWLEDGE" | string;
  started_at: string;
  last_notified_at: string;
  message: string;
};
export type HistoryAlarm = {
  created_at: string;
  metadata: { [key: string]: any };
} & Pick<ActiveAlarm, "alarm_name" | "target" | "status" | "message">;
export const useAlarmStore = defineStore("alarm", {
  state: () => ({
    activeAlarms: [] as ActiveAlarm[],
    historyAlarms: { data: [], pagination: undefined } as {
      data: HistoryAlarm[];
      pagination: { total: number; page: number; limit: number } | undefined;
    },
    activeAlarm: null,
    historyAlarm: null,
    action_result: { is_success: false, message: "", data: null } as {
      is_success: boolean;
      message: string;
      data: any;
    },
  }),
  actions: {
    async fetchActiveAlarms() {
      try {
        const response = await axios.get("/alarms/active");
        this.activeAlarms = response.data.data;
      } catch (error) {
        console.error("Fetch users failed:", error);
      }
    },
    async fetchHistoryAlarms(page: number, limit: number = 10) {
      try {
        const response = await axios.get(
          `/alarms/history?page=${page}&limit=${limit}`
        );
        this.historyAlarms = response.data.data;
      } catch (error) {
        console.error("Fetch users failed:", error);
      }
    },
  },
});
