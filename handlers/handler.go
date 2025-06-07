package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)


// Handler untuk menampung semua dependensi (DB, JWT Secret)
type Handler struct {
	DB        *pgxpool.Pool
	JWTSecret string
	CookieDomain string
}

// Claims untuk data di dalam JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Model untuk User
type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
}

// Model untuk Expense
type Expense struct {
	ID          int
	UserID      int
	Description string
	Amount      float64
	Category    string
	ExpenseDate time.Time
}

//  CONSTRUCTOR UNTUK HANDLER

func NewHandler(db *pgxpool.Pool, jwtSecret string) *Handler {
	return &Handler{
		DB:        db,
		JWTSecret: jwtSecret,
		CookieDomain: cookieDomain,
	}
}


// HANDLER UNTUK AUTHENNTIKASI 

func (h *Handler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

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
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
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
		log.Printf("Gagal menyimpan user: %v\n", err)
		c.String(http.StatusInternalServerError, "Terjadi kesalahan internal saat registrasi")
		return
	}

	c.Redirect(http.StatusFound, "/login?status=reg_success")
}

func (h *Handler) LoginPage(c *gin.Context) {
	errorType := c.Query("error")
	statusType := c.Query("status")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"ErrorType":  errorType,
		"StatusType": statusType,
	})
}

func (h *Handler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user User
	err := h.DB.QueryRow(context.Background(), "SELECT id, email, username, password_hash FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
    if err != nil {
        if err == pgx.ErrNoRows {
            _ = bcrypt.CompareHashAndPassword([]byte("$2a$10$fakesaltandhashforsecurity."), []byte(password))
            c.Redirect(http.StatusFound, "/login?error=true")
            return
        }
        c.String(http.StatusInternalServerError, "Terjadi kesalahan internal")
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        c.Redirect(http.StatusFound, "/login?error=true")
		return
	}

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   strconv.Itoa(user.ID),
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(h.JWTSecret))
    if err != nil {
        c.String(http.StatusInternalServerError, "Gagal membuat token otentikasi")
        return
    }
    c.SetCookie("token", tokenString, 3600*24, "/", h.CookieDomain, false, true)
    c.Redirect(http.StatusFound, "/")
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", h.CookieDomain, false, true)
	c.Redirect(http.StatusFound, "/login")
}

// MIDDLEWARE

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=unauthorized")
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("METODE SIGNIN ANEH")
			}
			return []byte(h.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login?error=invalid_token")
			c.Abort()
			return
		}
		
		c.Set("userID", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}


// HANDLER UNTUK EXPENSE

func (h *Handler) RenderIndexPage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna")
		return
	}

	username, usernameExists := c.Get("username")
	if !usernameExists {
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna (Username tidak ditemukan)")
		return
	}

	var expenses []Expense
	rows, err := h.DB.Query(context.Background(), "SELECT id, user_id, description, amount, category, expense_date FROM expenses WHERE user_id = $1 ORDER BY expense_date DESC", userID)
	if err != nil {
		log.Printf("Query untuk mengambil expenses gagal: %v\n", err)
		c.String(http.StatusInternalServerError, "Gagal mengambil data pengeluaran")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exp Expense
		err := rows.Scan(&exp.ID, &exp.UserID, &exp.Description, &exp.Amount, &exp.Category, &exp.ExpenseDate)
		if err != nil {
			log.Printf("Gagal memindai baris data: %v\n", err)
			continue
		}
		expenses = append(expenses, exp)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Expenses": expenses,
		"Username": username,
	})
}

func (h *Handler) AddExpense(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Gagal mendapatkan informasi pengguna")
		return
	}

	description := c.PostForm("description")
	amountStr := c.PostForm("amount")
	category := c.PostForm("category")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Println("Error parsing amount:", err)
		c.String(http.StatusBadRequest, "Invalid amount")
		return
	}

	_, err = h.DB.Exec(context.Background(), "INSERT INTO expenses (description, amount, category, user_id) VALUES ($1, $2, $3, $4)", description, amount, category, userID)
	if err != nil {
		log.Printf("Gagal memasukkan data expense: %v\n", err)
		c.String(http.StatusInternalServerError, "Gagal menyimpan pengeluaran")
		return
	}

	// Setelah insert, redirect ke halaman utama untuk melihat daftar yang sudah update
	c.Redirect(http.StatusFound, "/")
}
