package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool" // Import driver database
)

// Model Expense sekarang memiliki UserID untuk menautkan ke pengguna
type Expense struct {
	ID          int
	UserID      int // Foreign Key ke tabel users
	Description string
	Amount      float64
	Category    string
	ExpenseDate time.Time
}

// Handler sekarang berisi koneksi ke Database (Dependency Injection)
type Handler struct {
	DB *pgxpool.Pool
}

// RenderIndexPage merender halaman lengkap dengan data dari database untuk pengguna yang login
func (h *Handler) RenderIndexPage(c *gin.Context) {
	// Ambil userID dari context yang sudah di-set oleh middleware
	userID, exists := c.Get("userID")
	if !exists {
		// Seharusnya tidak pernah terjadi jika middleware bekerja, tapi ini adalah pengaman
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna")
		return
	}

	username, usernameExists := c.Get("username")
	if !usernameExists {
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna (Username tidak ditemukan)")
		return
	}

	var expenses []Expense

	// Kueri ke database untuk mengambil pengeluaran HANYA milik pengguna ini
	rows, err := h.DB.Query(context.Background(), "SELECT id, description, amount, category, expense_date FROM expenses WHERE user_id = $1 ORDER BY expense_date DESC", userID)
	if err != nil {
		log.Printf("Query untuk mengambil expenses gagal: %v\n", err)
		c.String(http.StatusInternalServerError, "Gagal mengambil data pengeluaran")
		return
	}
	defer rows.Close()

	// Loop melalui hasil query dan scan ke dalam struct Expense
	for rows.Next() {
		var exp Expense
		// Scan juga user_id, meskipun kita tidak menampilkannya di tabel utama saat ini
		err := rows.Scan(&exp.ID, &exp.Description, &exp.Amount, &exp.Category, &exp.ExpenseDate)
		if err != nil {
			log.Printf("Gagal memindai baris data: %v\n", err)
			continue
		}
		expenses = append(expenses, exp)
	}

	// Render halaman utama dengan data yang sudah difilter
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Expenses": expenses,
		"Username": username,
	})
}

// AddExpense menangani form submission, menyimpan data ke database, dan mengembalikan daftar baru
func (h *Handler) AddExpense(c *gin.Context) {
	// Ambil userID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna")
		return
	}

	// Parse data dari form
	description := c.PostForm("description")
	amountStr := c.PostForm("amount")
	category := c.PostForm("category")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Println("Error parsing amount:", err)
		c.String(http.StatusBadRequest, "Invalid amount")
		return
	}

	// Kueri INSERT ke database, menyertakan user_id
	// Menggunakan parameter ($1, $2, dll.) untuk mencegah SQL Injection
	_, err = h.DB.Exec(context.Background(), "INSERT INTO expenses (description, amount, category, user_id) VALUES ($1, $2, $3, $4)", description, amount, category, userID)
	if err != nil {
		log.Printf("Gagal memasukkan data expense: %v\n", err)
		c.String(http.StatusInternalServerError, "Gagal menyimpan pengeluaran")
		return
	}

	// Setelah berhasil menyimpan, ambil kembali daftar pengeluaran yang sudah update
	var expenses []Expense
	rows, err := h.DB.Query(context.Background(), "SELECT id, description, amount, category, expense_date FROM expenses WHERE user_id = $1 ORDER BY expense_date DESC", userID)
	if err != nil {
		log.Printf("Query setelah insert gagal: %v\n", err)
		// Kirimkan saja daftar kosong jika gagal mengambil data baru
		c.HTML(http.StatusOK, "_expense-list.html", gin.H{"Expenses": []Expense{}})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exp Expense
		err := rows.Scan(&exp.ID, &exp.Description, &exp.Amount, &exp.Category, &exp.ExpenseDate)
		if err != nil {
			log.Printf("Gagal memindai baris data setelah insert: %v\n", err)
			continue
		}
		expenses = append(expenses, exp)
	}

	// Render HANYA partial `_expense-list.html` untuk dikirim kembali ke HTMX
	c.HTML(http.StatusOK, "_expense-list.html", gin.H{
		"Expenses": expenses,
	})
}