package mud

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// CreatureTypes is a mapping of string IDs to creature types
var CreatureTypes map[string]CreatureType

// CreatureType is the type of creature (Hostile: true is monster, false is NPC)
type CreatureType struct {
	ID      string `json:""`
	Name    string `json:""`
	Hostile bool   `json:""`
	MaxHP   uint64 `json:""`
	MaxMP   uint64 `json:""`
	MaxAP   uint64 `json:""`
	MaxRP   uint64 `json:""`
}

// Creature is an instance of a Creature
type Creature struct {
	ID                 string       `json:""`
	CreatureType       string       `json:""`
	HP                 uint64       `json:""`
	MP                 uint64       `json:""`
	AP                 uint64       `json:""`
	RP                 uint64       `json:""`
	CreatureTypeStruct CreatureType `json:"-"`
	world              World
}

// CreatureList represents the creatures in a DB
type CreatureList struct {
	CreatureIDs []string `json:""`
}

func init() {
	CreatureTypes = make(map[string]CreatureType)

	creatureInfoFile := "./bestiary.json"
	data, err := ioutil.ReadFile(creatureInfoFile)

	if err == nil {
		err = json.Unmarshal(data, &CreatureTypes)
	}

	if err != nil {
		log.Printf("Error parsing %s: %v", creatureInfoFile, err)
	}
}