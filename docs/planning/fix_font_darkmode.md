# Planning: Fix Bug Font di Donut & Graph pada Dark Mode

## Context
Di halaman `index.vue`, komponen chart yang digunakan (seperti `DiskUsageChart` dan `WidgetLineChart`) masih menggunakan warna font default (hitam/gelap) meskipun aplikasi sedang dalam mode dark mode. Hal ini menyebabkan teks pada sumbu, legend, maupun tooltip menjadi sulit atau tidak bisa dibaca.

## Step by Step
1. **Identifikasi State Tema (Theme)**:
   - Aplikasi menggunakan `useAppStore()` dari `pinia` (file `@/stores/app`) untuk menyimpan state `theme` (`light` atau `dark`).

2. **Perbarui `DiskUsageChart.vue`**:
   - Import `useAppStore` dari `@/stores/app`.
   - Inisialisasi `const appStore = useAppStore();`.
   - Modifikasi `chartOptions` (yang saat ini menggunakan `computed`) dengan menambahkan pengaturan `theme: { mode: appStore.theme as "light" | "dark" }` ke dalam `ApexOptions`.

3. **Perbarui `WidgetLineChart.vue`**:
   - Import `useAppStore` dari `@/stores/app`.
   - Inisialisasi `const appStore = useAppStore();`.
   - Update pengaturan `theme: { mode: appStore.theme as "light" | "dark" }` ke dalam `chartOptions.value` di dalam fungsi `updateChartData`.
   - Tambahkan `watch` baru untuk memantau perubahan `appStore.theme` sehingga jika tema diubah secara real-time, fungsi `updateChartData(props.data)` akan dipanggil lagi untuk mere-render chart dengan tema yang baru.

4. **Testing & Validasi**:
   - Periksa halaman dashboard (`index.vue`).
   - Ubah tema melalui UI dan pastikan font legend, label axis, dan label donut otomatis berubah kontras (putih saat gelap, gelap saat terang).
