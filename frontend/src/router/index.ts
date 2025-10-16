/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from "vue-router";
import { setupLayouts } from "virtual:generated-layouts";
import { routes } from "vue-router/auto-routes";
import { useAuthStore } from "@/stores/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: setupLayouts(routes),
});

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.("Failed to fetch dynamically imported module")) {
    if (localStorage.getItem("vuetify:dynamic-reload")) {
      console.error("Dynamic import error, reloading page did not fix it", err);
    } else {
      console.log("Reloading page to fix dynamic import error");
      localStorage.setItem("vuetify:dynamic-reload", "true");
      location.assign(to.fullPath);
    }
  } else {
    console.error(err);
  }
});
router.beforeEach(async (to, from, next) => {
  // `to` adalah rute tujuan
  // `from` adalah rute asal
  // `next` adalah fungsi yang harus dipanggil untuk melanjutkan navigasi

  const authStore = useAuthStore(); // Dapatkan akses ke store di luar komponen
  await authStore.me(); // Pastikan status login diperbarui sebelum melanjutkan
  const requiresAuth = to.meta.requiresAuth;
  const guestOnly = to.meta.guestOnly;
  // 1. Cek rute yang butuh login
  if (requiresAuth && !authStore.isLoggedIn) {
    // Jika butuh login tapi belum login, lempar ke halaman login
    next({ name: "/login" }); // Gunakan nama rute jika ada
  }
  // 2. Cek rute yang hanya untuk tamu (belum login)
  else if (guestOnly && authStore.isLoggedIn) {
    // Jika sudah login tapi mencoba akses halaman login/register, lempar ke dashboard
    next({ name: "/" }); // Asumsi nama rute dashboard adalah 'dashboard'
  }
  // 3. Jika semua kondisi aman, izinkan masuk
  else {
    next();
  }
});
router.isReady().then(() => {
  localStorage.removeItem("vuetify:dynamic-reload");
});

export default router;
