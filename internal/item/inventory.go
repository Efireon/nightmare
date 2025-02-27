package item

import (
	"fmt"
	"math"

	"nightmare/internal/entity"
	"nightmare/internal/event"
)

// MaxInventoryWeight максимальный вес, который может нести игрок
const MaxInventoryWeight = 50.0

// Inventory представляет инвентарь игрока
type Inventory struct {
	Items        []*Item
	EquippedItem *Item
	Owner        *entity.Player
	MaxWeight    float64
	EventManager *event.EventManager
}

// NewInventory создает новый инвентарь
func NewInventory(owner *entity.Player, eventManager *event.EventManager) *Inventory {
	return &Inventory{
		Items:        make([]*Item, 0),
		EquippedItem: nil,
		Owner:        owner,
		MaxWeight:    MaxInventoryWeight,
		EventManager: eventManager,
	}
}

// AddItem добавляет предмет в инвентарь
func (inv *Inventory) AddItem(item *Item) bool {
	// Проверяем, не превышен ли вес
	if inv.GetTotalWeight()+item.Weight > inv.MaxWeight {
		return false
	}

	// Если предмет складируемый, пытаемся объединить с имеющимися
	if item.Stackable {
		for _, existingItem := range inv.Items {
			if existingItem.ID == item.ID {
				existingItem.Quantity += item.Quantity

				// Генерируем событие обновления инвентаря
				if inv.EventManager != nil {
					inv.EventManager.TriggerCustom(
						"inventory_updated",
						inv,
						existingItem,
						nil,
						existingItem.Quantity,
						nil,
					)
				}

				return true
			}
		}
	}

	// Добавляем новый предмет
	inv.Items = append(inv.Items, item)

	// Генерируем событие добавления предмета
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"item_added",
			inv,
			item,
			nil,
			item.Quantity,
			nil,
		)
	}

	return true
}

// RemoveItem удаляет предмет из инвентаря
func (inv *Inventory) RemoveItem(item *Item, quantity int) bool {
	for i, existingItem := range inv.Items {
		if existingItem.ID == item.ID {
			// Уменьшаем количество
			if existingItem.Stackable {
				if existingItem.Quantity >= quantity {
					existingItem.Quantity -= quantity

					// Если количество стало нулевым, удаляем предмет
					if existingItem.Quantity <= 0 {
						// Если предмет экипирован, снимаем его
						if inv.EquippedItem == existingItem {
							inv.UnequipItem()
						}

						// Удаляем предмет из списка
						inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
					}

					// Генерируем событие обновления инвентаря
					if inv.EventManager != nil {
						inv.EventManager.TriggerCustom(
							"inventory_updated",
							inv,
							existingItem,
							nil,
							existingItem.Quantity,
							nil,
						)
					}

					return true
				}
			} else {
				// Если предмет не складируемый, удаляем его
				if inv.EquippedItem == existingItem {
					inv.UnequipItem()
				}

				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)

				// Генерируем событие удаления предмета
				if inv.EventManager != nil {
					inv.EventManager.TriggerCustom(
						"item_removed",
						inv,
						existingItem,
						nil,
						1,
						nil,
					)
				}

				return true
			}
		}
	}

	return false
}

// GetItemByID находит предмет по ID
func (inv *Inventory) GetItemByID(id int) *Item {
	for _, item := range inv.Items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

// GetItemByName находит предмет по имени
func (inv *Inventory) GetItemByName(name string) *Item {
	for _, item := range inv.Items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

// GetTotalWeight возвращает общий вес инвентаря
func (inv *Inventory) GetTotalWeight() float64 {
	totalWeight := 0.0
	for _, item := range inv.Items {
		totalWeight += item.Weight * float64(item.Quantity)
	}
	return totalWeight
}

// GetItemsOfType возвращает все предметы указанного типа
func (inv *Inventory) GetItemsOfType(itemType ItemType) []*Item {
	var items []*Item
	for _, item := range inv.Items {
		if item.Type == itemType {
			items = append(items, item)
		}
	}
	return items
}

// HasItem проверяет, есть ли предмет в инвентаре
func (inv *Inventory) HasItem(id int, quantity int) bool {
	for _, item := range inv.Items {
		if item.ID == id {
			return item.Quantity >= quantity
		}
	}
	return false
}

// UseItem использует предмет
func (inv *Inventory) UseItem(item *Item, target interface{}) bool {
	if !inv.HasItem(item.ID, 1) {
		return false
	}

	// Используем предмет
	if item.Use(inv.Owner, target) {
		// Если предмет расходуемый, уменьшаем количество
		if item.Consumable {
			inv.RemoveItem(item, 1)
		}

		// Генерируем событие использования предмета
		if inv.EventManager != nil {
			inv.EventManager.TriggerCustom(
				"item_used",
				inv.Owner,
				item,
				target,
				1,
				nil,
			)
		}

		return true
	}

	return false
}

// EquipItem экипирует предмет
func (inv *Inventory) EquipItem(item *Item) bool {
	if !inv.HasItem(item.ID, 1) || !item.Equippable {
		return false
	}

	// Снимаем текущий экипированный предмет, если он есть
	if inv.EquippedItem != nil {
		inv.UnequipItem()
	}

	// Экипируем предмет
	inv.EquippedItem = item
	item.OnEquip(inv.Owner)

	// Генерируем событие экипировки предмета
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"item_equipped",
			inv.Owner,
			item,
			nil,
			1,
			nil,
		)
	}

	return true
}

// UnequipItem снимает экипированный предмет
func (inv *Inventory) UnequipItem() bool {
	if inv.EquippedItem == nil {
		return false
	}

	// Снимаем предмет
	item := inv.EquippedItem
	item.OnUnequip(inv.Owner)
	inv.EquippedItem = nil

	// Генерируем событие снятия предмета
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"item_unequipped",
			inv.Owner,
			item,
			nil,
			1,
			nil,
		)
	}

	return true
}

// GetEquippedItem возвращает экипированный предмет
func (inv *Inventory) GetEquippedItem() *Item {
	return inv.EquippedItem
}

// SortItems сортирует предметы по типу
func (inv *Inventory) SortItems() {
	// Сортируем предметы по типу
	for i := 0; i < len(inv.Items); i++ {
		for j := i + 1; j < len(inv.Items); j++ {
			if inv.Items[i].Type > inv.Items[j].Type {
				inv.Items[i], inv.Items[j] = inv.Items[j], inv.Items[i]
			}
		}
	}

	// Генерируем событие сортировки инвентаря
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"inventory_sorted",
			inv,
			nil,
			nil,
			0,
			nil,
		)
	}
}

// DropItem выбрасывает предмет из инвентаря
func (inv *Inventory) DropItem(item *Item, quantity int) bool {
	if !inv.HasItem(item.ID, quantity) {
		return false
	}

	// Снимаем предмет, если он экипирован
	if inv.EquippedItem == item {
		inv.UnequipItem()
	}

	// Удаляем предмет из инвентаря
	inv.RemoveItem(item, quantity)

	// Генерируем событие выбрасывания предмета
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"item_dropped",
			inv.Owner,
			item,
			inv.Owner.Position,
			quantity,
			nil,
		)
	}

	return true
}

// GetInventoryWeight возвращает текущий вес инвентаря и максимальный вес
func (inv *Inventory) GetInventoryWeight() (float64, float64) {
	return inv.GetTotalWeight(), inv.MaxWeight
}

// ClearInventory очищает инвентарь
func (inv *Inventory) ClearInventory() {
	// Снимаем экипированный предмет
	if inv.EquippedItem != nil {
		inv.UnequipItem()
	}

	// Очищаем инвентарь
	inv.Items = make([]*Item, 0)

	// Генерируем событие очистки инвентаря
	if inv.EventManager != nil {
		inv.EventManager.TriggerCustom(
			"inventory_cleared",
			inv,
			nil,
			nil,
			0,
			nil,
		)
	}
}

// ToString возвращает строковое представление инвентаря
func (inv *Inventory) ToString() string {
	result := fmt.Sprintf("Inventory: %d/%d items, %.1f/%.1f weight\n",
		len(inv.Items),
		math.MaxInt32, // Без ограничения количества предметов
		inv.GetTotalWeight(),
		inv.MaxWeight)

	if inv.EquippedItem != nil {
		result += fmt.Sprintf("Equipped: %s\n", inv.EquippedItem.Name)
	}

	result += "Items:\n"
	for _, item := range inv.Items {
		result += fmt.Sprintf("- %s (x%d): %.1f weight\n",
			item.Name,
			item.Quantity,
			item.Weight*float64(item.Quantity))
	}

	return result
}
