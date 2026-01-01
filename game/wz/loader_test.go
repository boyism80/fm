package wz

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGeneralItems(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Etc", "0400.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadGeneralItems(path)
	if err != nil {
		t.Fatalf("Failed to load general items: %v", err)
	}

	if items == nil {
		t.Fatal("loadGeneralItems returned nil")
	}

	if len(*items) == 0 {
		t.Fatal("No items loaded")
	}

	// Test specific item (4000001) that we know should have slotMax = 200
	var foundItem *GeneralItem
	for _, item := range *items {
		if item.ID == 4000001 {
			foundItem = item
			break
		}
	}

	if foundItem == nil {
		t.Fatal("Item 4000001 not found")
	}

	if foundItem.SlotMax != 200 {
		t.Errorf("Item 4000001: expected slotMax=200, got %d", foundItem.SlotMax)
	}

	if foundItem.Price == 0 && foundItem.ID != 0 {
		t.Logf("Warning: Item %d has price=0 (might be expected)", foundItem.ID)
	}

	// Test that all items have valid IDs
	for _, item := range *items {
		if item.ID == 0 {
			t.Errorf("Found item with ID=0")
		}
		if item.ItemCore == nil {
			t.Errorf("Item %d has nil ItemCore", item.ID)
		}
	}

	t.Logf("Successfully loaded %d general items", len(*items))
}

func TestLoadConsumes(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Consume", "0200.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadConsumes(path)
	if err != nil {
		t.Fatalf("Failed to load consumes: %v", err)
	}

	if items == nil {
		t.Fatal("loadConsumes returned nil")
	}

	if len(*items) == 0 {
		t.Fatal("No items loaded")
	}

	// Test that items have valid data
	itemsWithSlotMax := 0
	for _, item := range *items {
		if item.ID == 0 {
			t.Errorf("Found item with ID=0")
		}
		if item.ItemCore == nil {
			t.Errorf("Item %d has nil ItemCore", item.ID)
		}
		if item.SlotMax > 0 {
			itemsWithSlotMax++
		}
	}

	t.Logf("Successfully loaded %d consume items (%d with slotMax > 0)", len(*items), itemsWithSlotMax)
}

func TestLoadCashItems(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Cash", "0501.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadCashItems(path)
	if err != nil {
		t.Fatalf("Failed to load cash items: %v", err)
	}

	if items == nil {
		t.Fatal("loadCashItems returned nil")
	}

	if len(*items) == 0 {
		t.Fatal("No items loaded")
	}

	// Test that items have valid data
	for _, item := range *items {
		if item.ID == 0 {
			t.Errorf("Found item with ID=0")
		}
		if item.ItemCore == nil {
			t.Errorf("Item %d has nil ItemCore", item.ID)
		}
	}

	t.Logf("Successfully loaded %d cash items", len(*items))
}

func TestLoadInstallations(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Install", "0301.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadInstallations(path)
	if err != nil {
		t.Fatalf("Failed to load installations: %v", err)
	}

	if items == nil {
		t.Fatal("loadInstallations returned nil")
	}

	if len(*items) == 0 {
		t.Fatal("No items loaded")
	}

	// Test that items have valid data
	itemsWithSlotMax := 0
	for _, item := range *items {
		if item.ID == 0 {
			t.Errorf("Found item with ID=0")
		}
		if item.ItemCore == nil {
			t.Errorf("Item %d has nil ItemCore", item.ID)
		}
		if item.SlotMax > 0 {
			itemsWithSlotMax++
		}
	}

	t.Logf("Successfully loaded %d installation items (%d with slotMax > 0)", len(*items), itemsWithSlotMax)
}

func TestLoadWeapons(t *testing.T) {
	// Find a weapon file
	weaponDir := filepath.Join("..", "..", "resources", "wz", "Character.wz", "Weapon")
	if _, err := os.Stat(weaponDir); os.IsNotExist(err) {
		t.Skipf("Weapon directory not found: %s", weaponDir)
		return
	}

	// Try to find any weapon file
	entries, err := os.ReadDir(weaponDir)
	if err != nil {
		t.Skipf("Failed to read weapon directory: %v", err)
		return
	}

	var weaponPath string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".xml" {
			weaponPath = filepath.Join(weaponDir, entry.Name())
			break
		}
	}

	if weaponPath == "" {
		t.Skip("No weapon XML files found")
		return
	}

	weapon, err := loadWeapons(weaponPath)
	if err != nil {
		t.Fatalf("Failed to load weapon: %v", err)
	}

	// loadWeapons can return nil for category files (hit, bow, axe, etc.)
	if weapon == nil {
		t.Log("loadWeapons returned nil (likely a category file, skipping)")
		return
	}

	if weapon.ID == 0 {
		t.Errorf("Weapon has ID=0")
	}
	if weapon.ItemCore == nil {
		t.Errorf("Weapon has nil ItemCore")
	}

	// slotMax is optional for weapons (many weapons don't have it in XML)
	// We only verify that if slotMax is present in XML, it's loaded correctly
	// This is tested by checking that weapons with slotMax > 0 exist in other test files

	t.Logf("Successfully loaded weapon ID=%d, slotMax=%d, price=%d", weapon.ID, weapon.SlotMax, weapon.Price)
}

func TestLoadGeneralItemsSlotMax(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Etc", "0400.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadGeneralItems(path)
	if err != nil {
		t.Fatalf("Failed to load general items: %v", err)
	}

	// Test multiple items with known slotMax values
	testCases := []struct {
		itemID   uint32
		slotMax  uint16
		hasPrice bool
	}{
		{4000001, 200, true},
		{4000000, 200, true},
	}

	for _, tc := range testCases {
		var foundItem *GeneralItem
		for _, item := range *items {
			if item.ID == tc.itemID {
				foundItem = item
				break
			}
		}

		if foundItem == nil {
			t.Logf("Item %d not found (might not exist in test file)", tc.itemID)
			continue
		}

		if foundItem.SlotMax != tc.slotMax {
			t.Errorf("Item %d: expected slotMax=%d, got %d", tc.itemID, tc.slotMax, foundItem.SlotMax)
		}

		if tc.hasPrice && foundItem.Price == 0 {
			t.Logf("Warning: Item %d has price=0 (might be expected)", tc.itemID)
		}
	}
}

func TestLoadConsumesPriceAndSlotMax(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Item.wz", "Consume", "0200.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	items, err := loadConsumes(path)
	if err != nil {
		t.Fatalf("Failed to load consumes: %v", err)
	}

	// Test that price and slotMax are loaded from Ints
	hasPrice := false
	hasSlotMax := false

	for _, item := range *items {
		if item.Price > 0 {
			hasPrice = true
		}
		if item.SlotMax > 0 {
			hasSlotMax = true
		}
	}

	if !hasPrice {
		t.Error("No items with price > 0 found (price should be loaded from Ints)")
	}

	if !hasSlotMax {
		t.Error("No items with slotMax > 0 found (slotMax should be loaded from Ints)")
	}

	t.Logf("Verified: price and slotMax are loaded correctly from Ints")
}
