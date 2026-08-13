# Dashboard Feature (System Monitoring)

## Deskripsi Singkat
Fitur Dashboard merupakan halaman utama yang memberikan visualisasi real-time dan historis terkait performa sistem. Fitur ini memonitoring metrik kunci seperti penggunaan CPU, Memori, dan Disk.

## Komponen & Fungsionalitas
1. **Visualisasi Metrik**:
   - **Disk Usage Chart**: Menampilkan penggunaan disk saat ini dibandingkan dengan total kapasitas.
   - **CPU Usage Chart**: Line chart yang menampilkan rata-rata penggunaan CPU dalam bentuk persentase (%).
   - **Memory Usage Chart**: Line chart yang menampilkan penggunaan memori dalam satuan MiB.
2. **Filter Waktu**:
   - Terdapat filter rentang waktu untuk grafik, mulai dari `Last 1 hour`, `Last 6 hour`, `Last 12 hour`, `Last 1 day`, `Last 7 day`, hingga `Last 30 day`.
3. **Auto Refresh**:
   - Tombol toggle untuk mengaktifkan pembaruan data secara otomatis (polling setiap 30 detik) tanpa perlu memuat ulang halaman.

## Backend & Database
- Data didapatkan dari backend yang mencatat System Metrics.
- Menggunakan database SQLite khusus metrics untuk penyimpanan time-series data tanpa server database terpisah.
- Data disagregasi di repository menggunakan bucket waktu yang kompatibel dengan SQLite untuk efisiensi query grafik berjangka panjang.

## Teknologi
- Frontend: Vue 3, Vuetify, Pinia (`useWidgetStore`), Recharts/Echarts (untuk rendering chart).
- Backend: Go, Fiber (HTTP REST), SQLite.
