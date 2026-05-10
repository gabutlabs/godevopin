# App Settings Feature Migration Planning

## Deskripsi
Dokumen ini berisi panduan langkah demi langkah (step-by-step) untuk memindahkan konfigurasi `settings` yang saat ini berada di dalam file `configs/config.yaml` menjadi fitur dinamis yang tersimpan di dalam database PostgreSQL. Selain itu, dokumen ini mencakup panduan untuk membuat UI pengelolaan *settings* di frontend (Vue.js + Pinia).

Dokumen ini dirancang agar dapat diimplementasikan secara terstruktur oleh Junior Programmer atau AI Assistant, sesuai dengan prinsip arsitektur proyek saat ini.

---

## Tahap 1: Pembuatan Model Database (Backend)

1. **Buat file model baru**: Buat file `internal/model/setting.go`.
2. **Definisikan Struct `AppSetting`**:
   Karena `settings` bersifat global, kita akan menggunakan pola *Singleton* pada tabel database (hanya ada 1 baris record dengan ID = 1).
   ```go
   package model

   import "time"

   type AppSetting struct {
       ID                            uint      `gorm:"primaryKey"`
       MonitoringIntervalSeconds     int       `gorm:"default:10"`
       AlarmCheckIntervalSeconds     int       `gorm:"default:60"`
       AlarmRepeatIntervalMinutes    int       `gorm:"default:15"`
       SystemCPUCriticalPercent      float64   `gorm:"default:90.0"`
       SystemDiskCriticalPercent     float64   `gorm:"default:85.0"`
       SystemMemCriticalPercent      float64   `gorm:"default:85.0"`
       WorkerHeartbeatTimeoutSeconds int       `gorm:"default:300"`
       UpdatedAt                     time.Time
   }
   ```

## Tahap 2: Database Migration & Seeder (Backend)

1. **Update file `internal/database/db.go` (atau file inisialisasi GORM yang relevan)**:
   - Tambahkan `model.AppSetting{}` ke dalam list fungsi `AutoMigrate(...)`.
2. **Buat fungsi Seeder**:
   - Setelah proses `AutoMigrate` berjalan sukses, tambahkan logika pengecekan: Jika tabel `app_settings` kosong (jumlah record = 0), maka lakukan `INSERT` 1 baris data dengan ID = 1 yang berisi nilai default. Nilai default dapat mengacu pada struct di atas atau menyalin dari konfigurasi `config.yaml` lama.

## Tahap 3: Pembuatan Repository Layer (Backend)

1. **Buat file `internal/repository/setting_repository.go`**.
2. **Buat interface dan struct implementasinya**:
   ```go
   type SettingRepository interface {
       GetSettings() (*model.AppSetting, error)
       UpdateSettings(setting *model.AppSetting) error
   }
   ```
3. **Implementasikan Method GORM**:
   - `GetSettings()`: Lakukan query `db.First(&setting, 1)`.
   - `UpdateSettings()`: Lakukan operasi `db.Save(&setting)` (atau `Updates`) dengan memastikan nilai `setting.ID` selalu `1`.

## Tahap 4: Pembuatan Service Layer (Backend)

1. **Buat file `internal/services/setting_service.go`**.
2. **Buat interface dan struct implementasinya**:
   ```go
   type SettingService interface {
       GetSettings() (*model.AppSetting, error)
       UpdateSettings(req *UpdateSettingRequest) error
   }
   ```
3. **Tambahkan DTO/Request Struct `UpdateSettingRequest`**: 
   - Struct ini digunakan untuk memvalidasi input dari REST API (handler), yang mencakup semua field di `AppSetting` (tanpa field `ID`).

## Tahap 5: Pembuatan Handler dan Router (Backend)

1. **Buat file `internal/handler/setting_handler.go`**.
2. **Buat struct `SettingHandler`** dengan method:
   - `GetSettings`: Memanggil fungsi `GetSettings` di Service, dan mengembalikan response JSON `AppSetting`.
   - `UpdateSettings`: Menerima JSON body, melakukan validasi/binding ke `UpdateSettingRequest`, lalu memanggil service untuk menyimpan perubahan.
3. **Update Router (`internal/web/router.go` atau file router utama)**:
   - Daftarkan endpoint API baru:
     - `GET /api/settings` (Bisa bersifat publik atau hanya internal)
     - `PUT /api/settings` (Wajib dilindungi oleh middleware autentikasi).

## Tahap 6: Refactoring Penggunaan Config Lama (Backend)

1. **Cari semua referensi `config.Settings`**:
   Cari kode di package lain (misalnya di `internal/monitoring/`, `internal/worker/`, atau `internal/config/config.go`) yang saat ini menggunakan atau membaca nilai dari `viper` (contoh: `viper.GetInt("settings.monitoring_interval_seconds")`) atau field `Config.Settings`.
2. **Ganti Sumber Data**:
   Lakukan *dependency injection* dari `SettingRepository` (atau *Service*) ke komponen-komponen tersebut agar dapat mengambil nilai interval/thresholds terbaru langsung dari database setiap kali dibutuhkan, bukan dari memory config yang statis.
3. **Hapus Config Lama**:
   Setelah semua penggunaan direfaktor, hapus struct `Settings`, `Alarms`, dan `Thresholds` dari `internal/config/config.go`. Hapus juga scope `settings:` beserta seluruh isinya dari file `configs/config.yaml`.

## Tahap 7: Implementasi Frontend (Vue + Pinia)

1. **Buat API Service (`frontend/src/services/settingService.js` atau `.ts`)**:
   - Buat method `getSettings()` yang melakukan Axios `GET` ke `/api/settings`.
   - Buat method `updateSettings(data)` yang melakukan Axios `PUT` ke `/api/settings`.
2. **Buat Pinia Store (`frontend/src/stores/settingStore.js` atau `.ts`)**:
   - Deklarasikan state reaktif untuk menampung data konfigurasi.
   - Buat action `fetchSettings` untuk memuat data.
   - Buat action `saveSettings` untuk menyimpan dan memperbarui state ketika perubahan berhasil disubmit ke server.
3. **Buat Halaman UI (`frontend/src/views/Settings.vue` atau lokasi komponen pages yang sesuai)**:
   - Desain form input untuk memodifikasi parameter:
     - Monitoring Interval
     - Alarm Check Interval
     - Alarm Repeat Interval
     - CPU Critical Threshold (%)
     - Disk Critical Threshold (%)
     - RAM Critical Threshold (%)
     - Worker Heartbeat Timeout (detik)
   - Tambahkan tombol "Save" dengan animasi loading/sukses.
4. **Registrasi Rute dan Navigasi**:
   - Daftarkan rute `/settings` di router Vue.
   - Tambahkan menu tautan "Settings" di komponen Sidebar atau Navbar aplikasi.

---
**Catatan Penting untuk Implementator**:
- Selalu ikuti struktur *Clean Architecture* yang telah ada di repository ini (Handler -> Service -> Repository).
- Jangan menghapus file `configs/config.yaml` sepenuhnya. Cukup hapus blok kode di bawah kunci `settings:`. Kunci lain seperti `database:` dan `app:` harus dibiarkan karena masih dibutuhkan oleh sistem.
- Berikan respons log atau error message yang informatif jika proses update parameter ke database gagal.
