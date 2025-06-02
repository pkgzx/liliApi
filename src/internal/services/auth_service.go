package services

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/pkgzx/liliApi/src/pkg/data"
)

type AuthService struct {
    userService *UserService
    jwtSecret   string
}

func NewAuthService(userService *UserService, jwtSecret string) *AuthService {
    if jwtSecret == "" {
        jwtSecret = "default-secret-key-change-in-production"
    }
    return &AuthService{
        userService: userService,
        jwtSecret:   jwtSecret,
    }
}

type TokenClaims struct {
    UserID   int32  `json:"user_id"`
    Username string `json:"username"`
    FullName string `json:"full_name"`
    jwt.RegisteredClaims
}

type LoginResult struct {
    User  *data.User `json:"user"`
    Token string     `json:"token"`
}

func (s *AuthService) Login(username, password string) (*LoginResult, error) {
    // Autenticar usuario
    user, err := s.userService.AuthenticateUser(username, password)
    if err != nil {
        return nil, err
    }

    // Generar token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, fmt.Errorf("error generating token: %w", err)
    }

    return &LoginResult{
        User:  user,
        Token: token,
    }, nil
}

func (s *AuthService) generateToken(user *data.User) (string, error) {
    // Crear claims
    claims := TokenClaims{
        UserID:   user.ID,
        Username: user.Username,
        FullName: user.FullName,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 horas
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "liliapi",
            Subject:   fmt.Sprintf("user_%d", user.ID),
        },
    }

    // Crear token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Firmar token
    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
    // Parsear token
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verificar método de firma
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    // Verificar validez del token
    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func (s *AuthService) RefreshToken(tokenString string) (string, error) {
    claims, err := s.ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    // Verificar si el token expira en menos de 1 hora para permitir refresh
    if claims.ExpiresAt.Time.Sub(time.Now()) > time.Hour {
        return "", errors.New("token still valid, refresh not needed")
    }

    // Obtener usuario actualizado
    user, err := s.userService.userRepo.GetByUsername(claims.Username)
    if err != nil {
        return "", err
    }

    if user == nil {
        return "", errors.New("user not found")
    }

    // Generar nuevo token
    return s.generateToken(user)
}