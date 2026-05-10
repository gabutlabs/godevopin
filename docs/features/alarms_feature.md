# Alarms Feature

## Deskripsi Singkat
Fitur Alarms digunakan untuk memonitor dan mengelola peringatan (alerts) sistem yang dipicu oleh berbagai kondisi atau ambang batas (thresholds) monitoring yang telah ditentukan sebelumnya.

## Komponen & Fungsionalitas
1. **Manajemen Status Alarm**:
   - Fitur ini dibagi menjadi dua tab utama: **Active** (Alarm yang saat ini sedang aktif) dan **History** (Riwayat alarm yang sudah berlalu).
2. **Filter & Pencarian**:
   - Menyediakan dialog filter untuk menyaring alarm berdasarkan statusnya.
   - Opsi status utama: `FIRING` (sedang terjadi) dan `ACKNOWLEDGED` (telah diakui/ditangani).
3. **Komponen UI**:
   - `ActiveAlarmComponent`: Untuk merender daftar alarm aktif.
   - `HistoryAlarm`: Untuk merender riwayat alarm.

## Backend & Database
- Backend terus mengevaluasi metrik (seperti CPU dan Memory) terhadap threshold yang telah dikonfigurasi.
- Jika threshold terlampaui, backend akan menghasilkan record alarm yang dikirim ke antarmuka ini.
