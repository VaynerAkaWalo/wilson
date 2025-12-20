package location

import (
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"math/rand/v2"
)

const (
	OCEAN     TerrainType = "OCEAN"
	BEACH     TerrainType = "BEACH"
	FOREST    TerrainType = "FOREST"
	MOUNTAINS TerrainType = "MOUNTAINS"
)

var (
	terrains = []string{
		"Cove", "Lagoon", "Bluff", "Estuary",
		"Beach", "Jungle", "Reef", "River",
		"Peak", "Swamp", "Delta", "Forest",
		"Shore", "Cave", "Flats", "Slope",
		"Clearing", "Spring", "Tide", "Crotto",
		"Wreckage", "Falls", "Ridge", "Thicket",
		"Volcano", "Valley", "Basin", "Plateau",
		"Strand", "Islet", "Canyon", "Crest"}

	conditions = []string{
		"Hidden", "Sunken", "Roaring", "Dead",
		"Feral", "Serpent", "Salt", "Dry",
		"Scorched", "Obsidian", "Coral", "Poison",
		"Whispering", "Shallow", "Crab", "Damp",
		"Crimson", "High", "Black", "Dense",
		"Winding", "Iron", "Turtle", "Green",
		"Storm", "Fog", "Silent", "Lost",
		"Woven", "Bramble", "Shifting", "First"}

	terrainTypes = []TerrainType{
		OCEAN, BEACH, FOREST, MOUNTAINS,
	}
)

type (
	TerrainType string

	Id       string
	Location struct {
		Id               Id
		Name             string
		Latitude         int
		Longitude        int
		RewardMultiplier float64
		Type             TerrainType
	}
)

func New(latitude, longitude int, multiplier float64) *Location {
	terrain := terrains[rand.IntN(len(terrains))]
	condition := conditions[rand.IntN(len(conditions))]
	terrainType := terrainTypes[rand.IntN(len(terrainTypes))]

	return &Location{
		Id:               Id(xuuid.UUID()),
		Name:             condition + " " + terrain,
		Latitude:         latitude,
		Longitude:        longitude,
		RewardMultiplier: multiplier,
		Type:             terrainType,
	}
}
