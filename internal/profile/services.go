package profile

import (
	"context"
	"errors"
	"regexp"

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

	// 	profile := &UserProfile{
	// 		UserID:      userID,
	// 		Username:    defUsername,
	// 		FullName:    fullName,
	// 		DisplayName: fullName,
	// 	}

	// 	ok, err := ps.store.InsertDefaultProfile(ctx, profile)
	// 	if err == nil && ok {
	// 		return nil // successfully inserted the profile
	// 	}

	// 	if err == nil && !ok {
	// 		continue // username was taken, try again with a new username
	// 	}

	// 	// error not due to username being taken, return the error
	// 	if !errors.Is(err, ErrUsernameTaken) {
	// 		return err
	// 	}
	// }

	// return errors.New("failed to generate a unique username after maximum retries")

	return nil
}

// func generateDefaultUsername(fullName string) string {
// 	sliceFullName := strings.Fields(strings.ToLower(fullName))

// 	var username string

// 	if len(sliceFullName) >= 2 {
// 		username = strings.Join(sliceFullName[:2], "")
// 	} else {
// 		username = strings.Join(sliceFullName, "")
// 	}

// 	username = nonAlphaNumericRegex.ReplaceAllString(username, "")

// 	if utf8.RuneCountInString(username) < 3 {
// 		username = "real_" + username
// 	}

// 	runes := []rune(username)
// 	if len(runes) > 24 {
// 		username = string(runes[:24])
// 	}

// 	return username
// }

// // RandomSuffix generates a short alphanumeric suffix (e.g., "9b2a").
// func RandomSuffix(length int) string {
// 	b := make([]byte, length)
// 	rand.Read(b)

// 	return hex.EncodeToString(b)
// }
