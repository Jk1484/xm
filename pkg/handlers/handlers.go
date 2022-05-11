package handlers

import (
	"xm/pkg/logger"
	userRepository "xm/pkg/repositories/user"
	"xm/pkg/services/company"
	userService "xm/pkg/services/user"

	"encoding/json"
	"net/http"
	"time"
	"xm/pkg/services/utils"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var Module = fx.Provide(New)

type Handlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Middleware(next http.Handler) http.Handler
	LogRequest(next http.Handler) http.Handler

	CreateCompany(w http.ResponseWriter, r *http.Request)
	GetCompanyByID(w http.ResponseWriter, r *http.Request)
	GetAllCompanies(w http.ResponseWriter, r *http.Request)
	UpdateCompany(w http.ResponseWriter, r *http.Request)
	DeleteCompany(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	userService    userService.Service
	companyService company.Service
	logger         logger.Logger
}

type Params struct {
	fx.In
	UserService    userService.Service
	CompanyService company.Service
	Logger         logger.Logger
}

func New(p Params) Handlers {
	return &handlers{
		userService:    p.UserService,
		companyService: p.CompanyService,
		logger:         p.Logger,
	}
}

var jwtKey = []byte("secret_key")

type Claims struct {
	userRepository.User
	jwt.StandardClaims
}

func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	var credentials userRepository.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "bad credentials")
		return
	}

	if credentials.Username == "" {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "no username provided")
		return
	}

	if len(credentials.Password) < 8 {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "password minimum length should be at least 8")
		return
	}

	hashedPassword, err := HashPassword(credentials.Password)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	user := &userRepository.User{
		Username: credentials.Username,
		Password: hashedPassword,
	}

	err = h.userService.Create(user)
	if err == utils.ErrAlreadyExists {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "already registered")
		return
	}

	if err != nil {
		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), "sign up completed")
}

func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	var credentials userRepository.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	if credentials.Username == "" {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "no username provided")
		return
	}

	user, err := h.userService.GetByUsername(credentials.Username)
	if err != nil {
		if err == utils.ErrNotFound {
			apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "incorrect username or password")
			return
		}

		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}

	if !CheckPasswordHash(credentials.Password, user.Password) {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "incorrect username or password")
		return
	}

	expirationTime := time.Now().Add(time.Hour * 2)

	claims := &Claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), tokenString)
}

func (h *handlers) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var apiResp ApiResp

		tokenStr := r.Header.Get("token")

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				apiResp.Set(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
				apiResp.Respond(w)
				return
			}

			apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
			apiResp.Respond(w)
			return
		}

		if !tkn.Valid {
			apiResp.Set(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			apiResp.Respond(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *handlers) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer h.logger.Logger().Infof("time taken %v", time.Since(start))

		next.ServeHTTP(w, r)
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetClaims(r *http.Request) (*Claims, error) {
	claims := &Claims{}

	tokenStr := r.Header.Get("token")

	_, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

type ApiResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func (a *ApiResp) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(a.Code)

	resp, _ := json.Marshal(a)

	w.Write(resp)
}

func (a *ApiResp) Set(code int, message string, payload interface{}) {
	a.Code = code
	a.Message = message
	a.Payload = payload
}
