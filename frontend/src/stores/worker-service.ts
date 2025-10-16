import axios from "@/plugins/axios";
import { AxiosError } from "axios";
import { defineStore } from "pinia";

// Define TypeScript types for WorkerService
export type WorkerService = {
  id: number;
  name: string;
  description: string;
  desired_state: "enabled" | "disabled";
  current_status: "starting" | "running" | "stopped" | "failed" | "degraded";
  health_status: "healthy" | "unhealthy" | "unknown";
  last_heartbeat_at?: string;
  last_success_at?: string;
  last_failure_at?: string;
  last_error_message?: string;
  pid?: number;
  version?: string;
  created_at: string;
  updated_at: string;
};

// Type for creating a new worker service
export type CreateWorkerServiceRequest = {
  name: string;
  description: string;
  desired_state: "enabled" | "disabled";
};

// Type for updating a worker service
export type UpdateWorkerServiceRequest = {
  name: string;
  description: string;
  desired_state: "enabled" | "disabled";
};

// Type for updating status
export type UpdateWorkerServiceStatusRequest = {
  current_status: "starting" | "running" | "stopped" | "failed" | "degraded";
  health_status: "healthy" | "unhealthy" | "unknown";
};

export const useWorkerServiceStore = defineStore("worker-service", {
  state: () => ({
    workerServices: [] as Array<WorkerService>,
    workerService: null as WorkerService | null,
    action_result: { is_success: false, message: "", data: null } as {
      is_success: boolean;
      message: string;
      data: any;
    },
  }),
  getters: {
    getWorkerServices: (state) => state.workerServices,
    getWorkerService: (state) => state.workerService,
  },
  actions: {
    // Filter worker services by name
    async setFilterWorkerService(name: string) {
      if (name.trim() === "") {
        await this.fetchAllWorkerServices();
      } else {
        this.workerServices = this.workerServices.filter((ws) =>
          ws.name.toLowerCase().includes(name.toLowerCase())
        );
      }
    },

    // Fetch all worker services
    async fetchAllWorkerServices() {
      try {
        const response = await axios.get("/worker-services");
        this.workerServices = response.data.data;
      } catch (error) {
        console.error("Fetch worker services failed:", error);
      }
    },

    // Fetch worker service by ID
    async fetchWorkerService(id: number) {
      try {
        const response = await axios.get(`/worker-services/${id}`);
        this.workerService = response.data.data;
      } catch (error) {
        console.error("Fetch worker service failed:", error);
      }
    },

    // Fetch worker service by name
    async fetchWorkerServiceByName(name: string) {
      try {
        const response = await axios.get(`/worker-services/name/${name}`);
        this.workerService = response.data.data;
      } catch (error) {
        console.error("Fetch worker service by name failed:", error);
      }
    },

    // Create a new worker service
    async createWorkerService(
      payload: Omit<
        CreateWorkerServiceRequest,
        "id" | "created_at" | "updated_at"
      >
    ) {
      try {
        const response = await axios.post("/worker-services", payload);
        this.action_result = {
          is_success: true,
          message: response.data.message,
          data: null,
        };
        // Refresh the list after creation
        await this.fetchAllWorkerServices();
      } catch (error) {
        let message = "Failed creating worker service";
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

    // Update an existing worker service
    async updateWorkerService(
      id: number,
      payload: Omit<
        UpdateWorkerServiceRequest,
        "id" | "created_at" | "updated_at"
      >
    ) {
      try {
        const response = await axios.put(`/worker-services/${id}`, payload);
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        // Refresh the specific worker service after update
        await this.fetchWorkerService(id);
        // Refresh the list too
        await this.fetchAllWorkerServices();
      } catch (error) {
        console.error("Update worker service failed:", error);
        let message = "Failed updating worker service";
        let errors = undefined;
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
          errors = error.response?.data.errors;
        }
        this.action_result = {
          is_success: false,
          message,
          data: errors,
        };
      }
    },

    // Delete a worker service
    async deleteWorkerService(id: number) {
      try {
        const response = await axios.delete(`/worker-services/${id}`);
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        // Refresh the list after deletion
        await this.fetchAllWorkerServices();
      } catch (error) {
        console.error("Delete worker service failed:", error);
        let message = "Failed deleting worker service";
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
        }
        this.action_result = { is_success: false, data: null, message };
      }
    },

    // Update worker service status
    async updateWorkerServiceStatus(
      id: number,
      payload: UpdateWorkerServiceStatusRequest
    ) {
      try {
        const response = await axios.put(
          `/worker-services/${id}/status`,
          payload
        );
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        // Refresh the specific worker service after status update
        await this.fetchWorkerService(id);
      } catch (error) {
        console.error("Update worker service status failed:", error);
        let message = "Failed updating worker service status";
        let errors = undefined;
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
          errors = error.response?.data.errors;
        }
        this.action_result = {
          is_success: false,
          message,
          data: errors,
        };
      }
    },

    // Update worker service heartbeat
    async updateWorkerServiceHeartbeat(id: number) {
      try {
        const response = await axios.put(
          `/worker-services/${id}/heartbeat`,
          {}
        );
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        // Refresh the specific worker service after heartbeat update
        await this.fetchWorkerService(id);
      } catch (error) {
        console.error("Update worker service heartbeat failed:", error);
        let message = "Failed updating worker service heartbeat";
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
        }
        this.action_result = { is_success: false, data: null, message };
      }
    },
    // Restart worker service - updates the status to starting (to simulate restart)
    async updateStatusWorkerService(id: number, status: string) {
      try {
        const response = await axios.put(`/worker-services/${id}/status`, {
          current_status: status,
          health_status: "unknown",
        });
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
      } catch (error) {
        console.error("Restart worker service failed:", error);
        let message = "Failed restarting worker service";
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
        }
        this.action_result = { is_success: false, data: null, message };
      }
    },
  },
});
