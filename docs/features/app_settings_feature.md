# App Settings Feature

## Deskripsi Singkat
Fitur **App Settings** memungkinkan administrator untuk mengelola konfigurasi operasional aplikasi secara dinamis melalui antarmuka web. Konfigurasi ini disimpan di database SQLite utama dan langsung diterapkan pada sistem monitoring serta background workers tanpa perlu melakukan restart aplikasi atau mengubah file konfigurasi YAML secara manual.

## Komponen & Fungsionalitas
1. **Manajemen Konfigurasi Global**:
   - Mengatur interval pengumpulan metrik sistem.
   - Mengatur parameter evaluasi alarm (check interval dan repeat interval).
   - Mengatur ambang batas (*thresholds*) kritis untuk CPU, Disk, dan Memory.
   - Mengatur timeout detak jantung (*heartbeat*) untuk worker services.
2. **Sinkronisasi Otomatis**:
   - Perubahan yang disimpan melalui UI akan langsung dibaca oleh `ThresholdAutomation` pada iterasi pengecekan berikutnya.
3. **Komponen UI**:
   - **Halaman Settings**: Formulir terpusat untuk mengubah parameter sistem dengan validasi tipe data.
   - **Setting Store (Pinia)**: Mengelola state konfigurasi secara reaktif di sisi frontend dan menyediakan sinkronisasi dengan API.

## Arsitektur Teknis (Backend)
- **Model**: `model.AppSetting` (Singleton pattern, hanya menggunakan satu record dengan ID=1).
- **Service**: `service.SettingService` menyediakan antarmuka logika bisnis untuk pengambilan dan pembaruan pengaturan.
- **Repository**: `repository.SettingRepository` menangani persistensi data ke database menggunakan GORM.
- **API Endpoints**:
  - `GET /api/settings`: Mengambil konfigurasi saat ini.
  - `PUT /api/settings`: Memperbarui konfigurasi sistem.

## Konfigurasi yang Dikelola
| Nama Parameter | Deskripsi | Default |
| --- | --- | --- |
| `monitoring_interval_seconds` | Interval pengumpulan metrik sistem (detik). | 10 |
| `alarm_check_interval_seconds` | Seberapa sering sistem mengecek kondisi alarm (detik). | 60 |
| `alarm_repeat_interval_minutes`| Jeda pengulangan notifikasi untuk alarm aktif (menit). | 15 |
| `system_cpu_critical_percent` | Ambang batas penggunaan CPU untuk memicu alarm. | 90.0% |
| `system_disk_critical_percent`| Ambang batas penggunaan Disk untuk memicu alarm. | 85.0% |
| `system_mem_critical_percent` | Ambang batas penggunaan RAM untuk memicu alarm. | 85.0% |
| `worker_heartbeat_timeout`   | Batas waktu worker dianggap *inactive* jika tidak ada detak jantung (detik). | 300 |

## Migrasi dari Config File
Fitur ini secara resmi menggantikan blok `settings:` yang sebelumnya bersifat statis di dalam file `configs/config.yaml`. Dengan pemindahan ke database, sistem menjadi lebih *scalable* dan memungkinkan pengelolaan multi-user yang lebih aman melalui dashboard admin.
