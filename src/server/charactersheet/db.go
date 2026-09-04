package charactersheet

import (
	"context"
	"errors"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

var ErrSheetNotFound = errors.New("character sheet not found")

type ICharacterSheetService interface {
	Create(sheet *CharacterSheet) (*CharacterSheet, error)
	Read(sheetUuid uuid.UUID) (*CharacterSheet, error)
	ReadAll() ([]CharacterSheet, error)
	ReadByUser(userUuid uuid.UUID) ([]CharacterSheet, error)
	Update(sheetUuid uuid.UUID, sheet *CharacterSheet) (*CharacterSheet, error)
	Delete(sheetUuid uuid.UUID) error
}

type CharacterSheetService struct {
	db     *mongo.Database
	logger *zap.SugaredLogger
}

func NewCharacterSheetCrudService() ICharacterSheetService {
	var service ICharacterSheetService
	app.Invoke(func(db *mongo.Database, logger *zap.SugaredLogger) {
		service = &CharacterSheetService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (s *CharacterSheetService) sheetsColl() *mongo.Collection {
	return s.db.Collection("character_sheets")
}

func (s *CharacterSheetService) Create(sheet *CharacterSheet) (*CharacterSheet, error) {
	if sheet.Uuid == uuid.Nil {
		sheet.Uuid = uuid.New()
	}
	if sheet.CreatedAt.IsZero() {
		sheet.CreatedAt = time.Now()
	}
	sheet.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.sheetsColl().InsertOne(ctx, sheet)
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		sheet.ID = oid
	}

	return sheet, nil
}

func (s *CharacterSheetService) Read(sheetUuid uuid.UUID) (*CharacterSheet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sheet CharacterSheet
	err := s.sheetsColl().FindOne(ctx, bson.M{"uuid": sheetUuid}).Decode(&sheet)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSheetNotFound
		}
		return nil, err
	}
	return &sheet, nil
}

func (s *CharacterSheetService) ReadAll() ([]CharacterSheet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.sheetsColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sheets []CharacterSheet
	if err := cursor.All(ctx, &sheets); err != nil {
		return nil, err
	}
	return sheets, nil
}

func (s *CharacterSheetService) ReadByUser(userUuid uuid.UUID) ([]CharacterSheet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.sheetsColl().Find(ctx, bson.M{"user_uuid": userUuid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sheets []CharacterSheet
	if err := cursor.All(ctx, &sheets); err != nil {
		return nil, err
	}
	return sheets, nil
}

func (s *CharacterSheetService) Update(sheetUuid uuid.UUID, sheet *CharacterSheet) (*CharacterSheet, error) {
	existing, err := s.Read(sheetUuid)
	if err != nil {
		return nil, err
	}

	sheet.ID = existing.ID
	sheet.Uuid = sheetUuid
	sheet.CreatedAt = existing.CreatedAt
	sheet.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.sheetsColl().ReplaceOne(ctx, bson.M{"uuid": sheetUuid}, sheet)
	if err != nil {
		return nil, err
	}

	return sheet, nil
}

func (s *CharacterSheetService) Delete(sheetUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.sheetsColl().DeleteOne(ctx, bson.M{"uuid": sheetUuid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrSheetNotFound
	}
	return nil
}
