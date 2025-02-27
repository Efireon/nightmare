package item

import (
	"fmt"
	"math/rand"

	"nightmare/internal/entity"
	"nightmare/internal/util"
)

// ItemType представляет тип предмета
type ItemType int

const (
	ItemWeapon ItemType = iota
	ItemLight
	ItemMedical
	ItemFood
	ItemKey
	ItemNote
	ItemArtifact
	ItemRelic
	ItemMemento
	ItemMisc
)

// ItemRarity представляет редкость предмета
type ItemRarity int

const (
	RarityCommon ItemRarity = iota
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
	RarityUnique
)

// Item представляет предмет в игре
type Item struct {
	ID            int
	Name          string
	Description   string
	Type          ItemType
	Rarity        ItemRarity
	Weight        float64
	Value         int
	Stackable     bool
	Equippable    bool
	Consumable    bool
	Quantity      int
	Durability    float64
	MaxDurability float64
	Stats         map[string]float64
	Effects       []ItemEffect
	TextureID     int
	IconID        int
	UseSound      string
	PickupSound   string
	DropSound     string
	ExamineText   string
	Lore          string
	Tags          []string
}

// ItemEffect представляет эффект, применяемый предметом
type ItemEffect struct {
	Type        string
	Value       float64
	Duration    float64
	Probability float64
	Target      string
}

// NewItem создает новый предмет
func NewItem(id int, name, description string, itemType ItemType) *Item {
	return &Item{
		ID:            id,
		Name:          name,
		Description:   description,
		Type:          itemType,
		Rarity:        RarityCommon,
		Weight:        1.0,
		Value:         1,
		Stackable:     false,
		Equippable:    false,
		Consumable:    false,
		Quantity:      1,
		Durability:    100.0,
		MaxDurability: 100.0,
		Stats:         make(map[string]float64),
		Effects:       []ItemEffect{},
		TextureID:     0,
		IconID:        0,
		UseSound:      "",
		PickupSound:   "",
		DropSound:     "",
		ExamineText:   "",
		Lore:          "",
		Tags:          []string{},
	}
}

// Use использует предмет
func (i *Item) Use(user *entity.Player, target interface{}) bool {
	// Если предмет неиспользуемый, возвращаем false
	if !i.Consumable && !i.Equippable {
		return false
	}

	// Применяем эффекты предмета
	for _, effect := range i.Effects {
		// Проверяем вероятность срабатывания эффекта
		if rand.Float64() <= effect.Probability {
			i.applyEffect(effect, user, target)
		}
	}

	// Уменьшаем прочность, если предмет не расходуемый
	if !i.Consumable && i.Durability > 0 {
		i.Durability -= 1.0

		// Если прочность достигла нуля, предмет ломается
		if i.Durability <= 0 {
			i.Durability = 0
			// Здесь можно добавить дополнительную логику поломки предмета
		}
	}

	return true
}

// OnEquip вызывается при экипировке предмета
func (i *Item) OnEquip(user *entity.Player) {
	// Применяем постоянные эффекты предмета
	for _, effect := range i.Effects {
		if effect.Target == "equip" {
			i.applyEffect(effect, user, nil)
		}
	}
}

// OnUnequip вызывается при снятии предмета
func (i *Item) OnUnequip(user *entity.Player) {
	// Отменяем постоянные эффекты предмета
	for _, effect := range i.Effects {
		if effect.Target == "equip" {
			i.removeEffect(effect, user)
		}
	}
}

// applyEffect применяет эффект предмета
func (i *Item) applyEffect(effect ItemEffect, user *entity.Player, target interface{}) {
	switch effect.Type {
	case "heal":
		user.Health += effect.Value
		if user.Health > entity.MaxHealth {
			user.Health = entity.MaxHealth
		}

	case "sanity":
		user.Sanity += effect.Value
		if user.Sanity > entity.MaxSanity {
			user.Sanity = entity.MaxSanity
		}

	case "damage":
		if target != nil {
			if creature, ok := target.(*entity.Creature); ok {
				creature.TakeDamage(effect.Value)
			}
		}

	case "light":
		// Эффект освещения будет реализован в системе освещения

	case "speed":
		// Временное изменение скорости
		// Это потребует добавления системы временных эффектов

	case "reveal_map":
		// Раскрытие карты будет реализовано в системе карты

	case "scare":
		// Пугающий эффект для существ
		if target != nil {
			if creature, ok := target.(*entity.Creature); ok {
				creature.CurrentState = "flee"
				creature.StateTime = 0
			}
		}
	}
}

// removeEffect отменяет эффект предмета
func (i *Item) removeEffect(effect ItemEffect, user *entity.Player) {
	switch effect.Type {
	case "speed":
		// Отмена временного изменения скорости

	case "light":
		// Отмена эффекта освещения
	}
}

// Examine возвращает информацию о предмете при осмотре
func (i *Item) Examine() string {
	if i.ExamineText != "" {
		return i.ExamineText
	}

	return fmt.Sprintf("%s: %s", i.Name, i.Description)
}

// GetStats возвращает статистику предмета
func (i *Item) GetStats() map[string]float64 {
	return i.Stats
}

// HasTag проверяет, имеет ли предмет указанный тег
func (i *Item) HasTag(tag string) bool {
	for _, t := range i.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Copy создает копию предмета
func (i *Item) Copy() *Item {
	newItem := &Item{
		ID:            i.ID,
		Name:          i.Name,
		Description:   i.Description,
		Type:          i.Type,
		Rarity:        i.Rarity,
		Weight:        i.Weight,
		Value:         i.Value,
		Stackable:     i.Stackable,
		Equippable:    i.Equippable,
		Consumable:    i.Consumable,
		Quantity:      i.Quantity,
		Durability:    i.Durability,
		MaxDurability: i.MaxDurability,
		TextureID:     i.TextureID,
		IconID:        i.IconID,
		UseSound:      i.UseSound,
		PickupSound:   i.PickupSound,
		DropSound:     i.DropSound,
		ExamineText:   i.ExamineText,
		Lore:          i.Lore,
	}

	// Копируем статистику
	newItem.Stats = make(map[string]float64)
	for key, value := range i.Stats {
		newItem.Stats[key] = value
	}

	// Копируем эффекты
	newItem.Effects = make([]ItemEffect, len(i.Effects))
	copy(newItem.Effects, i.Effects)

	// Копируем теги
	newItem.Tags = make([]string, len(i.Tags))
	copy(newItem.Tags, i.Tags)

	return newItem
}

// ToString возвращает строковое представление предмета
func (i *Item) ToString() string {
	result := fmt.Sprintf("%s (ID: %d, Type: %s, Rarity: %s)\n",
		i.Name,
		i.ID,
		i.getTypeString(),
		i.getRarityString())

	result += fmt.Sprintf("Description: %s\n", i.Description)
	result += fmt.Sprintf("Weight: %.1f, Value: %d\n", i.Weight, i.Value)

	if i.Stackable {
		result += fmt.Sprintf("Quantity: %d\n", i.Quantity)
	}

	if i.Equippable {
		result += fmt.Sprintf("Durability: %.1f/%.1f\n", i.Durability, i.MaxDurability)
	}

	if len(i.Stats) > 0 {
		result += "Stats:\n"
		for stat, value := range i.Stats {
			result += fmt.Sprintf("- %s: %.1f\n", stat, value)
		}
	}

	if len(i.Effects) > 0 {
		result += "Effects:\n"
		for _, effect := range i.Effects {
			result += fmt.Sprintf("- %s: %.1f (%.1f%% chance, %.1fs duration)\n",
				effect.Type,
				effect.Value,
				effect.Probability*100,
				effect.Duration)
		}
	}

	if len(i.Tags) > 0 {
		result += fmt.Sprintf("Tags: %v\n", i.Tags)
	}

	return result
}

// getTypeString возвращает строковое представление типа предмета
func (i *Item) getTypeString() string {
	switch i.Type {
	case ItemWeapon:
		return "Weapon"
	case ItemLight:
		return "Light"
	case ItemMedical:
		return "Medical"
	case ItemFood:
		return "Food"
	case ItemKey:
		return "Key"
	case ItemNote:
		return "Note"
	case ItemArtifact:
		return "Artifact"
	case ItemRelic:
		return "Relic"
	case ItemMemento:
		return "Memento"
	case ItemMisc:
		return "Misc"
	default:
		return "Unknown"
	}
}

// getRarityString возвращает строковое представление редкости предмета
func (i *Item) getRarityString() string {
	switch i.Rarity {
	case RarityCommon:
		return "Common"
	case RarityUncommon:
		return "Uncommon"
	case RarityRare:
		return "Rare"
	case RarityEpic:
		return "Epic"
	case RarityLegendary:
		return "Legendary"
	case RarityUnique:
		return "Unique"
	default:
		return "Unknown"
	}
}

// ItemFactory фабрика для создания предметов
type ItemFactory struct {
	nextID  int
	random  *util.RandomGenerator
	itemsDB map[string]*Item // База шаблонов предметов
}

// NewItemFactory создает новую фабрику предметов
func NewItemFactory() *ItemFactory {
	return &ItemFactory{
		nextID:  1,
		random:  util.NewRandomGenerator(0),
		itemsDB: make(map[string]*Item),
	}
}

// RegisterItemTemplate регистрирует шаблон предмета в базе
func (f *ItemFactory) RegisterItemTemplate(templateID string, item *Item) {
	f.itemsDB[templateID] = item
}

// CreateItem создает новый предмет на основе шаблона
func (f *ItemFactory) CreateItem(templateID string) *Item {
	template, ok := f.itemsDB[templateID]
	if !ok {
		return nil
	}

	// Создаем копию шаблона
	item := template.Copy()

	// Устанавливаем уникальный ID
	item.ID = f.nextID
	f.nextID++

	return item
}

// CreateRandomItem создает случайный предмет указанного типа и редкости
func (f *ItemFactory) CreateRandomItem(itemType ItemType, rarity ItemRarity) *Item {
	// Собираем все шаблоны указанного типа и редкости
	var templates []*Item
	for _, item := range f.itemsDB {
		if item.Type == itemType && item.Rarity == rarity {
			templates = append(templates, item)
		}
	}

	// Если нет подходящих шаблонов, возвращаем nil
	if len(templates) == 0 {
		return nil
	}

	// Выбираем случайный шаблон
	template := templates[f.random.RangeInt(0, len(templates))]

	// Создаем копию шаблона
	item := template.Copy()

	// Устанавливаем уникальный ID
	item.ID = f.nextID
	f.nextID++

	// Добавляем случайную вариацию
	f.addRandomVariation(item)

	return item
}

// addRandomVariation добавляет случайную вариацию к предмету
func (f *ItemFactory) addRandomVariation(item *Item) {
	// Вариация прочности
	if item.Durability > 0 {
		variation := 0.8 + f.random.Float64()*0.4 // 80-120%
		item.Durability = item.MaxDurability * variation
		if item.Durability > item.MaxDurability {
			item.Durability = item.MaxDurability
		}
	}

	// Вариация статистик
	for stat, value := range item.Stats {
		variation := 0.9 + f.random.Float64()*0.2 // 90-110%
		item.Stats[stat] = value * variation
	}

	// Вариация эффектов
	for i := range item.Effects {
		variation := 0.9 + f.random.Float64()*0.2 // 90-110%
		item.Effects[i].Value *= variation
	}

	// Иногда добавляем дополнительный тег
	if f.random.Chance(0.1) {
		extraTags := []string{"Cursed", "Blessed", "Ancient", "Broken", "Repaired", "Modified"}
		item.Tags = append(item.Tags, extraTags[f.random.RangeInt(0, len(extraTags))])
	}
}

// CreateItemsDatabase создает базу данных предметов с шаблонами
func (f *ItemFactory) CreateItemsDatabase() {
	// Создаем шаблоны для оружия
	weapon := NewItem(0, "Rusty Pipe", "A rusty metal pipe. Not very effective, but better than nothing.", ItemWeapon)
	weapon.Equippable = true
	weapon.Weight = 2.0
	weapon.Stats["damage"] = 10.0
	weapon.Effects = append(weapon.Effects, ItemEffect{
		Type:        "damage",
		Value:       10.0,
		Probability: 1.0,
		Target:      "target",
	})
	f.RegisterItemTemplate("weapon_pipe", weapon)

	flashlight := NewItem(0, "Flashlight", "A battery-powered flashlight. Essential for navigating dark areas.", ItemLight)
	flashlight.Equippable = true
	flashlight.Weight = 0.5
	flashlight.Effects = append(flashlight.Effects, ItemEffect{
		Type:        "light",
		Value:       10.0,
		Probability: 1.0,
		Target:      "equip",
	})
	f.RegisterItemTemplate("light_flashlight", flashlight)

	medkit := NewItem(0, "First Aid Kit", "Contains basic medical supplies to treat wounds.", ItemMedical)
	medkit.Consumable = true
	medkit.Weight = 1.0
	medkit.Effects = append(medkit.Effects, ItemEffect{
		Type:        "heal",
		Value:       50.0,
		Probability: 1.0,
		Target:      "self",
	})
	f.RegisterItemTemplate("medical_first_aid", medkit)

	pills := NewItem(0, "Sanity Pills", "Helps stabilize your mental state.", ItemMedical)
	pills.Consumable = true
	pills.Stackable = true
	pills.Quantity = 3
	pills.Weight = 0.1
	pills.Effects = append(pills.Effects, ItemEffect{
		Type:        "sanity",
		Value:       20.0,
		Probability: 1.0,
		Target:      "self",
	})
	f.RegisterItemTemplate("medical_pills", pills)

	key := NewItem(0, "Rusty Key", "An old, rusty key. It might open something nearby.", ItemKey)
	key.Weight = 0.1
	key.Tags = append(key.Tags, "Key")
	f.RegisterItemTemplate("key_rusty", key)

	note := NewItem(0, "Torn Page", "A page torn from a journal. It contains disturbing writings.", ItemNote)
	note.Weight = 0.1
	note.ExamineText = "The page reads: 'They're watching me. I can feel their eyes following me through the trees. The forest itself seems alive, breathing with malice.'"
	f.RegisterItemTemplate("note_torn", note)

	artifact := NewItem(0, "Strange Amulet", "An amulet with unusual symbols. It gives off an unnerving aura.", ItemArtifact)
	artifact.Equippable = true
	artifact.Weight = 0.2
	artifact.Rarity = RarityRare
	artifact.Effects = append(artifact.Effects, ItemEffect{
		Type:        "sanity",
		Value:       -5.0,
		Probability: 1.0,
		Target:      "equip",
	})
	artifact.Effects = append(artifact.Effects, ItemEffect{
		Type:        "scare",
		Value:       10.0,
		Probability: 0.5,
		Target:      "creature",
	})
	f.RegisterItemTemplate("artifact_amulet", artifact)

	// И так далее для других типов предметов...
}
