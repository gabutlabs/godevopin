import axios from "@/plugins/axios";
import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    // Kita asumsikan token disimpan di localStorage saat login
    token: localStorage.getItem("token") || null,
    user: localStorage.getItem("user") || null, // Bisa diisi data user setelah login
    loggedIn: false,
  }),
  getters: {
    // Getter ini akan menjadi "pengecek status" utama kita
    isLoggedIn: (state) => state.loggedIn,
  },
  actions: {
    async login(payload: {
      email: string;
      password: string;
    }): Promise<boolean> {
      try {
        const response = await axios.post("/auth/login", payload);
        localStorage.setItem("token", response.data.data.token);
        localStorage.setItem("user", JSON.stringify(response.data.data.user));
        this.token = response.data.data.token;
        this.user = response.data.data.user;
        // Arahkan ke halaman dashboard atau halaman yang diinginkan setelah login
        return true;
      } catch (error) {
        console.error("Login failed:", error);
        return false;
      }
    },
    async me() {
      try {
        const response = await axios.get("/auth/me");
        this.user = response.data.data.user;
        localStorage.setItem("user", JSON.stringify(response.data.data.user));
        this.loggedIn = response.data.data.is_loggedin;
      } catch (error) {
        console.error("Fetch user profile failed:", error);
        this.loggedIn = false;
      }
    },
    logout() {
      this.token = null;
      localStorage.removeItem("token");
      this.user = null;
      localStorage.removeItem("user");
      this.loggedIn = false;
      return true;
      // Arahkan ke halaman login
      // this.router.push('/login')
    },
  },
});
