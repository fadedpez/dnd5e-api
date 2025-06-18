# D&D Room Generator Module Design

## Overview

A Go module/plugin system for generating dungeon rooms with monsters, loot, and environmental features for the D&D Discord bot.

## Core Concept

```go
// github.com/yourguy/dnd-room-generator
package roomgen

type Generator interface {
    GenerateRoom(config RoomConfig) (*Room, error)
    ValidateConfig(config RoomConfig) error
}

type RoomConfig struct {
    // Core parameters
    Difficulty   DifficultyLevel // easy, medium, hard, deadly
    PartyLevel   int             // Average party level
    PartySize    int             // Number of players
    
    // Theming
    DungeonTheme string          // "undead_crypt", "goblin_cave", "dragon_lair"
    RoomType     string          // "combat", "puzzle", "treasure", "mixed"
    
    // Optional constraints
    MaxMonsters  int
    MinLoot      int
    IncludeTraps bool
    
    // Dungeon progression
    DungeonDepth int             // How many rooms deep are we?
    RoomsCleared int             // How many rooms has the party cleared?
    BossDefeated bool            // Has a boss been encountered yet?
}

type Room struct {
    // Descriptive elements
    Name        string
    Description string
    Atmosphere  string // "eerie", "ancient", "foul-smelling"
    
    // Gameplay elements
    Monsters    []Monster
    Loot        []LootItem
    Traps       []Trap
    Features    []EnvironmentalFeature
    
    // Connections
    Exits       []Exit
    
    // Metadata
    ChallengeRating float64
    XPReward        int
    
    // Room type hints
    IsBossRoom  bool
    IsRestSafe  bool // Can the party safely rest here?
    HasShortcut bool // Does this room provide a shortcut back?
}
```

## Integration Approaches

### 1. Direct Module Import
```go
import "github.com/yourguy/dnd-room-generator"

func generateDungeonRoom(partyLevel int) {
    gen := roomgen.New()
    room, err := gen.GenerateRoom(roomgen.RoomConfig{
        Difficulty: roomgen.Medium,
        PartyLevel: partyLevel,
        DungeonTheme: "goblin_cave",
    })
}
```

### 2. Plugin System
```go
// Room generator as compiled plugin
type RoomGeneratorPlugin interface {
    Name() string
    Version() string
    Generate(config map[string]interface{}) ([]byte, error)
}

// Load at runtime
plugin, _ := plugin.Open("./plugins/roomgen.so")
symGenerator, _ := plugin.Lookup("Generator")
generator := symGenerator.(RoomGeneratorPlugin)
```

### 3. Configuration-Driven
```yaml
# room_templates.yaml
templates:
  goblin_ambush:
    description: "A narrow passage opens into a small chamber..."
    monsters:
      - type: "goblin"
        count: "1d4+1"
      - type: "wolf"
        count: "1d2"
        chance: 0.5
    loot:
      - table: "individual_treasure_cr_0_4"
        rolls: "1d4"
    features:
      - "crude_alarm_trap"
      - "hidden_escape_tunnel"
```

## Feature Ideas

### Smart Monster Selection
```go
type MonsterSelector struct {
    // Balance encounter based on party
    CalculateCR(party Party) float64
    
    // Themed groups that make sense together
    GetMonsterGroup(theme string) []Monster
    
    // Adjust difficulty dynamically
    ScaleEncounter(monsters []Monster, targetCR float64) []Monster
}
```

### Contextual Loot Generation
```go
type LootGenerator struct {
    // Loot makes sense for the monsters
    GenerateContextualLoot(monsters []Monster) []LootItem
    
    // Rarity based on difficulty
    RollRarity(difficultyLevel int) ItemRarity
    
    // Special loot tables
    RollMagicItem(tier string) MagicItem
}
```

### Environmental Storytelling
```go
type EnvironmentBuilder struct {
    // What happened here?
    GenerateRoomHistory() string
    
    // Interactive elements
    AddInteractables(roomType string) []Interactable
    
    // Clues and secrets
    HideSecrets(room *Room, difficultyDC int)
}
```

## Example Implementations

### Basic Combat Room
```go
func generateCombatRoom(level int) *Room {
    return &Room{
        Name: "Guard Chamber",
        Description: "Weapon racks line the walls of this chamber. A table with dice and coins suggests the guards were gambling.",
        Monsters: []Monster{
            {Name: "Bandit", Count: 2, CR: 0.125},
            {Name: "Bandit Captain", Count: 1, CR: 2},
        },
        Loot: []LootItem{
            {Name: "Gold pieces", Quantity: "3d6"},
            {Name: "Silvered dagger", Quantity: 1, Rarity: "uncommon"},
        },
        Features: []EnvironmentalFeature{
            {Name: "Weapon Rack", Description: "Contains basic weapons", Interactable: true},
            {Name: "Alarm Bell", Description: "Alerts nearby rooms if rung", TriggerEffect: "summon_reinforcements"},
        },
    }
}
```

### Dynamic Puzzle Room
```go
func generatePuzzleRoom(theme string) *Room {
    puzzles := map[string]Puzzle{
        "ancient_tomb": RiddleOfTheSphinx{},
        "wizard_tower": MagicWordDoor{},
        "dwarven_vault": CombinationLock{},
    }
    
    return &Room{
        Name: "The Sealed Door",
        Description: generatePuzzleDescription(theme),
        Puzzle: puzzles[theme],
        Reward: generatePuzzleReward(theme),
    }
}
```

## API Design

### REST Endpoints (if service-based)
```
POST /api/v1/rooms/generate
{
    "config": {
        "difficulty": "medium",
        "party_level": 5,
        "theme": "undead_crypt"
    }
}

GET /api/v1/rooms/templates
GET /api/v1/rooms/themes
GET /api/v1/rooms/validate-config
```

### Go Module Interface
```go
// Simple usage
room := roomgen.Quick(5) // Generate room for level 5 party

// Advanced usage
generator := roomgen.New(
    roomgen.WithTheme("elemental_chaos"),
    roomgen.WithMonsterSource(customMonsterDB),
    roomgen.WithLootTables(customLoot),
)

room := generator.Generate(config)
```

## Data Sources

1. **D&D 5e SRD Data**
   - Monster stats from existing dnd5e-api
   - Item tables
   - Spell effects for traps

2. **Custom Content**
   - User-created room templates
   - Community monster variants
   - Homebrew items

3. **Procedural Generation**
   - Room shape algorithms
   - Name generators
   - Description builders

## Extensibility

### Plugin Hooks
```go
type RoomGeneratorHook interface {
    BeforeGenerate(config RoomConfig) RoomConfig
    AfterGenerate(room *Room) *Room
    ValidateRoom(room *Room) error
}

// Allow custom logic injection
generator.RegisterHook("cursed_items", CursedItemHook{})
generator.RegisterHook("holiday_theme", HalloweenHook{})
```

### Theme Packs
```go
// Users can add theme packs
type ThemePack interface {
    Name() string
    RoomTemplates() []RoomTemplate
    MonsterGroups() []MonsterGroup
    LootTables() []LootTable
    Descriptions() DescriptionGenerator
}
```

## Next Steps

1. Define the core `Room` and `RoomConfig` structs
2. Create interface for monster/loot selection
3. Build basic generator with a few themes
4. Add plugin system for extensibility
5. Create REST API wrapper (optional)
6. Write comprehensive tests
7. Document theme creation guide

## Dungeon Progression System

### Boss Room Probability
```go
type DungeonProgression struct {
    BaseBossprobability    float64 // Starting chance (e.g., 0.05)
    ProbabilityIncrement   float64 // Increase per room (e.g., 0.1)
    GuaranteedBossRoom     int     // Force boss at room N (e.g., 10)
    
    // Modifiers
    ConsecutiveCombatPenalty float64 // Reduce boss chance after multiple combats
    RestRoomBonus           float64 // Increase boss chance after rest
}

func (dp *DungeonProgression) CalculateBossProbability(config RoomConfig) float64 {
    // Base calculation
    prob := dp.BaseBossprobability + (float64(config.RoomsCleared) * dp.ProbabilityIncrement)
    
    // Guaranteed boss room
    if config.RoomsCleared >= dp.GuaranteedBossRoom && !config.BossDefeated {
        return 1.0
    }
    
    // Cap at reasonable maximum before guaranteed
    return math.Min(prob, 0.8)
}
```

### Room Type Distribution
```go
func DetermineRoomType(config RoomConfig, progression DungeonProgression) string {
    bossProbability := progression.CalculateBossProbability(config)
    
    // Roll for boss room first
    if rand.Float64() < bossProbability && !config.BossDefeated {
        return "boss"
    }
    
    // Adjust probabilities based on dungeon depth
    weights := map[string]float64{
        "combat":   0.4 - (float64(config.DungeonDepth) * 0.02),
        "puzzle":   0.2 + (float64(config.DungeonDepth) * 0.01),
        "treasure": 0.2 + (float64(config.DungeonDepth) * 0.01),
        "rest":     0.1,
        "mixed":    0.1,
    }
    
    return weightedRandomChoice(weights)
}
```

### Progressive Difficulty Scaling
```go
type DifficultyScaler struct {
    BaseDifficulty     DifficultyLevel
    ScalingFactor      float64 // How much harder each room gets
    BossMultiplier     float64 // Boss room difficulty multiplier
}

func (ds *DifficultyScaler) AdjustRoomDifficulty(room *Room, depth int, isBoss bool) {
    if isBoss {
        room.ChallengeRating *= ds.BossMultiplier
        room.Monsters = ds.upgradeMonstersToBoss(room.Monsters)
    } else {
        // Gradual scaling
        scalingBonus := 1.0 + (float64(depth) * ds.ScalingFactor)
        room.ChallengeRating *= scalingBonus
    }
}
```

### Boss Room Templates
```go
type BossRoomTemplate struct {
    Name           string
    MinionsAllowed bool
    SpecialMechanics []BossMechanic
    LootMultiplier  float64
}

var bossTemplates = map[string]BossRoomTemplate{
    "goblin_cave": {
        Name: "Goblin King's Throne Room",
        MinionsAllowed: true,
        SpecialMechanics: []BossMechanic{
            CallReinforcements{cooldown: 3},
            RageMode{triggerHP: 0.5},
        },
        LootMultiplier: 3.0,
    },
    "undead_crypt": {
        Name: "Lich's Sanctum",
        MinionsAllowed: false,
        SpecialMechanics: []BossMechanic{
            PhaseShift{phases: 2},
            SummonUndead{maxSummons: 4},
        },
        LootMultiplier: 4.0,
    },
}
```

### Example Usage with Progression
```go
// In your Discord bot
type DungeonSession struct {
    Theme        string
    RoomsCleared int
    CurrentDepth int
    BossDefeated bool
    Progression  DungeonProgression
}

func (bot *Bot) generateNextRoom(session *DungeonSession) *Room {
    config := roomgen.RoomConfig{
        PartyLevel:   5,
        PartySize:    4,
        DungeonTheme: session.Theme,
        DungeonDepth: session.CurrentDepth,
        RoomsCleared: session.RoomsCleared,
        BossDefeated: session.BossDefeated,
    }
    
    // Generator handles progression internally
    room := generator.GenerateRoom(config)
    
    // Update session
    session.RoomsCleared++
    session.CurrentDepth++
    if room.IsBossRoom {
        session.BossDefeated = true
    }
    
    return room
}
```

### Foreshadowing and Hints
```go
func AddBossForeshadowing(room *Room, bossProximity float64) {
    if bossProximity > 0.6 {
        room.Description += " The air grows heavy with an oppressive presence."
        room.Features = append(room.Features, EnvironmentalFeature{
            Name: "Ominous Scratches",
            Description: "Deep claw marks score the walls, far larger than any creature you've encountered so far.",
        })
    }
    
    if bossProximity > 0.8 {
        room.Description += " You can hear distant roaring echoing through the corridors ahead."
    }
}
```

## Benefits of This Approach

- **Separation of Concerns**: Room generation logic separate from bot logic
- **Reusability**: Could be used by multiple bots or tools
- **Extensibility**: Easy to add new themes, monsters, items
- **Testability**: Generator can be tested independently
- **Community**: Others could contribute themes/content
- **Progressive Gameplay**: Natural dungeon progression with increasing challenge
- **Player Engagement**: Building tension towards inevitable boss encounter