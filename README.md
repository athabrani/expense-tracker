# expense-tracker
expense tracker using GO + HTMX for DevSecOps Final Project

# Expense Tracker (Go + HTMX)

Expense Tracker adalah aplikasi web sederhana yang dibangun menggunakan Go untuk backend dan HTMX untuk interaktivitas frontend. Aplikasi ini memungkinkan pengguna untuk mendaftar, login, mencatat pengeluaran harian mereka, dan melihat riwayat pengeluaran mereka. Proyek ini juga mendemonstrasikan beberapa praktik dasar DevSecOps.

## Fitur Utama

* Registrasi Pengguna
* Login & Logout Pengguna
* Sistem Sesi berbasis Cookie
* Tambah Data Pengeluaran (Deskripsi, Jumlah, Kategori)
* Lihat Riwayat Pengeluaran (per pengguna)
* Interaksi Dinamis menggunakan HTMX tanpa reload halaman penuh
* Penyimpanan data persisten menggunakan PostgreSQL
* HTTPS untuk development lokal (opsional, menggunakan mkcert)

## Tech Stack

* **Backend:** Go (Golang)
    * Framework Web: Gin (`github.com/gin-gonic/gin`)
    * Driver Database: pgx (`github.com/jackc/pgx/v5/pgxpool`)
    * Manajemen Lingkungan: godotenv (`github.com/joho/godotenv`)
    * Password Hashing: bcrypt (`golang.org/x/crypto/bcrypt`)
* **Frontend:**
    * HTMX (`htmx.org`)
    * HTML5 (Go Templates)
    * CSS3 (Bootstrap 5 untuk styling dasar)
* **Database:** PostgreSQL (dijalankan melalui Docker)
* **Alat Keamanan (DevSecOps):**
    * SAST: `gosec`
    * SCA: `govulncheck`
    * DAST: OWASP ZAP (manual scan)
* **Development Lokal HTTPS:** `mkcert` (opsional)

## Prasyarat

Sebelum memulai, pastikan Anda sudah menginstal perangkat lunak berikut di komputer Anda:

1.  **Go**: Versi 1.21 atau lebih baru. (Cek dengan `go version`)
2.  **Docker**: Untuk menjalankan database PostgreSQL. (Cek dengan `docker --version`)
3.  **Docker Compose**: Biasanya terinstal bersama Docker. (Cek dengan `docker-compose --version`)
4.  **Git**: Untuk mengkloning repositori. (Cek dengan `git --version`)
5.  **(Opsional) `mkcert`**: Jika Anda ingin menjalankan aplikasi dengan HTTPS di lingkungan lokal. Ikuti panduan instalasi di [halaman mkcert](https://github.com/FiloSottile/mkcert#installation).

## Setup dan Instalasi

Ikuti langkah-langkah berikut untuk menyiapkan dan menjalankan proyek ini di komputer Anda:

1.  **Kloning Repositori:**
    ```bash
    git clone <URL-repositori-GitHub-Anda>
    cd expense-tracker 
    ```
    *(Ganti `<URL-repositori-GitHub-Anda>` dengan URL repositori Anda sebenarnya)*

2.  **Siapkan File Environment (`.env`):**
    File ini akan menyimpan kredensial database dan konfigurasi lainnya.
    * Salin file contoh:
        ```bash
        cp .env.example .env
        ```
    * Buka file `.env` yang baru dibuat dan isi nilainya sesuai konfigurasi Anda. Contohnya:
        ```env
        DB_USER=user_expenses
        DB_PASSWORD=password_expenses
        DB_HOST=localhost
        DB_PORT=5432
        DB_NAME=db_expenses

        JWT_SECRET=passwordJWT
        COOKIE_DOMAIN=domain_web_kamu
        ```

3.  **Nyalakan Database (Docker):**
    Dari folder root proyek, jalankan perintah berikut untuk membuat dan menyalakan container database PostgreSQL:
    ```bash
    docker-compose up -d
    ```


4.  **Buat Tabel Database:**
    * Hubungkan ke database PostgreSQL Anda menggunakan DBeaver atau alat database lainnya dengan detail koneksi dari file `.env` Anda.
    * Jalankan skrip SQL berikut untuk membuat tabel `users` dan `expenses`:
        ```sql
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            username VARCHAR(50) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );

        CREATE TABLE expenses (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            description VARCHAR(255) NOT NULL,
            amount NUMERIC(10, 2) NOT NULL,
            category VARCHAR(100),
            expense_date TIMESTAMPTZ DEFAULT NOW(),
            
            CONSTRAINT fk_user
                FOREIGN KEY (user_id)
                REFERENCES users(id)
                ON DELETE CASCADE
        );
        ```

5.  **Instal Dependensi Go:**
    Perintah ini akan mengunduh semua paket Go yang dibutuhkan oleh proyek (seperti Gin, pgx, dll.) berdasarkan file `go.mod`.
    ```bash
    go mod tidy
    ```

6.  **(Opsional) Setup Local HTTPS dengan `mkcert`:**
    Jika Anda ingin menjalankan aplikasi dengan HTTPS di `localhost`:
    * Pastikan Anda sudah menginstal `mkcert`.
    * Instal CA lokal (hanya perlu sekali):
        ```bash
        # Untuk Windows (jika mkcert.exe di folder proyek)
        .\mkcert -install 
        ```
        
    * Buat sertifikat untuk `localhost` di folder proyek Anda:
        ```bash
        # Untuk Windows
        .\mkcert localhost
        ```
        Ini akan menghasilkan file `localhost.pem` dan `localhost-key.pem`.

## Menjalankan Aplikasi

1.  Pastikan database Anda sudah berjalan (`docker-compose up -d`).
2.  Jalankan server Go dari folder root proyek:
    ```bash
    go run main.go
    ```
3.  Aplikasi akan berjalan di:
    * **HTTP:** `http://localhost:8080`
    * **HTTPS (jika Anda setup mkcert dan mengubah `main.go`):** `https://localhost:8443` (atau port lain yang Anda tentukan).

    Buka URL tersebut di browser Anda.

## Menjalankan Pengecekan Keamanan (DevSecOps)

1.  **SAST (Static Application Security Testing) dengan `gosec`:**
    * Instal (jika belum): `go install github.com/securego/gosec/v2/cmd/gosec@latest`
    * Jalankan dari root proyek: `gosec ./...`

2.  **SCA (Software Composition Analysis) dengan `govulncheck`:**
    * Instal (jika belum): `go install golang.org/x/vuln/cmd/govulncheck@latest`
    * Jalankan dari root proyek: `govulncheck ./...`

3.  **DAST (Dynamic Application Security Testing) dengan OWASP ZAP:**
    * Pastikan aplikasi Go dan database Anda sedang berjalan.
