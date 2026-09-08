package charactersheet_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/charactersheet"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func sampleSheet() *charactersheet.CharacterSheet {
	now := time.Now().Truncate(time.Millisecond)
	sheetUuid := uuid.New()
	userUuid := uuid.New()

	return &charactersheet.CharacterSheet{
		Uuid:              sheetUuid,
		UserUuid:          userUuid,
		Name:              "Gotrek Gurnisson",
		Species:           "Dwarf",
		Appearance:        "Mohawk, eyepatch, muscular",
		Class:             "Warrior",
		Career:            "Slayer",
		CareerLevel:       2,
		CareerPath:        "Troll Slayer",
		CareerAdvancement: charactersheet.CareerAdvancementTracker{Tier1: true, Tier2: true},
		Advances2:         []bool{true, false, true},
		Advances3:         []bool{false},
		Advances4:         []bool{false},
		Status:            "Brass 2",
		XpCurrent:         120,
		XpSpent:           880,
		XpTotal:           1000,
		Characteristics: charactersheet.Characteristics{
			WS:  charactersheet.StatValue{Initial: 40, Advances: 15, Current: 55},
			BS:  charactersheet.StatValue{Initial: 25, Advances: 0, Current: 25},
			S:   charactersheet.StatValue{Initial: 45, Advances: 10, Current: 55},
			T:   charactersheet.StatValue{Initial: 50, Advances: 10, Current: 60},
			I:   charactersheet.StatValue{Initial: 30, Advances: 5, Current: 35},
			Ag:  charactersheet.StatValue{Initial: 25, Advances: 5, Current: 30},
			Dex: charactersheet.StatValue{Initial: 30, Advances: 0, Current: 30},
			Int: charactersheet.StatValue{Initial: 20, Advances: 0, Current: 20},
			WP:  charactersheet.StatValue{Initial: 60, Advances: 10, Current: 70},
			Fel: charactersheet.StatValue{Initial: 15, Advances: 0, Current: 15},
		},
		Fate:             2,
		Fortune:          2,
		Resilience:       3,
		Resolve:          3,
		Movement:         3,
		Walk:             6,
		Run:              12,
		PersonalAmbition: "Find a glorious death",
		PartyAmbition:    "Cleanse the mines",
		Languages: []charactersheet.Language{
			{Name: "Khazalid", Int: 20, Adv: 10, Skill: 30},
			{Name: "Reikspiel", Int: 20, Adv: 5, Skill: 25},
		},
		Skills: []charactersheet.Skill{
			{Name: "Melee (Basic)", Characteristic: "WS", Adv: 15, Total: 70, Type: "Basic"},
			{Name: "Lore (Trolls)", Characteristic: "Int", Adv: 10, Total: 30, Type: "Advanced"},
		},
		Talents: []charactersheet.Talent{
			{Name: "Hardy", Description: "Add TB to wounds", Page: "135"},
			{Name: "Furious Assault", Description: "Extra attack on hit", Page: "136"},
		},
		Wounds: charactersheet.Wounds{
			Current: 18,
			Max:     22,
			SB:      5,
			TBx2:    12,
			WPB:     7,
			Hardy:   6,
		},
		ArmourPoints: charactersheet.ArmourLocation{
			Head:         0,
			PrimaryArm:   1,
			SecondaryArm: 1,
			Body:         2,
			PrimaryLeg:   0,
			SecondaryLeg: 0,
			Shield:       0,
		},
		Wealth: charactersheet.Wealth{
			GC: 5,
			SS: 14,
			BP: 22,
		},
		Encumbrance: charactersheet.Encumbrance{
			Weapons:   4,
			Armour:    3,
			Trappings: 5,
			Other:     0,
			Max:       15,
			Total:     12,
		},
		CorruptionPoints: 1,
		Sin:              0,
		Mutations: []charactersheet.Mutation{
			{Name: "Iron Skin", Effect: "+1 AP all locations"},
		},
		Weapons: []charactersheet.WeaponItem{
			{Name: "Rune Axe", Group: "Two-Handed", Enc: 3, RangeReach: "Average", Damage: "+SB+6", Qualities: "Impact, Hack"},
		},
		Armour: []charactersheet.ArmourItem{
			{Name: "Leather Jerkin", Locations: "Body", Enc: 1, AP: 1, Qualities: ""},
		},
		Trappings: []charactersheet.TrappingItem{
			{Name: "Troll skull mug", Category: "Keepsake", Enc: 1, Description: "A mug made of bone"},
		},
		SpellsAndPrayers: []charactersheet.SpellItem{
			{Name: "Oath of Vengeance", CN: 0, Range: "Self", Target: "Self", Duration: "1 round", Description: "Gain hatred", Sin: 0},
		},
		Notes:     "A great companion to Felix.",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestCharacterSheet_JsonSerialization(t *testing.T) {
	orig := sampleSheet()

	data, err := json.Marshal(orig)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded charactersheet.CharacterSheet
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, orig.Uuid, decoded.Uuid)
	assert.Equal(t, orig.UserUuid, decoded.UserUuid)
	assert.Equal(t, orig.Name, decoded.Name)
	assert.Equal(t, orig.Species, decoded.Species)
	assert.Equal(t, orig.Career, decoded.Career)
	assert.Equal(t, orig.CareerAdvancement.Tier1, decoded.CareerAdvancement.Tier1)
	assert.Equal(t, orig.Characteristics.WS.Current, decoded.Characteristics.WS.Current)
	assert.Equal(t, orig.Characteristics.WP.Current, decoded.Characteristics.WP.Current)
	assert.Len(t, decoded.Languages, 2)
	assert.Len(t, decoded.Skills, 2)
	assert.Len(t, decoded.Talents, 2)
	assert.Len(t, decoded.Weapons, 1)
	assert.Len(t, decoded.Armour, 1)
	assert.Len(t, decoded.Trappings, 1)
	assert.Len(t, decoded.SpellsAndPrayers, 1)
	assert.Len(t, decoded.Mutations, 1)
	assert.Equal(t, orig.Wealth.GC, decoded.Wealth.GC)
	assert.Equal(t, orig.Encumbrance.Total, decoded.Encumbrance.Total)
}

func TestCharacterSheet_BsonSerialization(t *testing.T) {
	orig := sampleSheet()
	orig.ID = bson.NewObjectID()

	data, err := bson.Marshal(orig)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded charactersheet.CharacterSheet
	err = bson.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, orig.ID, decoded.ID)
	assert.Equal(t, orig.Uuid, decoded.Uuid)
	assert.Equal(t, orig.UserUuid, decoded.UserUuid)
	assert.Equal(t, orig.Name, decoded.Name)
	assert.Equal(t, orig.Wounds.Current, decoded.Wounds.Current)
}
