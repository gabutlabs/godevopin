import axios from "@/plugins/axios";
import { AxiosError } from "axios";
import { defineStore } from "pinia";

export type User = {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
};
export const useUserStore = defineStore("user", {
  state: () => ({
    profile: localStorage.getItem("user")
      ? JSON.parse(localStorage.getItem("user") as string)
      : null,
    users: [] as Array<User>,
    user: null as User | null,
    action_result: { is_success: false, message: "", data: null } as {
      is_success: boolean;
      message: string;
      data: any;
    },
  }),
  getters: {
    getProfile: (state) => state.profile,
  },
  actions: {
    async setFilterUser(name: string) {
      if (name.trim() === "") {
        await this.fetchAllUsers();
      } else {
        this.users = this.users.filter((user) =>
          user.name.toLowerCase().includes(name.toLowerCase())
        );
      }
    },
    setProfile(profile: object) {
      this.profile = profile;
      localStorage.setItem("user", JSON.stringify(profile));
    },
    clearProfile() {
      this.profile = null;
      localStorage.removeItem("user");
    },
    async fetchAllUsers() {
      try {
        const response = await axios.get("/users");
        this.users = response.data.data;
      } catch (error) {
        console.error("Fetch users failed:", error);
      }
    },
    async fetchUser(id: number) {
      try {
        const response = await axios.get(`/users/${id}`);
        this.user = response.data.data;
      } catch (error) {
        console.error("Fetch users failed:", error);
      }
    },
    async createUser(payload: Omit<User, "id" | "created_at" | "updated_at">) {
      try {
        const response = await axios.post(
          "/users",
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
    async updateUser(
      id: number,
      payload: Omit<User, "id" | "created_at" | "updated_at">
    ) {
      try {
        const response = await axios.put(
          `/users/${id}`,
          payload
        );
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
      } catch (error) {
        console.error("Update user failed:", error);
        let message = "Failed creating user";
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
    async deleteUser(id: number) {
      try {
        const response = await axios.delete(
          `/users/${id}`
        );
        this.action_result = {
          is_success: true,
          data: null,
          message: response.data.message,
        };
      } catch (error) {
        console.error("Delete user failed:", error);
        let message = "Failed creating user";
        if (error instanceof AxiosError) {
          message = error.response?.data.message;
        }
        this.action_result = { is_success: false, data: null, message };
      }
    },
  },
});
