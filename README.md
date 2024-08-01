# PRBCareScheduler

PRBCareScheduler adalah aplikasi scheduler yang mendukung [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi). Aplikasi ini memanfaatkan Firebase Cloud Messaging (FCM) dan Golang untuk mengirimkan data notifikasi ke pasien. Data tersebut bertujuan untuk mengingatkan pasien tentang jadwal kontrol balik dan pengambilan obat. Selain itu, aplikasi menggunakan database PostgreSQL, yang juga digunakan oleh [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi), untuk menyimpan dan mengelola data terkait jadwal kontrol balik dan pengambilan obat, serta mengelola status pembatalan hingga pengembalian stok obat.

## Fitur Utama

- **Pengiriman Data Notifikasi**: Data dikirim melalui Firebase Cloud Messaging ke pasien, berisi pengingat tentang jadwal kontrol balik dan pengambilan obat.
- **Pembatalan Otomatis**: Jika pasien tidak mengunjungi sesuai dengan jadwal yang ditetapkan, sistem akan membatalkan jadwal tersebut secara otomatis.

## Tech Stack

- **Programming Language**: Golang
- **Database**: PostgreSQL
- **Task Scheduling**: Robfig/Cron

## Environment Variables
PRBCareScheduler akan menggunakan environment variables sebagai konfigurasi utama menggantikan `config.json` jika variabel-variabel tersebut diset sebelum menjalankan proyek:
* **DB_USERNAME**: Nama pengguna database.
* **DB_PASSWORD**: Kata sandi database.
* **DB_HOST**: Host database.
* **DB_PORT**: Port koneksi database.
* **DB_NAME**: Nama database.

Cara Set Environment Variables:
- **Windows**: Gunakan System Properties > Advanced > Environment Variables, atau command setx.
- **Linux/macOS**: Tambahkan export VARIABLE="value" ke file .bashrc atau .profile dan jalankan source ~/.bashrc.