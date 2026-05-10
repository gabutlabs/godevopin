# Users Management Feature

## Deskripsi Singkat
Fitur Users Management adalah modul otentikasi dan otorisasi dasar yang memungkinkan pengelolaan data pengguna yang memiliki akses ke dashboard Devopin.

## Komponen & Fungsionalitas
1. **Tabel Pengguna**:
   - Menampilkan daftar semua pengguna terdaftar dengan kolom: Name, Email, Created At, dan Updated At.
   - Fitur pencarian (`Search name...`) untuk menyaring daftar pengguna secara instan.
2. **Form Manipulasi Pengguna (CRUD)**:
   - **Add User**: Dialog untuk menambahkan pengguna baru. Membutuhkan input `Name` dan `Email`. Password secara default akan diatur ke `Password123!` (yang ditandai di dalam UI).
   - **Edit User**: Memperbarui informasi `Name` dan `Email` dari pengguna yang sudah ada.
   - **Delete User**: Menghapus pengguna secara permanen dari sistem.
3. **Validasi Form**:
   - Menggunakan Vuelidate untuk memvalidasi kolom input agar `Name` wajib diisi, dan format `Email` harus valid sebelum disubmit ke server.

## Backend & Database
- Disimpan di tabel database (model `User` di `internal/model/user.go`).
- Kata sandi (password) yang masuk ke database akan di-hash menggunakan bcrypt oleh Go Backend.
- Manajemen sesi menggunakan JWT Token yang diinisialisasi melalui login.
