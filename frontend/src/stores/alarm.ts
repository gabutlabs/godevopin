import axios from "@/plugins/axios";
import { AxiosError } from "axios";
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
    async fetchActiveAlarmByStatus(status: string) {
      try {
        const response = await axios.get(`/alarms/active/status/${status}`);
        this.activeAlarms = response.data.data;
      } catch (error) {
        console.error("Fetch users failed:", error);
      }
    },
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
    async acknowledgeAlarm(
      param: { alarmName: string; target: string },
      payload: { acknowledged_by: string }
    ) {
      try {
        console.log("Payload:", payload);
        const response = await axios.post(
          `/alarms/active/acknowledge/${param.alarmName}/${param.target}`,
          payload
        );
        this.action_result = {
          is_success: true,
          message: response.data.message,
          data: null,
        };
      } catch (error) {
        let message = "Failed creating user";
        let errors = undefined;
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
          errors = error.response?.data.errors;
        }
        this.action_result = {
          is_success: false,
          message: message,
          data: errors,
        };
      }
    },
  },
});
