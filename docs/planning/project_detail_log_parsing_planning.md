# Planning: Project Detail & Log Parsing Worker

## Deskripsi Singkat
Dokumen ini merencanakan implementasi fitur "Project Detail" yang akan menampilkan data riwayat log dari suatu proyek. Fitur ini ditenagai oleh sebuah *background worker* yang bertugas mengambil daftar proyek, membaca file log (`.log`) di masing-masing direktori proyek, memparsing isi log berdasarkan pola regex di field `log_format` dari tabel `projects`, dan menyimpannya ke tabel `log_histories`.

Sesuai prinsip arsitektur yang bersih dan efisien (*best practice*), implementasi harus menghindari *N+1 query* saat operasi database dan diwajibkan menggunakan *Batch Insert* untuk memaksimalkan performa saat menyimpan ribuan baris log.

---

## Tahap 1: Backend - Database & Model

**1. Membuat Model `LogHistory`**
- Buat file `internal/model/log_history.go`.
- Definisikan struct `LogHistory` yang merepresentasikan hasil ekstraksi baris log.
  - `ID` (Primary Key, uint)
  - `ProjectID` (uint) - *Foreign Key* ke tabel `Project`.
  - `Timestamp` (time.Time) - Waktu log dicatat.
  - `Level` (string) - Tingkat *severity* log (INFO, ERROR, DEBUG, dll).
  - `Message` (text) - Isi konten dari pesan log.
  - `Source` (string) - Sumber file atau informasi *trace*.
  - `CreatedAt` (time.Time)
- **Relasi**: Perbarui model `Project` (`internal/model/project.go`) agar memiliki `LogHistories []LogHistory` (*One-to-Many*).

**2. Pembaruan Migrasi**
- Tambahkan `model.LogHistory{}` ke `AutoMigrate` di inisialisasi GORM (`cmd/cli/serve.go`).

---

## Tahap 2: Update Engine Log Parser

**1. Modifikasi `internal/logparser/parser.go`**
- Saat ini metode `ParseFile` masih meng-hardcode `Level: "INFO"` dan menetapkan `Timestamp` secara statis dengan `time.Now()`.
- **Tugas**: Perbarui logika regex untuk menangkap *Named Capture Groups* yang sesuai dengan pola `log_format` dari database (contoh: `(?P<timestamp>.*?) \[(?P<level>.*?)\] (?P<message>.*)`).
- Gunakan `regex.SubexpNames()` untuk mengekstrak grup dan memetakannya secara spesifik ke field `LogEntry.Timestamp`, `LogEntry.Level`, dan `LogEntry.Message`.
- Lakukan *parsing* string waktu menjadi tipe `time.Time`.

---

## Tahap 3: Backend - Worker & Batch Processing

**1. Pembuatan Worker Log Parser**
- Buat file `internal/worker/log_parser_worker.go`.
- Alur kerja worker secara berkala:
  1. Lakukan query `GetAllProjects()` (Pastikan hanya query data dasar tanpa memicu N+1 ke tabel lain).
  2. Iterasi pada setiap proyek:
     - Lakukan pemindaian file dengan ekstensi `.log` pada *path* yang didefinisikan proyek.
     - Panggil `logparser.ParseFile()` menggunakan *path* file log dan `project.LogFormat`.
  3. Tampung semua `LogEntry` yang valid dan konversi ke dalam entitas array of `model.LogHistory`.
  4. **Penting**: Simpan data ke database menggunakan operasi *Bulk/Batch Insert* (Contoh: `gorm.DB.CreateInBatches(&logHistories, 1000)`). Jangan menyimpan menggunakan *loop* tunggal.

**2. Integrasi Worker**
- Registrasikan worker ini di dalam `StartWorkers` (`internal/worker/initiator-worker.go`) agar berjalan secara berkala menggunakan `time.Ticker` (interval disarankan 1-5 menit).

---

## Tahap 4: Backend - Repository, Service, & Handler

**1. Repository & Service Layer**
- Buat `internal/repository/log_history_repository.go`. Sediakan fungsionalitas:
  - `BatchInsert(logs []model.LogHistory) error`
  - `GetLogsByProjectID(projectID uint, limit int, offset int) ([]model.LogHistory, error)`
- Buat `internal/services/log_history_service.go` untuk menangani logika penyediaan data ke *Controller*.

**2. HTTP Handler API**
- Tambahkan *endpoint* baru di `project_handler.go` atau handler terpisah:
  - `GET /api/projects/:id/logs` -> Mengembalikan riwayat log dalam format JSON beserta metadata paginasi.

---

## Tahap 5: Frontend - Halaman Project Detail

**1. State Management API**
- Tambahkan *action* `fetchProjectLogs(projectId, page, limit)` di Pinia store (seperti `useProjectStore` atau store baru `useLogHistoryStore`).

**2. UI Halaman Project Detail (`projects/[id].vue`)**
- Buat komponen antarmuka yang menampilkan detail proyek statis di bagian atas.
- Di bagian bawah, sediakan komponen tabel (`v-data-table-server` dari Vuetify) untuk menampilkan daftar log.
- Kolom Tabel: `Timestamp`, `Level`, `Message`.
- **Desain Khusus**: Berikan variasi warna teks atau *chip* (badge) pada kolom `Level` (misalnya: Merah untuk ERROR, Kuning untuk WARN, Biru untuk INFO).
- Sediakan tombol manual untuk *Refresh Log Data* atau sistem auto-polling.

---
**Instruksi Kritis untuk Implementator**:
- Ikuti standar penulisan kode proyek.
- Pastikan logika GORM di worker menghindari query individu (*N+1 issue*) saat integrasi.
- Terapkan mekanisme *error handling* pada pemrosesan file log, sehingga apabila ada satu proyek yang file log-nya rusak atau path tidak valid, *worker* tetap lanjut ke proyek lainnya.
- **Dilarang keras memulai implementasi kode aktual sebelum mendapatkan komando eksekusi dari instruktur/pengguna**.
