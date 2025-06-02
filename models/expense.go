package models

import "time"

type Expense struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"` // Untuk mengaitkan dengan pengguna di Laravel
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	ExpenseDate time.Time `json:"expense_date"`
}