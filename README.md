# PRBCareScheduler

PRBCareScheduler adalah aplikasi scheduler yang mendukung [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi). Aplikasi ini memanfaatkan Firebase Cloud Messaging (FCM) dan Golang untuk mengirimkan data notifikasi ke pasien. Data tersebut bertujuan untuk mengingatkan pasien tentang jadwal kontrol balik dan pengambilan obat. Selain itu, aplikasi menggunakan database PostgreSQL, yang juga digunakan oleh [PRBCareAPI](https://github.com/scrkiddie/PRBCareApi), untuk menyimpan dan mengelola data terkait jadwal kontrol balik dan pengambilan obat, serta mengelola status pembatalan hingga pengembalian stok obat.

## Fitur Utama

- **Pengiriman Data Notifikasi**: Data dikirim melalui Firebase Cloud Messaging ke pasien, berisi pengingat tentang jadwal kontrol balik dan pengambilan obat.
- **Pembatalan Otomatis**: Jika pasien tidak mengunjungi sesuai dengan jadwal yang ditetapkan, sistem akan membatalkan jadwal tersebut secara otomatis.
