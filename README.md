# PubSub Sandbox (Kafka + Go)

Proyek ini adalah sandbox sederhana untuk eksperimen publish-subscribe menggunakan Apache Kafka dan bahasa pemrograman Go. Proyek ini terdiri dari dua bagian utama:

- **Kafka Broker & UI**: Dikelola menggunakan Docker Compose, termasuk layanan Kafka dan Kafka UI untuk monitoring.
- **Aplikasi Go**: Terdiri dari dua aplikasi, yaitu publisher dan subscriber, yang berkomunikasi dengan Kafka.

## Struktur Direktori

```
docker-compose.yml
app/
  publisher/
    go.mod
    main.go
  subscriber/
    go.mod
    main.go
data-kafka/
```

## Cara Menjalankan

1. **Jalankan Kafka dan Kafka UI**
   
   Pastikan Docker sudah terinstal. Jalankan perintah berikut di root folder:
   
   ```powershell
   docker-compose up -d
   ```
   
   Kafka akan berjalan di port 9092, dan Kafka UI di port 8080.

2. **Jalankan Publisher**
   
   Masuk ke folder `app/publisher` dan jalankan:
   
   ```powershell
   go run main.go
   ```
   
   Publisher akan mengirim beberapa pesan ke Kafka.

3. **Jalankan Subscriber**
   
   Masuk ke folder `app/subscriber` dan jalankan:
   
   ```powershell
   go run main.go
   ```
   
   Subscriber akan menerima pesan dari Kafka.

## Konfigurasi Kafka

- Konfigurasi Kafka menggunakan mode KRaft (tanpa Zookeeper).
- Dalam eksperimen kali ini, harus membuat topic kafka sesuai dengan nama topic yang didefinisikan di dalam program.
- Selain itu, kita juga harus mendefinisikan broker kafka yang kita gunakan bertindak sebagai `coordinator`.
- Data Kafka disimpan di folder `data-kafka`.
- Kafka UI dapat diakses di [http://localhost:8080](http://localhost:8080).

### Penjelasan Terkait Dead Letter Queue
> Dead Letter Queue (DLQ) adalah suatu mekanisme dalam penanganan error (error hendling) dalam Kafka

Dalam aplikasi ini, mekanisme yang diimplementasikan adalah melakukan _retry_ sebanyak 3 kali, jika sudah dilakukan sebanyak 3 kali maka message yang diterima akan dimasukkan ke dalam topic DLQ.

Untuk tujuan simulasi, maka aplikasi publisher akan mengirimkan pesan yang mengandung string 'error'. Lalu pada subscriber, akan melakukan validasi jika message yang diterima mengandung string 'error' maka aplikasi akan mengkategorikan pesan ini sebagai pesan yang gagal diproses dalam subscriber.

### Contoh Script Kafka (Jika Melakukan Konfigurasi Melalui Shell)

- Membuat Topik
```Shell
kafka-topics \
  --create \
  --topic contoh-topic \
  --bootstrap-server localhost:9092 \
  --partitions 2 \
  --replication-factor 1
```
- Membuat Topik Dead Letter Queue
```Shell
kafka-topics \
  --create \
  --topic contoh-topic-dlq \
  --bootstrap-server localhost:9092 \
  --partitions 2 \
  --replication-factor 1
```
- Menjadikan Kafka Coordinator (opsional)
```Shell
kafka-topics --create \
  --topic __consumer_offsets \
  --bootstrap-server localhost:9092 \
  --partitions 50 \
  --replication-factor 1
```

## Dependencies

- [Go](https://golang.org/) (minimal versi 1.18)
- [Docker](https://www.docker.com/)
- [Apache Kafka](https://kafka.apache.org/) (via Docker)
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go) untuk aplikasi Go

## Catatan

- Pastikan port 9092 dan 8080 tidak digunakan oleh aplikasi lain.
- Untuk eksperimen, topik Kafka yang digunakan adalah `contoh-topic` dan `contoh-topic-dlq`.

---

Silakan gunakan proyek ini untuk belajar dan bereksperimen dengan Kafka dan Go!
