package dnd5e

import (
	"encoding/json"
	"testing"
)

func TestChoice(t *testing.T) {
	payload := `[
    {
      "desc": "(a) a light crossbow and 20 bolts or (b) any simple weapon",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "multiple",
            "items": [
              {
                "option_type": "counted_reference",
                "count": 1,
                "of": {
                  "index": "crossbow-light",
                  "name": "Crossbow, light",
                  "url": "/api/equipment/crossbow-light"
                }
              },
              {
                "option_type": "counted_reference",
                "count": 20,
                "of": {
                  "index": "crossbow-bolt",
                  "name": "Crossbow bolt",
                  "url": "/api/equipment/crossbow-bolt"
                }
              }
            ]
          },
          {
            "option_type": "choiceResult",
            "choiceResult": {
              "desc": "any simple weapon",
              "choose": 1,
              "type": "equipment",
              "from": {
                "option_set_type": "equipment_category",
                "equipment_category": {
                  "index": "simple-weapons",
                  "name": "Simple Weapons",
                  "url": "/api/equipment-categories/simple-weapons"
                }
              }
            }
          }
        ]
      }
    },
    {
      "desc": "(a) a component pouch or (b) an arcane focus",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "component-pouch",
              "name": "Component pouch",
              "url": "/api/equipment/component-pouch"
            }
          },
          {
            "option_type": "choiceResult",
            "choiceResult": {
              "desc": "arcane focus",
              "choose": 1,
              "type": "equipment",
              "from": {
                "option_set_type": "equipment_category",
                "equipment_category": {
                  "index": "arcane-foci",
                  "name": "Arcane Foci",
                  "url": "/api/equipment-categories/arcane-foci"
                }
              }
            }
          }
        ]
      }
    },
    {
      "desc": "(a) a scholar’s pack or (b) a dungeoneer’s pack",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "scholars-pack",
              "name": "Scholar's Pack",
              "url": "/api/equipment/scholars-pack"
            }
          },
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "dungeoneers-pack",
              "name": "Dungeoneer's Pack",
              "url": "/api/equipment/dungeoneers-pack"
            }
          }
        ]
      }
    },
    {
      "desc": "any simple weapon",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "equipment_category",
        "equipment_category": {
          "index": "simple-weapons",
          "name": "Simple Weapons",
          "url": "/api/equipment-categories/simple-weapons"
        }
      }
    }
  ]`
	c := []*choiceResult{}
	err := json.Unmarshal([]byte(payload), &c)
	if err != nil {
		t.Error(err)
	}

}
