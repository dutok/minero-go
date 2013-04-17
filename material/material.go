package material

type Material uint

const (
	MinBlock = 0
	MaxBlock = 158
	MinItem  = 256
	MaxItem  = 408
	MinDisc  = 2256
	MaxDisc  = 2267
)

func (m Material) String() string {
	return materials[m]
}

func (m Material) MaxStack() byte {
	return stacks[m]
}

func (m Material) MaxDurability() int {
	return durability[m]
}

// FromId attempts to match Material by its id.
// Matching order:
//   1: if id < MinBlock { return Air }
//   2: exact match
//   3: "Block" + name
//   4: "Item" + name
//   5: "Disc" + name
//   6: return Air
func FromId(id int) (m Material) {
	if id < MinBlock {
		return Air
	}

	m = Material(id)
	if m.IsBlock() || m.IsItem() || m.IsDisc() {
		return
	}
	return Air
}

// FromName attempts to match Material by its name.
// Matching order:
//   1: if name == "" { return Air }
//   2: exact match
//   3: "Item" + name
//   4: "Disc" + name
//   5: return Air
func FromName(name string) (m Material) {
	if name == "" {
		return Air
	}

	var ok bool
	if m, ok = names[name]; ok {
		return
	}
	if m, ok = names["Item"+name]; ok {
		return
	}
	if m, ok = names["Disc"+name]; ok {
		return
	}

	return Air
}

// IsBlock returns true if m is a block.
func (m Material) IsBlock() bool {
	return m >= MinBlock && m <= MaxBlock
}

// IsItem returns true if m is an item.
func (m Material) IsItem() bool {
	if m >= MinItem && m <= MaxItem {
		return true
	}
	if m.IsDisc() {
		return true
	}
	return false
}

// IsDisc returns true if m is a playable music disc.
func (m Material) IsDisc() bool {
	return m >= MinDisc && m <= MaxDisc
}

func (m Material) IsEdible() bool {
	switch m {
	case Bread, CarrotItem, PotatoBaked, PotatoItem, PotatoPoisonous,
		CarrotGolden, PumpkinPie, Cookie, MelonItem, MushroomStew, RawChicken,
		CookedChicken, RawSteak, CookedSteak, RawFish, CookedFish, RawPorkchop,
		CookedPorkchop, AppleItem, AppleGolden, RottenFlesh, EyeSpider:
		return true
	}
	return false
}

// IsSolid returns true if m is solid block.
func (m Material) IsSolid() bool {
	if !m.IsBlock() || m == Air {
		return false
	}
	switch m {
	case Stone, Grass, Dirt, Cobblestone, WoodPlanks, Bedrock, Sand, Gravel,
		OreGold, OreIron, OreCoal, Wood, Leaves, Sponge, Glass, OreLapisLazuli,
		LapisLazuli, Dispenser, Sandstone, Note, Bed, Piston, PistonExtension,
		Wool, PistonBlockMoved, Gold, Iron, SlabDouble, Slab, Brick, Tnt,
		Bookshelf, MossStone, Obsidian, MonsterSpawner, StairsWoodOak, Chest,
		OreDiamond, Diamond, CraftingTable, Wheat, Furnace, FurnaceOn, SignPost,
		DoorWood, StairsCobblestone, SignWall, PlateStone, DoorIron, PlateWood,
		OreRedstone, OreRedstoneOn, Ice, BlockSnow, Cactus, Clay, Jukebox,
		Fence, Pumpkin, Netherrack, SoulSand, Glowstone, JackOLantern, Cake,
		ChestLocked, Trapdoor, EggMonster, StoneBrick, MushroomBrownHuge,
		MushroomRedHuge, IronBars, GlassPane, Melon, FenceGate, StairsBrick,
		StairsStoneBrick, Mycelium, NetherBrick, NetherBrickFence,
		StairsNetherBrick, EnchantmentTable, BrewingStand, Cauldron, PortalEnd,
		EndStone, EggDragon, RedstoneLamp, RedstoneLampOn, SlabDoubleWood,
		SlabWood, StairsSandstone, OreEmerald, ChestEnder, Emerald,
		StairsSpruceWood, StairsBirchWood, StairsJungleWood, Command, Beacon,
		CobblestoneWall, Anvil, ChestTrapped, WeightedPlateLight,
		WeightedPlateHeavy, DaylightSensor, Redstone, OreQuartz, Hopper, Quartz,
		StairsQuartz, Dropper:
		return true
	}
	return false
}

// IsTransparent returns true if m is a block and does not block any light.
func (m Material) IsTransparent() bool {
	if !m.IsBlock() {
		return false
	}
	switch m {
	case Air, Sapling, RailPowered, RailDetector, TallGrass, DeadBush,
		FlowerYellow, FlowerRed, MushroomBrown, MushroomRed, Torch, Fire,
		RedstoneWire, Wheat, Ladder, Rail, Lever, RedstoneTorch,
		RedstoneTorchOn, ButtonStone, Snow, SugarCane, PortalNether,
		RedstoneRepeater, RedstoneRepeaterOn, PumpkinStem, MelonStem, Vine,
		LilyPad, NetherWart, PortalEnd, Cocoa, TripwireHook, Tripwire,
		FlowerPot, Carrot, Potato, ButtonWood, Skull, RedstoneComparator,
		RedstoneComparatorOn, RailActivator:
		return true
	}
	return false
}

// IsFlammable returns true if m is a block and can catch fire.
func (m Material) IsFlammable() bool {
	if !m.IsBlock() {
		return false
	}
	switch m {
	case WoodPlanks, Wood, Leaves, Note, Bed, TallGrass, DeadBush, Wool, Tnt,
		Bookshelf, StairsWoodOak, Chest, CraftingTable, SignPost, DoorWood,
		SignWall, PlateWood, Jukebox, Fence, ChestLocked, Trapdoor,
		MushroomBrownHuge, MushroomRedHuge, Vine, FenceGate, SlabDoubleWood,
		SlabWood, StairsSpruceWood, StairsBirchWood, StairsJungleWood,
		ChestTrapped, DaylightSensor:
		return true
	}
	return false
}

// IsBurnable returns true if m is a block and can burn away.
func (m Material) IsBurnable() bool {
	if !m.IsBlock() {
		return false
	}
	switch m {
	case WoodPlanks, Wood, Leaves, TallGrass, Wool, Tnt, Bookshelf,
		StairsWoodOak, Fence, Vine, SlabDoubleWood, SlabWood, StairsSpruceWood,
		StairsBirchWood, StairsJungleWood:
		return true
	}
	return false
}

// IsOccluding returns true if m is a block and completely blocks vision.
func (m Material) IsOccluding() bool {
	if !m.IsBlock() {
		return false
	}
	switch m {
	case Stone, Grass, Dirt, Cobblestone, WoodPlanks, Bedrock, Sand, Gravel,
		OreGold, OreIron, OreCoal, Wood, Sponge, OreLapisLazuli, LapisLazuli,
		Dispenser, Sandstone, Note, Wool, Gold, Iron, SlabDouble, Brick,
		Bookshelf, MossStone, Obsidian, MonsterSpawner, OreDiamond, Diamond,
		CraftingTable, Furnace, FurnaceOn, OreRedstone, OreRedstoneOn,
		BlockSnow, Clay, Jukebox, Pumpkin, Netherrack, SoulSand, JackOLantern,
		ChestLocked, EggMonster, StoneBrick, MushroomBrownHuge, MushroomRedHuge,
		Melon, Mycelium, NetherBrick, PortalEnd, EndStone, RedstoneLamp,
		RedstoneLampOn, SlabDoubleWood, OreEmerald, Emerald, Command, OreQuartz,
		Quartz, Dropper:
		return true
	}
	return false
}

// HasGravity returns true if m is affected by gravity.
func (m Material) HasGravity() bool {
	if !m.IsBlock() {
		return false
	}
	switch m {
	case Sand, Gravel, Anvil:
		return true
	}
	return false
}
