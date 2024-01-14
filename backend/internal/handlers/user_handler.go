package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"regexp"

	"github.com/gorilla/sessions"
	"github.com/kiloMIA/mediscan/backend/internal/controllers"
	"github.com/kiloMIA/mediscan/backend/internal/models"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userc   *controllers.UserController
	store   *sessions.CookieStore
	lg      *logrus.Logger
}

func NewUserHandler(userc *controllers.UserController, store *sessions.CookieStore, lg *logrus.Logger) *UserHandler {
	return &UserHandler{
		userc:   userc,
		store:   store,
		lg:      lg,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (userh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	userh.lg.Debugln("User Registration at handler level")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		userh.lg.Errorf("user handler - Register - json decoder - %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if IsValidPassword(user.Password) {
		http.Error(w, "Password must be at least 10 symbols, contain at least one uppercase letter and one lowercase letter.", http.StatusBadRequest)
		return
	}

	err = userh.userc.CreateUser(r.Context(), &user)
	if err != nil {
		userh.lg.Errorf("user handler - register - user create - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, "User registered sucessfully", http.StatusCreated)
}

func (userh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	userh.lg.Debugln("User Login at handler level")
	var user loginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		userh.lg.Errorf("user handler - Login - json decoder - %v", err)
		http.Error(w, err.Error(),http.StatusBadRequest)
		return
	}

	authenticatedUser, err := userh.userc.Authenticate(r.Context(), user.Email, user.Password)
	if err != nil {
		userh.lg.Errorf("user handler - Login - authenticate - %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, _ := userh.store.Get(r, "session-name")
	session.Values["user_id"] = authenticatedUser.ID
	session.Values["user_email"] = authenticatedUser.Email
	session.Save(r, w)

	http.Error(w, "User logged in sucessfully", http.StatusCreated)
}

func (userh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := userh.store.Get(r, "session-name")

	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Error(w, "User logged out successfully", http.StatusOK)
}

func (userh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	session, err := userh.store.Get(r, "session-name")
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	user, err := userh.userc.GetUserByID(r.Context(), userID)
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - get user by id - %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (userh *UserHandler) BookCheckUp(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	name := r.FormValue("name")
	cellphoneNumber := r.FormValue("cellphoneNumber")

	body := fmt.Sprintf("Date: %s\nName: %s\nCellphone Number: %s",date, name,  cellphoneNumber)

	recipientEmail := os.Getenv("RECIPIENT_EMAIL")

	if err := SendEmail(recipientEmail, "New Booking for a Check Up", body); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inquiry submitted successfully"))
}

func IsValidPassword(password string) bool {
    if len(password) < 10 {
        return false
    }

    uppercase := regexp.MustCompile(`[A-Z]`)
    if !uppercase.MatchString(password) {
        return false
    }

    lowercase := regexp.MustCompile(`[a-z]`)
	return lowercase.MatchString(password)
}

func (userh *UserHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := userh.store.Get(r, "session-name")
		_, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SendEmail(to, subject, body string) error {
	from := os.Getenv("MAIL")
	password := os.Getenv("MAIL_PASSWORD")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

