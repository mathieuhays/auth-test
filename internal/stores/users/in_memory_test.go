package users

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestUserStore(t *testing.T) {
	userStore := NewUserMemoryStore()
	user := User{
		ID:             uuid.New(),
		Email:          "some@email.com",
		EmailConfirmed: nil,
		PasswordHash:   "password",
	}

	u, err := userStore.Create(user)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if u.ID != user.ID {
		t.Fatalf("id doesn't match. expected: %v. got: %v", user.ID, u.ID)
	}

	if u.CreatedAt.IsZero() {
		t.Fatalf("createdAt expected to default to 'now' if not defined. value return is a zero value")
	}

	u2, err := userStore.Get(user.ID)
	if err != nil {
		t.Fatalf("unexpected Get error on existing user: %s", err)
	}

	if u2.ID != user.ID {
		t.Fatalf("Get: IDs don't match. expected: %v. got: %v", user.ID, u2.ID)
	}

	wrongID := uuid.New()
	_, err = userStore.Get(wrongID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("Get wrong id. unexpected error: %v", err)
	}

	updatedUser := User{
		ID:             user.ID,
		Email:          "new@email.com",
		EmailConfirmed: nil,
		PasswordHash:   "password",
	}
	u3, err := userStore.Update(updatedUser)
	if err != nil {
		t.Fatalf("update: unexpected error: %s", err)
	}

	if u3.Email != updatedUser.Email {
		t.Fatalf("value not updated. expected: %s. got: %s", updatedUser.Email, u3.Email)
	}

	err = userStore.Delete(user.ID)
	if err != nil {
		t.Fatalf("delete: unexpected error: %s", err)
	}

	u4, err := userStore.Get(user.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("get after delete: unexpected error: %v", err)
	}

	if u4 != nil {
		t.Fatalf("get after delete: user should be nil, got: %v", u4)
	}

	userWithDate := User{
		ID:             uuid.New(),
		Email:          "random@example.com",
		EmailConfirmed: nil,
		PasswordHash:   "blahblahblah",
		CreatedAt:      time.Now().Add(time.Hour * -2),
	}
	u5, err := userStore.Create(userWithDate)
	if err != nil {
		t.Fatalf("unexpected error when creating user with date: %s", err)
	}

	if !u5.CreatedAt.Equal(userWithDate.CreatedAt) {
		t.Fatalf("CreatedAt different after create(). expected: %v. got: %v", userWithDate.CreatedAt, u5.CreatedAt)
	}
}
