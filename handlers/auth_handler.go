package handlers

import (
	"context"
	"log"
	"errors"
	"net/http"
	"strconv"


	"github.com/jackc/pgconn"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Model untuk User
type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
}

// Menampilkan halaman register
func (h *Handler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// Memproses form register
func (h *Handler) Register(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal mengenkripsi password")
		return
	}

	_, err = h.DB.Exec(context.Background(), "INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)", email, username, string(hashedPassword))
	if err != nil {
		var pgErr *pgconn.PgError
		// Cek apakah error ini adalah error spesifik dari PostgreSQL
		if errors.As(err, &pgErr) {
			// Kode 23505 adalah untuk 'unique_violation'
			if pgErr.Code == "23505" {
				// Cek constraint mana yang dilanggar untuk pesan yang lebih baik
				if pgErr.ConstraintName == "users_username_key" {
					c.Redirect(http.StatusFound, "/register?error=username_taken")
					return
				}
				if pgErr.ConstraintName == "users_email_key" {
					c.Redirect(http.StatusFound, "/register?error=email_taken")
					return
				}
			}
		}
		// Untuk error database lainnya
		log.Printf("Gagal menyimpan user: %v\n", err)
		c.String(http.StatusInternalServerError, "Terjadi kesalahan internal saat registrasi")
		return
	}

	c.Redirect(http.StatusFound, "/login?status=reg_success")
}

// Menampilkan halaman login
func (h *Handler) LoginPage(c *gin.Context) {
	errorType := c.Query("error")
	statusType := c.Query("status") // <-- TAMBAHKAN INI

	// Kirimkan keduanya ke template
	c.HTML(http.StatusOK, "login.html", gin.H{
		"ErrorType": errorType,
        "StatusType": statusType, // <-- TAMBAHKAN INI
	})
}

// Memproses form login
func (h *Handler) Login(c *gin.Context) {
    email := c.PostForm("email")
    password := c.PostForm("password")

    var user User
    // Perbaiki urutan Scan di sini
    err := h.DB.QueryRow(context.Background(), "SELECT id, email, username, password_hash FROM users WHERE email = $1", email).Scan(
        &user.ID,
        &user.Email,      
        &user.Username,   
        &user.PasswordHash,
    )
    if err != nil {
        c.Redirect(http.StatusFound, "/login?error=true")
        return
    }


	// Bandingkan password yang di-submit dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// JIKA GAGAL: Redirect kembali ke /login dengan penanda error
		c.Redirect(http.StatusFound, "/login?error=true")
		return
	}

	
	c.SetCookie("username", user.Username, 3600*24, "/", "localhost", false, true)
	c.SetCookie("user_id", strconv.Itoa(user.ID), 3600*24, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/")
}

// Logout
func (h *Handler) Logout(c *gin.Context) {
    // Hapus cookie dengan mengatur max_age ke -1
    c.SetCookie("user_id", "", -1, "/", "localhost", false, true)
	c.SetCookie("username", "", -1, "/", "localhost", false, true)
    c.Redirect(http.StatusFound, "/login")
}