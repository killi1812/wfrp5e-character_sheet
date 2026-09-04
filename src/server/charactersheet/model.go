package charactersheet

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StatValue struct {
	Initial  int `json:"initial" bson:"initial"`
	Advances int `json:"advances" bson:"advances"`
	Current  int `json:"current" bson:"current"`
}

type Characteristics struct {
	WS  StatValue `json:"ws" bson:"ws"`
	BS  StatValue `json:"bs" bson:"bs"`
	S   StatValue `json:"s" bson:"s"`
	T   StatValue `json:"t" bson:"t"`
	I   StatValue `json:"i" bson:"i"`
	Ag  StatValue `json:"ag" bson:"ag"`
	Dex StatValue `json:"dex" bson:"dex"`
	Int StatValue `json:"int" bson:"int"`
	WP  StatValue `json:"wp" bson:"wp"`
	Fel StatValue `json:"fel" bson:"fel"`
}

type Language struct {
	Name  string `json:"name" bson:"name"`
	Int   int    `json:"int" bson:"int"`
	Adv   int    `json:"adv" bson:"adv"`
	Skill int    `json:"skill" bson:"skill"`
}

type Skill struct {
	Name           string `json:"name" bson:"name"`
	Characteristic string `json:"characteristic" bson:"characteristic"`
	Adv            int    `json:"adv" bson:"adv"`
	Total          int    `json:"total" bson:"total"`
	Type           string `json:"type" bson:"type"` // Basic or Advanced
}

type Talent struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Page        string `json:"page" bson:"page"`
}

type Wounds struct {
	Current int `json:"current" bson:"current"`
	Max     int `json:"max" bson:"max"`
	SB      int `json:"sb" bson:"sb"`
	TBx2    int `json:"tbX2" bson:"tb_x2"`
	WPB     int `json:"wpb" bson:"wpb"`
	Hardy   int `json:"hardy" bson:"hardy"`
}

type ArmourLocation struct {
	Head         int `json:"head" bson:"head"`
	PrimaryArm   int `json:"primaryArm" bson:"primary_arm"`
	SecondaryArm int `json:"secondaryArm" bson:"secondary_arm"`
	Body         int `json:"body" bson:"body"`
	PrimaryLeg   int `json:"primaryLeg" bson:"primary_leg"`
	SecondaryLeg int `json:"secondaryLeg" bson:"secondary_leg"`
	Shield       int `json:"shield" bson:"shield"`
}

type Wealth struct {
	GC int `json:"gc" bson:"gc"` // Gold Crowns
	SS int `json:"ss" bson:"ss"` // Silver Shillings
	BP int `json:"bp" bson:"bp"` // Brass Pennies
}

type Encumbrance struct {
	Weapons   int `json:"weapons" bson:"weapons"`
	Armour    int `json:"armour" bson:"armour"`
	Trappings int `json:"trappings" bson:"trappings"`
	Other     int `json:"other" bson:"other"`
	Max       int `json:"max" bson:"max"`
	Total     int `json:"total" bson:"total"`
}

type Mutation struct {
	Name   string `json:"name" bson:"name"`
	Effect string `json:"effect" bson:"effect"`
}

type WeaponItem struct {
	Name       string `json:"name" bson:"name"`
	Group      string `json:"group" bson:"group"`
	Enc        int    `json:"enc" bson:"enc"`
	RangeReach string `json:"rangeReach" bson:"range_reach"`
	Damage     string `json:"damage" bson:"damage"`
	Qualities  string `json:"qualities" bson:"qualities"`
}

type ArmourItem struct {
	Name      string `json:"name" bson:"name"`
	Locations string `json:"locations" bson:"locations"`
	Enc       int    `json:"enc" bson:"enc"`
	AP        int    `json:"ap" bson:"ap"`
	Qualities string `json:"qualities" bson:"qualities"`
}

type TrappingItem struct {
	Name        string `json:"name" bson:"name"`
	Category    string `json:"category" bson:"category"`
	Enc         int    `json:"enc" bson:"enc"`
	Description string `json:"description" bson:"description"`
}

type SpellItem struct {
	Name        string `json:"name" bson:"name"`
	CN          int    `json:"cn" bson:"cn"`
	Range       string `json:"range" bson:"range"`
	Target      string `json:"target" bson:"target"`
	Duration    string `json:"duration" bson:"duration"`
	Description string `json:"description" bson:"description"`
	Sin         int    `json:"sin" bson:"sin"`
}

type CareerAdvancementTracker struct {
	Tier1 bool `json:"tier1" bson:"tier1"`
	Tier2 bool `json:"tier2" bson:"tier2"`
	Tier3 bool `json:"tier3" bson:"tier3"`
	Tier4 bool `json:"tier4" bson:"tier4"`
}

type CharacterSheet struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Uuid      uuid.UUID     `json:"uuid" bson:"uuid"`
	UserUuid  uuid.UUID     `json:"userUuid" bson:"user_uuid"`

	// Header Info
	Name               string `json:"name" bson:"name"`
	Species            string `json:"species" bson:"species"`
	Appearance         string `json:"appearance" bson:"appearance"`
	Class              string `json:"class" bson:"class"`
	Career             string `json:"career" bson:"career"`
	CareerLevel        int    `json:"careerLevel" bson:"career_level"`
	CareerPath         string `json:"careerPath" bson:"career_path"`
	CareerAdvancement  CareerAdvancementTracker `json:"careerAdvancement" bson:"career_advancement"`
	Status             string `json:"status" bson:"status"`

	// XP
	XpCurrent int `json:"xpCurrent" bson:"xp_current"`
	XpSpent   int `json:"xpSpent" bson:"xp_spent"`
	XpTotal   int `json:"xpTotal" bson:"xp_total"`

	// Characteristics
	Characteristics Characteristics `json:"characteristics" bson:"characteristics"`

	// Fate & Fortune, Resilience & Resolve, Movement
	Fate       int `json:"fate" bson:"fate"`
	Fortune    int `json:"fortune" bson:"fortune"`
	Resilience int `json:"resilience" bson:"resilience"`
	Resolve    int `json:"resolve" bson:"resolve"`
	Movement   int `json:"movement" bson:"movement"`
	Walk       int `json:"walk" bson:"walk"`
	Run        int `json:"run" bson:"run"`

	// Ambitions & Details
	PersonalAmbition string `json:"personalAmbition" bson:"personal_ambition"`
	PartyAmbition    string `json:"partyAmbition" bson:"party_ambition"`

	// Sub-lists
	Languages        []Language      `json:"languages" bson:"languages"`
	Skills           []Skill         `json:"skills" bson:"skills"`
	Talents          []Talent        `json:"talents" bson:"talents"`
	Wounds           Wounds          `json:"wounds" bson:"wounds"`
	ArmourPoints     ArmourLocation  `json:"armourPoints" bson:"armour_points"`
	Wealth           Wealth          `json:"wealth" bson:"wealth"`
	Encumbrance      Encumbrance     `json:"encumbrance" bson:"encumbrance"`
	CorruptionPoints int             `json:"corruptionPoints" bson:"corruption_points"`
	Mutations        []Mutation      `json:"mutations" bson:"mutations"`
	Weapons          []WeaponItem    `json:"weapons" bson:"weapons"`
	Armour           []ArmourItem    `json:"armour" bson:"armour"`
	Trappings        []TrappingItem  `json:"trappings" bson:"trappings"`
	SpellsAndPrayers []SpellItem     `json:"spellsAndPrayers" bson:"spells_and_prayers"`
	Notes            string          `json:"notes" bson:"notes"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
