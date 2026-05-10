import axios from "@/plugins/axios";
import { AxiosError } from "axios";
import { defineStore } from "pinia";

export type Project = {
  id: number;
  name: string;
  path: string;
  project_type: string;
  log_format: string;
  created_at: string;
  updated_at: string;
};

export type LogHistory = {
  id: string;
  project_id: number;
  timestamp: string;
  level: string;
  message: string;
  source: string;
  created_at: string;
};

export const useProjectStore = defineStore("project", {
  state: () => ({
    projects: [] as Array<Project>,
    project: null as Project | null,
    logs: [] as Array<LogHistory>,
    logsTotal: 0,
    loading: false,
    action_result: { is_success: false, message: "", data: null } as {
      is_success: boolean;
      message: string;
      data: any;
    },
  }),
  actions: {
    async fetchProjectLogs(id: number, page: number = 1, limit: number = 100) {
      this.loading = true;
      try {
        const response = await axios.get(`/projects/${id}/logs`, {
          params: { page, limit }
        });
        this.logs = response.data.data.data;
        this.logsTotal = response.data.data.total;
      } catch (error) {
        console.error("Fetch project logs failed:", error);
      } finally {
        this.loading = false;
      }
    },
    async fetchAllProjects() {
      this.loading = true;
      try {
        const response = await axios.get("/projects");
        this.projects = response.data.data;
      } catch (error) {
        console.error("Fetch projects failed:", error);
      } finally {
        this.loading = false;
      }
    },
    async fetchProject(id: number) {
      try {
        const response = await axios.get(`/projects/${id}`);
        this.project = response.data.data;
      } catch (error) {
        console.error("Fetch project failed:", error);
      }
    },
    async createProject(payload: Omit<Project, "id" | "created_at" | "updated_at">) {
      try {
        const response = await axios.post("/projects", payload);
        this.action_result = {
          is_success: true,
          message: response.data.message,
          data: null,
        };
        await this.fetchAllProjects();
      } catch (error) {
        let message = "Failed creating project";
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
    async updateProject(id: number, payload: Omit<Project, "id" | "created_at" | "updated_at">) {
      try {
        const response = await axios.put(`/projects/${id}`, payload);
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        await this.fetchAllProjects();
      } catch (error) {
        console.error("Update project failed:", error);
        let message = "Failed updating project";
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
    async deleteProject(id: number) {
      try {
        const response = await axios.delete(`/projects/${id}`);
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
        await this.fetchAllProjects();
      } catch (error) {
        console.error("Delete project failed:", error);
        let message = "Failed deleting project";
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
        }
        this.action_result = { is_success: false, data: null, message };
      }
    },
  },
});
