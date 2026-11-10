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

func (ps *ProfileService) CreateDefaultProfile(ctx context.Context, userID uuid.UUID, fullName string) error {
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

		err := ps.store.InsertDefaultProfile(ctx, profile)
		if err == nil {
			return nil // Successfully inserted the profile
		}

		// error not due to username being taken, return the error
		if !errors.Is(err, ErrUsernameTaken) {
			return err
		}
	}

	return errors.New("failed to generate a unique username after maximum retries")
}

func generateDefaultUsername(fullName string) string {
	sliceFullName := strings.Fields(strings.ToLower(fullName))

	var username string

	if len(sliceFullName) >= 2 {
		username = strings.Join(sliceFullName[:2], "")
	} else {
		username = strings.Join(sliceFullName, "")
	}

	username = nonAlphaNumericRegex.ReplaceAllString(username, "")

	if utf8.RuneCountInString(username) < 3 {
		username = "real_" + username
	}

	runes := []rune(username)
	if len(runes) > 24 {
		username = string(runes[:24])
	}

	return username
}

// RandomSuffix generates a short alphanumeric suffix (e.g., "9x2a").
func RandomSuffix(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	return hex.EncodeToString(b)
}
