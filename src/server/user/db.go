package user

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"

	"github.com/google/uuid"
	"github.com/xrash/smetrics"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

var ErrRecordNotFound = errors.New("record not found")

type IUserCrudService interface {
	Create(user *User, password string) (*User, error)
	Read(uuid uuid.UUID) (*User, error)
	ReadAll() ([]User, error)
	Update(uuid uuid.UUID, user *User) (*User, error)
	Delete(uuid uuid.UUID) error
	GetAllUsers() ([]User, error)
	SearchUsersByName(query string) ([]User, error)
	GetUserByEmail(email string) (*User, error)
}

type UserCrudService struct {
	db     *mongo.Database
	logger *zap.SugaredLogger
}

type UserWithScore struct {
	User  User
	Score float64
}

func NewUserCrudService() IUserCrudService {
	var service IUserCrudService
	app.Invoke(func(db *mongo.Database, logger *zap.SugaredLogger) {
		service = &UserCrudService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (u *UserCrudService) usersColl() *mongo.Collection {
	return u.db.Collection("users")
}

func (u *UserCrudService) ReadAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := u.usersColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserCrudService) Delete(userUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing User
	err := u.usersColl().FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&existing)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrRecordNotFound
		}
		return err
	}

	existing.Username = fmt.Sprintf("deleted_user_%s", userUuid.String()[:8])
	existing.Email = fmt.Sprintf("deleted_%s@example.com", userUuid.String()[:8])
	existing.PasswordHash = ""
	existing.UpdatedAt = time.Now()

	_, err = u.usersColl().ReplaceOne(ctx, bson.M{"uuid": userUuid}, existing)
	if err != nil {
		return err
	}

	_, err = u.usersColl().DeleteOne(ctx, bson.M{"uuid": userUuid})
	return err
}

func (u *UserCrudService) Read(userUuid uuid.UUID) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := u.usersColl().FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserCrudService) Update(userUuid uuid.UUID, user *User) (*User, error) {
	userOld, err := u.Read(userUuid)
	if err != nil {
		return nil, err
	}

	userOld = userOld.Update(user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = u.usersColl().ReplaceOne(ctx, bson.M{"uuid": userUuid}, userOld)
	if err != nil {
		return nil, err
	}
	return userOld, nil
}

func (u *UserCrudService) Create(user *User, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing User
	err := u.usersColl().FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = hash
	if user.Uuid == uuid.Nil {
		user.Uuid = uuid.New()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	res, err := u.usersColl().InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}

func (u *UserCrudService) GetAllUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := u.usersColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserCrudService) SearchUsersByName(query string) ([]User, error) {
	users, err := u.ReadAll()
	if err != nil {
		return nil, err
	}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	var scoredUsers []UserWithScore
	for _, user := range users {
		fullName := strings.ToLower(user.Username + " " + user.Email)
		score := smetrics.JaroWinkler(normalizedQuery, fullName, 0.7, 4)
		scoredUsers = append(scoredUsers, UserWithScore{User: user, Score: score})
	}

	sort.Slice(scoredUsers, func(i, j int) bool {
		return scoredUsers[i].Score > scoredUsers[j].Score
	})

	var filteredUsers []User
	for _, scoredUser := range scoredUsers {
		if scoredUser.Score >= 0.8 {
			filteredUsers = append(filteredUsers, scoredUser.User)
		}
	}

	return filteredUsers, nil
}

func (u *UserCrudService) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := u.usersColl().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
