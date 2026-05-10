# Planning: Project Management Feature

**Deskripsi Singkat:**
Fitur ini bertujuan untuk mengelola daftar project (aplikasi) yang ada di dalam server (seperti Node.js, Laravel, Python, dll). Fitur "Project" ini akan menjadi acuan (sumber log) bagi fitur Log Parser nantinya. 

Mengingat format log dari berbagai platform bisa berbeda (contoh: Laravel memiliki standar tersendiri, sementara Node.js/Python format log-nya seringkali ditentukan secara custom oleh developer), maka entitas `Project` ini harus menyimpan konfigurasi terkait bagaimana log tersebut diformat/diparsing.

Berikut adalah langkah-langkah detail implementasi yang bisa dikerjakan secara berurutan.

---

## Tahap 1: Backend - Database & Model

**1. Membuat Model `Project`**
- Buat file baru di `internal/model/project.go`.
- Definisikan struct `Project` yang merepresentasikan tabel di database.
- **Field Wajib:**
  - `ID` (Primary Key, Auto Increment/UUID)
  - `Name` (string) - Nama project.
  - `Path` (string) - Absolute path dari project di server (contoh: `/home/user/myproject`).
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)
- **Field Tambahan (Sangat Direkomendasikan untuk fleksibilitas Log Parser):**
  - `ProjectType` (string) - Jenis project (misal: `laravel`, `nodejs`, `python`, `custom`). Ini berguna untuk menentukan metode parsing otomatis.
  - `LogFormat` (string) - Jika `ProjectType` adalah custom/nodejs, field ini menyimpan Regex pattern atau string format yang digunakan oleh Log Parser untuk membaca file log project tersebut.

**2. Update Migrasi Database**
- Tambahkan pembuatan tabel `projects` ke file migrasi database menggunakan GORM.

---

## Tahap 2: Backend - Repository & Service Layer

**1. Membuat Repository Layer**
- Buat file `internal/repository/project_repository.go`.
- Implementasikan fungsi-fungsi CRUD dasar:
  - `CreateProject(project *model.Project) error`
  - `GetAllProjects() ([]model.Project, error)`
  - `GetProjectByID(id uint) (*model.Project, error)`
  - `UpdateProject(project *model.Project) error`
  - `DeleteProject(id uint) error`

**2. Membuat Service Layer**
- Buat file `internal/services/project_service.go`.
- Handle validasi bisnis:
  - Pastikan input `Name` dan `Path` tidak kosong.
  - (Opsional) Validasi apakah absolute `Path` yang dimasukkan benar-benar ada di server atau memiliki izin baca (read permission).
  - Validasi dependensi `LogFormat` berdasarkan `ProjectType`.

---

## Tahap 3: Backend - HTTP Handlers & Router

**1. Membuat Handler API**
- Buat file `internal/handler/http/project_handler.go`.
- Implementasikan endpoint berikut:
  - `GET /api/projects` (Mengambil semua project)
  - `POST /api/projects` (Menambahkan project baru)
  - `PUT /api/projects/:id` (Mengubah data project)
  - `DELETE /api/projects/:id` (Menghapus project)

**2. Mendaftarkan Route**
- Buka file routing backend (biasanya di `internal/handler/initiatior.go` atau file router terkait).
- Daftarkan endpoint-endpoint yang baru dibuat ke dalam grup `/api` (dan pastikan dilindungi oleh JWT Middleware jika diwajibkan).

---

## Tahap 4: Frontend - Store & API Integration

**1. Membuat Pinia Store (atau sejenisnya)**
- Buat file seperti `frontend/src/stores/useProjectStore.js` atau `useProjectStore.ts`.
- Buat action (API calls dengan Axios/Fetch) untuk:
  - `fetchProjects()`
  - `addProject(payload)`
  - `updateProject(id, payload)`
  - `deleteProject(id)`
- Buat state `projects` (array) dan `loading` (boolean).

---

## Tahap 5: Frontend - UI Components & Pages

**1. Membuat Halaman Utama (`projects/index.vue`)**
- Buat folder & file baru: `frontend/src/pages/projects/index.vue`.
- Ikuti standar layout yang ada (seperti pada fitur Users):
  - Gunakan komponen *Breadcrumbs* di bagian atas.
  - Gunakan `v-data-table` (Vuetify) untuk menampilkan list project.
  - Kolom tabel: `Name`, `Path`, `Project Type`, `Created At`, `Actions` (Edit & Delete).

**2. Membuat Dialog Form (Add/Edit)**
- Tambahkan komponen `<v-dialog>` di dalam `index.vue` (atau sebagai komponen terpisah).
- Form input harus berisi:
  - `Name` (Text Field)
  - `Path` (Text Field, berikan hint: contoh `/home/bla/ba/blaaa`)
  - `Project Type` (Select Dropdown: `Laravel`, `Node.js`, `Python`, `Custom`)
  - `Log Format Builder` (Dynamic Components). **Logika UI**: 
    - User tidak menulis regex secara manual.
    - Disediakan builder point-per-point menggunakan dropdown untuk memilih komponen log (`timestamp`, `level`, `message`, `source`) dan pemisah (`separator/text`).
    - Sistem akan merangkai komponen tersebut menjadi sebuah format string terstandar (contoh: `{timestamp} - {level} : {message}`).

**3. Sidebar Navigation**
- Tambahkan menu "Projects" ke dalam sidebar navigasi (biasanya berada di komponen layout utama seperti `App.vue` atau komponen Sidebar khusus).

---

## Tahap 6: Log Parser Integration (Langkah Selanjutnya)
*(Note: Ini untuk gambaran kedepan setelah fitur CRUD selesai)*
- Engine `internal/logparser` nantinya akan di-update agar saat melakukan tugasnya, ia akan melakukan query ke tabel `projects`.
- Menggunakan `Path` project untuk menemukan letak direktori/file log (misal `/home/.../storage/logs/laravel.log`).
- **Dynamic Regex Generation**: Backend akan mengonversi format string dari builder (misal `{timestamp} - {level}`) menjadi *Regex Pattern* yang valid secara otomatis sebelum melakukan parsing log.
