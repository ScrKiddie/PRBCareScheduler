# PRBCareScheduler

<p align="center">
<img src="https://github.com/user-attachments/assets/50eea6b6-e922-4dda-a036-3fbf1704458d" alt="prbcare" width="400">
</p>

PRBCareScheduler adalah aplikasi scheduler yang mendukung [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi). Aplikasi ini memanfaatkan Firebase Cloud Messaging (FCM) dan Golang untuk mengirimkan data notifikasi ke pasien. Data tersebut bertujuan untuk mengingatkan pasien tentang jadwal kontrol balik dan pengambilan obat. Aplikasi ini menggunakan database yang juga digunakan oleh [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi), untuk menyimpan dan mengelola data terkait jadwal kontrol balik dan pengambilan obat, serta mengelola status pembatalan hingga pengembalian stok obat.

## Fitur Utama

- **Pengiriman Data Notifikasi**: Data dikirim melalui Firebase Cloud Messaging ke pasien, berisi pengingat tentang jadwal kontrol balik dan pengambilan obat.
- **Pembatalan Otomatis**: Jika pasien tidak mengunjungi sesuai dengan jadwal yang ditetapkan, sistem akan membatalkan jadwal tersebut secara otomatis.

## Tech Stack

- **Programming Language**: Golang
- **Task Scheduling**: Robfig/Cron
- **Database**: PostgreSQL


## Environment Variables
PRBCareScheduler akan menggunakan environment variables sebagai konfigurasi utama menggantikan `config.json` jika variabel-variabel tersebut diset sebelum menjalankan proyek:

| Key                                                                | Type    | Deskripsi                           | Contoh                                                       |
|--------------------------------------------------------------------|---------|-------------------------------------|--------------------------------------------------------------|
| **TIME_NOTIFY_KONTROL,TIME_NOTIFY_PROLANIS,TIME_NOTIFY_OBAT** | `string`| Waktu notifikasi dalam format cron. | `* * * * *`                                                       |
| **TIME_CANCEL**                                                    | `string`| Waktu pembatalan dalam format cron. | `* * * * *`                                                       |
| **DB_USERNAME**                                                    | `string`| Nama pengguna database.             | `root`                                                       |
| **DB_PASSWORD**                                                    | `string`| Kata sandi database.                | `password123`                                                |
| **DB_HOST**                                                        | `string`| Host database.                      | `localhost`                                                  |
| **DB_PORT**                                                        | `int`   | Port koneksi database.              | `3306`                                                       |
| **DB_NAME**                                                        | `string`| Nama database.                      | `prbcare`                                                    |

Cara Set Environment Variables:
- **Windows**: Gunakan System Properties > Advanced > Environment Variables, atau command setx.
- **Linux/macOS**: Tambahkan export VARIABLE="value" ke file .bashrc atau .profile dan jalankan source ~/.bashrc.