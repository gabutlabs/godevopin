# Worker Services Feature

## Deskripsi Singkat
Fitur Worker Services berfungsi sebagai control panel untuk mengatur background task dan layanan pekerja (worker) bawaan di dalam Devopin. Fitur ini memungkinkan admin untuk membuat, memonitor, dan mengendalikan siklus hidup daemon/worker process.

## Komponen & Fungsionalitas
1. **Daftar Worker Service**:
   - Menampilkan tabel mendetail dari setiap worker service.
   - Kolom informasi meliputi: Name, Description, Desired State, Current Status, Health Status, PID, Last Heartbeat, Last Success, dan timestamp pembuatan/pembaruan.
2. **Siklus Hidup (Lifecycle) Actions**:
   - **Start**: Menjalankan worker yang berstatus stopped.
   - **Stop**: Menghentikan worker yang berstatus running.
   - **Restart**: Memulai ulang worker.
   - Action akan men-disable (abu-abu) secara pintar tergantung dari state saat ini (contoh: tidak bisa start jika status sudah running).
3. **Status Indicators (Chips)**:
   - **Desired State**: `enabled` (hijau), `disabled` (kuning/warning).
   - **Current Status**: `running` (hijau), `stopped` (abu-abu), `failed` (merah), `starting` (biru), `degraded` (kuning).
   - **Health Status**: `healthy` (hijau), `unhealthy` (merah).
4. **Manajemen Data Worker (CRUD)**:
   - Form dialog untuk membuat (`Add Worker Service`) dan mengedit (`Edit`) konfigurasi worker. Validasi disertakan menggunakan Vuelidate (Name, Description, Desired State).
   - Fitur penghapusan (`Delete`) worker service.

## Backend Integration
- Frontend berkomunikasi dengan state manager Pinia (`useWorkerServiceStore`) yang mengirimkan REST/Socket request ke backend Go.
- Backend menggunakan Cobra CLI (`worker` command) untuk meluncurkan background process secara riil, lalu mencatat heartbeat dan PIDs di dalam sistem.
