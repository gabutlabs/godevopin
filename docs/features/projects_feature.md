# Projects Management Feature

## Deskripsi Singkat
Fitur **Projects Management** dirancang untuk mengelola daftar aplikasi atau proyek yang berjalan di server (misalnya aplikasi Node.js, Laravel, Python, dll). Fitur ini berfungsi sebagai fondasi utama bagi sistem monitoring log, di mana setiap proyek yang terdaftar akan menjadi sumber data (log source) yang dipantau secara otomatis oleh sistem.

## Komponen & Fungsionalitas
1. **Manajemen Proyek (CRUD)**:
   - Administrator dapat menambah, melihat, mengubah, dan menghapus data proyek.
   - Menyimpan informasi krusial seperti lokasi direktori absolut (*Path*) dan kategori aplikasi (*Project Type*).
2. **Log Format Builder**:
   - Memungkinkan administrator untuk mendefinisikan format log aplikasi melalui antarmuka visual.
   - Menghindari penulisan *Regular Expression* manual dengan menyediakan pilihan komponen log standar seperti `timestamp`, `level`, `message`, dan `source`.
3. **Persiapan Log Parser**:
   - Menyediakan metadata yang diperlukan bagi *engine* Log Parser untuk menemukan lokasi file log dan menentukan metode parsing yang tepat secara dinamis.

## Arsitektur Teknis (Backend)
- **Model**: `model.Project` merepresentasikan tabel di database yang menyimpan detail proyek termasuk `Path`, `ProjectType`, dan `LogFormat`.
- **Service**: `service.ProjectService` menangani validasi bisnis, termasuk pengecekan kelengkapan data dan format log.
- **Repository**: `repository.ProjectRepository` melakukan operasi CRUD ke database menggunakan GORM.
- **API Endpoints**:
  - `GET /api/projects`: Mengambil daftar seluruh proyek.
  - `POST /api/projects`: Mendaftarkan proyek baru ke sistem.
  - `PUT /api/projects/:id`: Memperbarui data proyek yang sudah ada.
  - `DELETE /api/projects/:id`: Menghapus data proyek dari sistem.

## Struktur Data (Project Model)
| Field | Tipe Data | Deskripsi |
| --- | --- | --- |
| `id` | `uint` | ID unik proyek (Primary Key). |
| `name` | `string` | Nama identitas proyek. |
| `path` | `string` | Lokasi absolut direktori proyek di sistem server. |
| `project_type`| `string` | Jenis aplikasi (contoh: Laravel, Node.js, Python, atau Custom). |
| `log_format`  | `string` | Konfigurasi pola log yang dirangkai melalui UI Builder. |
| `created_at`  | `time.Time`| Waktu pendaftaran proyek. |
| `updated_at`  | `time.Time`| Waktu pembaruan data terakhir. |

## Komponen UI (Frontend)
- **Halaman Proyek**: Menampilkan daftar proyek menggunakan tabel interaktif (`v-data-table`) dengan fitur pencarian dan aksi cepat.
- **Form Dialog (Add/Edit)**: Antarmuka input data proyek yang dilengkapi dengan *builder* visual untuk memudahkan penentuan pola log.
- **Project Store (Pinia)**: Store pusat yang mengelola sinkronisasi data proyek antara backend dan UI.

## Rencana Pengembangan (Next Steps)
Setelah manajemen CRUD selesai, fitur ini akan diintegrasikan dengan **Log Parser Engine**. Backend akan mengonversi konfigurasi `log_format` menjadi pola *Regex* yang valid untuk membaca dan mengekstraksi informasi dari file log proyek secara real-time.
