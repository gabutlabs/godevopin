# Planning: Migrasi Primary Key Log History ke UUID

## Deskripsi Singkat
Dokumen ini merencanakan implementasi perubahan tipe data *Primary Key* (`ID`) pada tabel `log_histories`. Awalnya menggunakan tipe `bigserial` (`uint`), akan diubah menjadi tipe `UUID`. 

**Alasan Best Practice:** 
Tabel log memiliki tingkat pertumbuhan data yang sangat masif (*high insert rate*). Penggunaan UUID (khususnya *gen_random_uuid()* di PostgreSQL) mencegah tabrakan ID (collision), mendukung pembuatan ID terdistribusi di level aplikasi, dan tidak akan pernah kehabisan kapasitas nilai (meskipun `bigint` sangat besar, UUID adalah standar industri untuk tabel log).

---

## Tahap 1: Pembaruan Model Backend

**1. Update Model `LogHistory`**
- Buka file `internal/model/log_history.go`.
- Ubah tipe data field `ID` dari `uint` menjadi `string`.
- Tambahkan tag GORM untuk memerintahkan PostgreSQL menggunakan UUID dan men-generate secara otomatis.
```go
type LogHistory struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    // ... field lainnya tetap sama
}
```
*Catatan:* Fungsi `gen_random_uuid()` adalah fungsi bawaan PostgreSQL untuk menghasilkan UUIDv4.

---

## Tahap 2: Strategi Migrasi Database

Mengubah tipe data *Primary Key* yang sudah memiliki data dan relasi di PostgreSQL tidak bisa dilakukan dengan `AutoMigrate` biasa. Ada dua pendekatan yang bisa diambil:

**Opsi A: Drop Table (Paling Mudah, Berlaku jika data log saat ini boleh dihapus)**
Karena ini adalah data riwayat log (yang bisa di-generate ulang oleh worker), cara paling efisien adalah:
1. Buka `cmd/cli/serve.go`.
2. Sebelum pemanggilan `AutoMigrate(...)`, tambahkan perintah sementara untuk men-drop tabel lama:
   ```go
   db.Migrator().DropTable(&model.LogHistory{})
   ```
3. Biarkan `AutoMigrate` berjalan untuk membuat ulang tabel `log_histories` dengan skema UUID yang baru.
4. *(Penting)* Hapus kembali baris `DropTable` tersebut setelah dijalankan sekali agar tabel tidak terus-menerus dihapus saat aplikasi direstart.

**Opsi B: Raw SQL Migration (Data dipertahankan)**
Jika data log mutlak harus dipertahankan, Anda harus menjalankan *Raw SQL* untuk menambahkan kolom uuid baru, memindahkan data, menghapus kolom id lama, lalu menjadikan kolom uuid sebagai primary key. (Opsi A lebih disarankan untuk kasus log ini).

---

## Tahap 3: Pembaruan Frontend (TypeScript)

**1. Update Pinia Store**
- Buka file `frontend/src/stores/project.ts`.
- Cari deklarasi `export type LogHistory`.
- Ubah tipe data `id` dari `number` menjadi `string`.
```typescript
export type LogHistory = {
  id: string; // <-- Ubah disini
  project_id: number;
  timestamp: string;
  // ...
};
```

**2. Update UI (opsional)**
- Karena tabel komponen Vuetify (`v-data-table-server`) pada `frontend/src/pages/projects/[id].vue` sudah menggunakan `item-value="id"`, tidak ada perubahan yang signifikan yang diperlukan di UI asalkan struktur JSON dari API merespon `id` sebagai string.

---

## Tahap 4: Verifikasi & Testing
1. Jalankan ulang aplikasi (`go run cmd/main.go serve`).
2. Pastikan tabel `log_histories` berhasil dibuat ulang dengan kolom `id` bertipe `uuid`.
3. Jalankan worker (`go run cmd/main.go worker`) dan biarkan worker membaca log.
4. Buka halaman Project Detail di Frontend dan pastikan daftar log muncul tanpa error, dan jika Anda memeriksa *Network Tab*, nilai `id` berupa string UUID yang panjang.

---
**Instruksi Kritis untuk Implementator**:
- Terapkan **Opsi A** (Drop Table) untuk migrasi database, karena data log bersifat *volatile* dan bisa di-parsing ulang oleh worker kapan saja.
- Jangan kerjakan (*implement*) rencana ini sebelum diberikan instruksi eksplisit dari pengguna.
