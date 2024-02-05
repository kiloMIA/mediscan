package internal

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
    DB *pgxpool.Pool
}

func NewUserController(db *pgxpool.Pool) *UserController {
    return &UserController{DB: db}
}

func (uc *UserController) BookAppointment(discordID, date, time, details string) error {
    sql := `INSERT INTO appointments (discord_id, timestamp, details) VALUES ($1, $2, $3)`
    _, err := uc.DB.Exec(context.Background(), sql, discordID, date+" "+time, details)
    return err
}

func (uc *UserController) UploadDocument(discordID, filePath string) error {
    sql := `INSERT INTO documents (discord_id, file_path) VALUES ($1, $2)`
    _, err := uc.DB.Exec(context.Background(), sql, discordID, filePath)
    return err
}