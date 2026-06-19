# Property Marketplace REST API

Sebuah backend service untuk ekosistem marketplace properti, dirancang untuk memfasilitasi transaksi real estat (pencarian, pemasangan, promosi, dan pemesanan properti). Mengambil rujukan standar industri dari platform seperti OLX dan Rumah123, API ini dibangun menggunakan **Golang** dengan implementasi **Layered Architecture** untuk memastikan skalabilitas, kemudahan pengujian, dan pemeliharaan sistem dalam jangka panjang. Seluruh interaksi endpoint didemonstrasikan dan diuji melalui Postman.

## Tentang Dokumen Ini

README ini bertindak sebagai blueprint teknis dan fungsional dari Property Marketplace API. Dokumen ini memetakan secara sistematis cakupan fitur, spesifikasi endpoint, desain model data, alur transaksi (payment flow), serta metodologi pengujian yang diterapkan. Tujuannya adalah memberikan konteks yang transparan bagi para penguji, kolaborator, maupun maintainer lanjutan mengenai arsitektur sistem dan rasionalisasi di balik setiap keputusan desain (Design Rationale). Deployment dokumen ini pada https://p2-individual-project-production.up.railway.app dan tautan export JSON https://web.postman.co/workspace/Kevin-Yusuf-Briliantama's-Works~982dc8fe-817c-4ce5-8663-932f52adcd8c/collection/55382287-46cf81ae-4f4f-4f86-a581-e4472f6ab820

## 1. Ringkasan Eksekutif Produk

Platform ini adalah marketplace properti yang mempertemukan penjual yang ingin mendaftarkan propertinya dengan pencari yang sedang mencari hunian, baik untuk dibeli maupun disewa. Penjual memasang listing, dan pencari menelusurinya menggunakan filter realistis yang menyerupai standar OLX. Nilai inti sistem ini adalah menciptakan titik temu yang efisien antara listing yang baik dan orang yang tepat.

Agar platform tidak berhenti sebagai katalog statis, melainkan terasa hidup dan tepercaya, API ini dilengkapi tiga fitur pendukung yang saling melengkapi: verifikasi email untuk menjaga keamanan dan keaslian akun, mekanisme promote untuk meningkatkan visibilitas listing, dan sistem booking untuk memfasilitasi transaksi yang nyata.

---

## 2. Ekosistem dan Peran Pengguna

Sistem dibangun di atas tiga peran yang saling berinteraksi. Masing masing memiliki kapabilitas dan nilai yang jelas dalam alur produk.

| Peran       | Kapabilitas dan Nilai                                                                                                                                                                                           |
| :---------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Penjual** | Mendaftar dengan akun terverifikasi, mengelola listing propertinya, dan dapat menggunakan fitur promote berbayar untuk menaikkan visibilitas listing agar lebih cepat ditemukan pembeli.                        |
| **Pencari** | Menelusuri katalog tanpa perlu login, menyaring properti dengan spesifikasi realistis, dan ketika menemukan properti yang tepat, melanjutkan ke pemesanan melalui sistem booking yang terintegrasi.             |
| **Sistem**  | Menjaga integritas data melalui verifikasi pengguna, mengatur urutan tampilan agar listing yang dipromosikan muncul lebih dulu, dan mengelola perubahan status data selama proses pembayaran secara transparan. |

---

## 3. Cakupan dan integrasi

API ini menyediakan dua belas hingga empat belas endpoint yang mencakup autentikasi, pengelolaan listing, pencarian berfilter, serta tiga alur pembayaran. Terdapat dua integrasi layanan pihak ketiga:

1. Xendit, untuk membuat tautan pembayaran (payment link) pada promote, booking, dan pelunasan.
2. Mailjet, untuk verifikasi email saat pendaftaran.

Pengujian mencakup minimal satu unit test pada lapisan usecase menggunakan mock.

---

## 4. Keputusan rancangan

Tiap keputusan dicatat beserta alasannya agar mudah dijelaskan dan tidak berubah ubah.

1. Payment gateway: Xendit. Pola integrasinya sudah dikuasai dan menghasilkan invoice link yang dapat dibagikan ke pembayar.
2. Tiga alur pembayaran, semuanya berbagi mesin yang sama (membuat invoice link Xendit lalu menunggu konfirmasi webhook), dibedakan oleh jenis transaksi:
   1. promote, penjual menaikkan listing menjadi featured.
   2. booking, pembeli atau penyewa memesan properti dengan tanda jadi.
   3. pelunasan, pembeli atau penyewa melunasi pemesanan.
3. Bentuk pembayaran: tautan (invoice link). Setiap pembayaran menghasilkan tautan dari Xendit yang dibuka oleh pembayar, mengikuti pola yang sama dengan contoh invoice Xendit.
4. Email: Mailjet, untuk verifikasi akun saat register.
5. Auth: JWT. Pemilik resource selalu diambil dari token, bukan dari body request.
6. Foto properti: disimpan sebagai URL, bukan unggahan berkas.
7. Lokasi: disimpan datar sebagai kota dan kecamatan, bukan tabel hierarki provinsi.
8. Invoice: dibuat oleh gateway. Xendit yang menerbitkan invoice dan tautan pembayaran, sistem cukup menyimpan referensinya (id dan url) pada catatan transaksi.
9. Error contract: bentuk respons konsisten di seluruh endpoint, dengan penanda sukses, pesan, dan daftar kesalahan per field opsional.
10. Sumber kebenaran status lunas: selalu webhook Xendit, tidak pernah klaim dari klien.
11. Pengujian: menggunakan mock pada repository dan client Xendit agar usecase dapat diuji secara terisolasi.

---

## 5. Filterasi

Seluruh filter pada dasarnya adalah cara mengklasifikasi data. Perbedaannya terletak pada cara menyaring: nilai berupa pilihan (huruf) atau nilai berupa angka.

Filter pilih dari daftar:

1. Tipe properti: rumah, apartemen, tanah.
2. Tipe transaksi: dijual atau disewakan.
3. Kota dan kecamatan.
4. Jumlah kamar tidur.

Filter saring rentang:

1. Harga, disaring dengan nilai minimum dan maksimum.
2. Luas, disaring dengan nilai minimum dan maksimum.

Rentang harga dibedakan antara jual dan sewa karena skalanya jauh berbeda.

Untuk dijual:

1. di bawah 500 juta
2. 500 juta sampai 1 miliar
3. 1 sampai 2 miliar
4. 2 sampai 5 miliar
5. di atas 5 miliar

Untuk disewakan, per bulan:

1. di bawah 1 juta
2. 1 sampai 3 juta
3. 3 sampai 7 juta
4. 7 sampai 15 juta
5. di atas 15 juta

Harga disimpan sebagai angka mentah, lalu disaring dengan minimum dan maksimum. Daftar rentang di atas berfungsi sebagai tombol cepat yang mengeset minimum dan maksimum secara otomatis.

Mengenai luas, penggunaannya menyesuaikan tipe properti:

1. Rumah dan tanah memakai luas tanah.
2. Apartemen memakai luas bangunan.

Dengan begitu satu kolom luas yang relevan ditampilkan dan disaring sesuai tipe propertinya, tanpa membebani pengguna dengan dua rentang sekaligus.

Catatan tentang kategori OLX yang tampil tergabung, seperti dijual rumah dan apartemen: itu sebenarnya gabungan dua field, yaitu tipe transaksi dan tipe properti. Keduanya disimpan terpisah, dan penggabungan label hanya urusan tampilan.

---

## 6. Daftar endpoint

| #   | Method | Path                    | Auth         | Keterangan                                              |
| --- | ------ | ----------------------- | ------------ | ------------------------------------------------------- |
| 1   | POST   | /users/register         | publik       | buat akun, kirim email verifikasi (Mailjet)             |
| 2   | POST   | /users/login            | publik       | kembalikan JWT                                          |
| 3   | GET    | /users/verify           | publik       | validasi token dari email                               |
| 4   | POST   | /properties             | JWT          | pasang listing                                          |
| 5   | GET    | /properties             | publik       | daftar dengan filter, pencarian, rentang harga dan luas |
| 6   | GET    | /properties/:id         | publik       | detail satu listing                                     |
| 7   | PUT    | /properties/:id         | JWT, pemilik | ubah listing milik sendiri                              |
| 8   | DELETE | /properties/:id         | JWT, pemilik | hapus listing milik sendiri                             |
| 9   | GET    | /my/properties          | JWT          | lihat listing milik sendiri                             |
| 10  | POST   | /properties/:id/promote | JWT, pemilik | buat invoice link untuk featured                        |
| 11  | POST   | /properties/:id/book    | JWT          | pesan properti dengan tanda jadi, buat invoice link     |
| 12  | POST   | /bookings/:id/settle    | JWT          | lunasi pemesanan, buat invoice link                     |
| 13  | POST   | /payments/webhook       | signature    | Xendit konfirmasi lunas, jalankan efek transaksi        |
| 14  | GET    | /transactions           | JWT          | riwayat dan status pembayaran sendiri                   |

---

## 7. Model data

Tiga tabel inti.

users: identitas, email, kata sandi ter-hash, dan penanda terverifikasi.

properties: dimiliki satu user, berisi

1. judul listing
2. tipe properti: rumah, apartemen, atau tanah
3. tipe transaksi: dijual atau disewakan
4. harga
5. luas tanah, dipakai untuk rumah dan tanah
6. luas bangunan, dipakai untuk apartemen
7. jumlah kamar tidur
8. jumlah kamar mandi
9. sertifikat: SHM atau HGB
10. lokasi: kota dan kecamatan
11. deskripsi
12. daftar URL gambar
13. status: available, booked, atau sold
14. waktu kedaluwarsa featured, kosong bila tidak sedang dipromosikan

transactions: dimiliki satu user dan menunjuk satu properti, berisi jenis transaksi (promote, booking, atau pelunasan), nominal, status (pending, paid, atau expired), serta referensi Xendit berupa external id, invoice id, dan invoice url.

Relasi: satu user memiliki banyak properti, dan satu properti memiliki banyak transaksi.

Catatan tentang penanda featured: cukup gunakan satu kolom waktu kedaluwarsa yang boleh kosong. Sebuah listing dianggap featured selama waktu kedaluwarsanya masih di masa depan, dan otomatis berhenti featured begitu waktu itu lewat. Pendekatan satu kolom ini lebih sederhana daripada memakai penanda boolean terpisah, karena status aktif dapat dihitung langsung dari waktunya.

---

## 8. Alur pembayaran

Ketiga alur memakai mesin yang sama: sistem membuat catatan transaksi berstatus pending, meminta invoice link ke Xendit, mengembalikan tautan ke pengguna, lalu menunggu webhook. Status lunas hanya dipercaya dari webhook, tidak pernah dari klaim klien.

Promote, oleh penjual:

1. Penjual menekan promote pada listing miliknya.
2. Sistem membuat transaksi promote pending dan invoice link.
3. Setelah webhook mengonfirmasi lunas, listing diberi waktu kedaluwarsa featured.
4. Pada daftar properti, listing featured yang masih aktif diurutkan paling atas.

Durasi featured disetel singkat, sekitar satu sampai dua menit, agar dapat diperagakan langsung saat demo: listing naik ke atas saat dibayar, lalu turun kembali setelah masa featured habis.

Booking lalu pelunasan, oleh pembeli atau penyewa:

1. Pencari yang serius menekan booking pada properti.
2. Sistem membuat transaksi booking pending dan invoice link untuk tanda jadi.
3. Setelah webhook mengonfirmasi lunas, status properti menjadi booked.
4. Pembeli kemudian menekan pelunasan untuk membayar sisa.
5. Setelah webhook mengonfirmasi lunas, status properti menjadi sold.

---

## 9. Pengujian

Fokus pada lapisan usecase dengan mock agar terisolasi. Repository dan client Xendit di-mock. Skenario contoh:

1. Promote atau ubah listing yang bukan milik sendiri harus gagal dengan kesalahan kepemilikan.
2. Booking properti yang sudah berstatus booked harus ditolak.
3. Saat valid, usecase membuat transaksi pending dan memanggil client pembayaran satu kali.

---

## 10. Rencana pengerjaan

Diurutkan dari bagian paling berisiko ke paling ringan.

1. Fondasi: penyiapan proyek, koneksi basis data, migrasi, model, struktur lapisan, middleware JWT, dan error contract.
2. Inti: autentikasi lengkap dengan verifikasi email, lalu CRUD properti dengan pemeriksaan kepemilikan dan daftar berfilter.
3. Penerapan awal: deploy ke Railway begitu kerangka berjalan, karena webhook membutuhkan URL publik.
4. Pembayaran: promote, booking, pelunasan, dan webhook Xendit.
5. Pengujian: unit test usecase dengan mock.
6. Pelengkap: merapikan koleksi Postman, melengkapi dokumentasi, dan menyiapkan materi presentasi.

---

## 11. Risiko dan mitigasi

1. Webhook membutuhkan URL publik. Mitigasi: deploy lebih awal, atau gunakan ngrok saat demo.
2. Keaslian webhook. Verifikasi signature wajib agar notifikasi lunas tidak dapat dipalsukan.
3. Kerahasiaan kunci API. Kunci Xendit dan Mailjet disimpan pada environment, tidak pernah masuk ke repositori.
4. Basis data kosong saat demo. Sediakan data contoh, termasuk listing featured, agar pengurutan terlihat.
5. Cakupan melebar. Daftar fitur dan batasan sudah ditetapkan agar pengerjaan tetap terarah.

---

## 12. Di luar cakupan

Hal hal berikut sengaja tidak dikerjakan agar fokus terjaga:

1. Pembuatan deskripsi otomatis dengan AI.
2. Unggahan berkas gambar, cukup URL.
3. Langganan agen.
4. Tabel hierarki lokasi provinsi.
5. Favorit, wishlist, dan percakapan antar pengguna.

---

## 13. Daftar deliverable

1. API berjalan dan ter-deploy.
2. Seluruh endpoint berfungsi dan teruji di Postman.
3. Koleksi Postman rapi per endpoint.
4. Unit test usecase lulus.
5. Data contoh termasuk listing featured.
6. Dokumentasi penyiapan dan cara menjalankan.
7. Materi presentasi.
