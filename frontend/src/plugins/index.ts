/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from "./vuetify";
import pinia from "../stores";
import router from "../router";
import axios from "./axios";

// Types
import type { App } from "vue";
import VueApexCharts from "vue3-apexcharts";
import { createNotivue } from "notivue";
import "notivue/notification.css"; // Only needed if using built-in notifications
import "notivue/animations.css"; // Only needed if using built-in animations
export function registerPlugins(app: App) {
  const notivue = createNotivue({
    notifications: {
      global: {
        duration: 10000,
      },
    },
  });
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(notivue)
    .component("apexchart", VueApexCharts);

  // Make axios available globally
  app.config.globalProperties.$http = axios;
}
