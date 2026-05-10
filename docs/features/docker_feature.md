# Docker Management Feature

## Deskripsi Singkat
Fitur Docker Management memungkinkan pengguna untuk memantau dan mengelola sumber daya Docker langsung dari dalam aplikasi Devopin. Ini memberikan visibilitas penuh terhadap environment Docker yang berjalan.

## Komponen & Fungsionalitas
1. **Container Management**:
   - Tab `Container` untuk melihat daftar container yang ada, baik yang sedang berjalan maupun yang berhenti.
2. **Images Management**:
   - Tab `Images` untuk melihat list image Docker yang tersedia di dalam sistem lokal.
3. **Network Management**:
   - Tab `Network` untuk melihat list network Docker yang telah dibuat beserta detailnya.
4. **Volume Management**:
   - Tab `Volume` untuk memonitor volume Docker.
5. **Real-time / Refresh**:
   - Dilengkapi dengan fitur refresh (`refreshData()`) untuk memuat ulang status terbaru semua entitas Docker.
   - Pada integrasi lebih lanjut, backend Go menghubungkan docker events menggunakan koneksi Socket.io / WebSockets untuk live update.

## Backend Integration
- Backend Go menggunakan Docker Engine API (SDK) untuk melakukan fetch dan manajemen data yang berinteraksi dengan daemon Docker lokal/remote.
