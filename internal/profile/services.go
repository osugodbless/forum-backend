package profile

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

var nonAlphaNumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrProfileExists = errors.New("profile already exists for this user")
)

type ProfileService struct {
	store profileStore
}

func NewProfileService(store profileStore) *ProfileService {
	return &ProfileService{store: store}
}

func (ps *ProfileService) CreateProfile(ctx context.Context, userID uuid.UUID, fullName string) error {
	base := generateDefaultUsername(fullName)
	defUsername := base

	const maxRetries = 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			defUsername = fmt.Sprintf("%s%s", base, RandomSuffix(2))
		}
		profile := &UserProfile{
			UserID:      userID,
			Username:    defUsername,
			FullName:    fullName,
			DisplayName: fullName,
		}
		err := ps.store.InsertProfile(ctx, profile)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateDefaultUsername(fullName string) string {
	s := strings.ToLower(fullName)
	s = nonAlphaNumericRegex.ReplaceAllString(s, "")
	if utf8.RuneCountInString(s) < 3 {
		s = "real_" + s
	}

	runes := []rune(s)
	if len(runes) > 24 {
		s = string(runes[:24])
	}
	return s
}

// RandomSuffix generates a short alphanumeric suffix (e.g., "9x2a").
func RandomSuffix(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	return hex.EncodeToString(b)
}
