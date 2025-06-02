package services

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "errors"
    "fmt"
    "strings"

    "golang.org/x/crypto/argon2"
    "github.com/pkgzx/liliApi/src/pkg/data"
    "github.com/pkgzx/liliApi/src/pkg/repository"
)

type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

func (s *UserService) AuthenticateUser(username, password string) (*data.User, error) {
    user, err := s.userRepo.GetByUsername(username)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    if !s.verifyPassword(password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}

func (s *UserService) CreateUser(username, password, fullName string) (*data.User, error) {
    // Verificar si el usuario ya existe
    existingUser, err := s.userRepo.GetByUsername(username)
    if err != nil {
        return nil, err
    }
    
    if existingUser != nil {
        return nil, errors.New("username already exists")
    }

    passwordHash, err := s.hashPassword(password)
    if err != nil {
        return nil, fmt.Errorf("error hashing password: %w", err)
    }

    user := &data.User{
        Username:     username,
        PasswordHash: passwordHash,
        FullName:     fullName,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// Hash password usando Argon2
func (s *UserService) hashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version, 64*1024, 1, 4, b64Salt, b64Hash), nil
}

// Verificar password
func (s *UserService) verifyPassword(password, hash string) bool {
    parts := strings.Split(hash, "$")
    if len(parts) != 6 {
        return false
    }

    var memory, time uint32
    var parallelism uint8
    _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
    if err != nil {
        return false
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false
    }

    decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false
    }

    comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, parallelism, uint32(len(decodedHash)))

    return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}