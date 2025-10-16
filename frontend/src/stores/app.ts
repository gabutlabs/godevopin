// Utilities
import { defineStore } from "pinia";

export const useAppStore = defineStore("app", {
  state: () => ({
    theme: localStorage.getItem("theme_mode") || "light",
  }),
  actions: {
    changeTheme() {
      this.theme = this.theme === "light" ? "dark" : "light";
      localStorage.setItem("theme_mode", this.theme);
    },
  },
});
