package material

const (
	Air                  Material = iota // 00
	Stone                                // 01
	Grass                                // 02
	Dirt                                 // 03
	Cobblestone                          // 04
	WoodPlanks                           // 05 // D
	Sapling                              // 06 // DB
	Bedrock                              // 07
	Water                                // 08 // D
	WaterStatic                          // 09 // D
	Lava                                 // 10 // D
	LavaStatic                           // 11 // D
	Sand                                 // 12
	Gravel                               // 13
	OreGold                              // 14
	OreIron                              // 15
	OreCoal                              // 16
	Wood                                 // 17 // DB
	Leaves                               // 18 // DB
	Sponge                               // 19
	Glass                                // 20
	OreLapisLazuli                       // 21
	LapisLazuli                          // 22
	Dispenser                            // 23 // DT
	Sandstone                            // 24 // D
	Note                                 // 25 // T
	Bed                                  // 26 // DI
	RailPowered                          // 27 // D
	RailDetector                         // 28 // D
	StickyPiston                         // 29 // D
	Cobweb                               // 30
	TallGrass                            // 31 // D
	DeadBush                             // 32
	Piston                               // 33 // D
	PistonExtension                      // 34 // D
	Wool                                 // 35 // DB
	PistonBlockMoved                     // 36 // T
	FlowerYellow                         // 37
	FlowerRed                            // 38
	MushroomBrown                        // 39
	MushroomRed                          // 40
	Gold                                 // 41
	Iron                                 // 42
	SlabDouble                           // 43 // DB
	Slab                                 // 44 // DB
	Brick                                // 45
	Tnt                                  // 46
	Bookshelf                            // 47
	MossStone                            // 48
	Obsidian                             // 49
	Torch                                // 50 // D
	Fire                                 // 51 // D
	MonsterSpawner                       // 52 // T
	StairsWoodOak                        // 53 // D
	Chest                                // 54 // DT
	RedstoneWire                         // 55 // DI
	OreDiamond                           // 56
	Diamond                              // 57
	CraftingTable                        // 58
	Wheat                                // 59 // D
	Farmland                             // 60 // D
	Furnace                              // 61 // DT
	FurnaceOn                            // 62 // DT
	SignPost                             // 63 // DIT
	DoorWood                             // 64 // DI
	Ladder                               // 65 // D
	Rail                                 // 66 // D
	StairsCobblestone                    // 67 // D
	SignWall                             // 68 // DT
	Lever                                // 69 // D
	PlateStone                           // 70 // D
	DoorIron                             // 71 // DI
	PlateWood                            // 72 // D
	OreRedstone                          // 73
	OreRedstoneOn                        // 74
	RedstoneTorch                        // 75 // D
	RedstoneTorchOn                      // 76 // D
	ButtonStone                          // 77 // D
	Snow                                 // 78 // D
	Ice                                  // 79
	BlockSnow                            // 80 // Snow is 1/8 of this
	Cactus                               // 81 // D
	Clay                                 // 82
	SugarCane                            // 83 // DI
	Jukebox                              // 84 // DT
	Fence                                // 85
	Pumpkin                              // 86 // D
	Netherrack                           // 87
	SoulSand                             // 88
	Glowstone                            // 89
	PortalNether                         // 90
	JackOLantern                         // 91 // D
	Cake                                 // 92 // DI
	RedstoneRepeater                     // 93 // DI
	RedstoneRepeaterOn                   // 94 // DI
	ChestLocked                          // 95
	Trapdoor                             // 96 // D
	EggMonster                           // 97 // D
	StoneBrick                           // 98 // DB
	MushroomBrownHuge                    // 99 // D
	MushroomRedHuge                      // 100 // D
	IronBars                             // 101
	GlassPane                            // 102
	Melon                                // 103
	PumpkinStem                          // 104 // D
	MelonStem                            // 105 // D
	Vine                                 // 106 // D
	FenceGate                            // 107 // D
	StairsBrick                          // 108 // D
	StairsStoneBrick                     // 109 // D
	Mycelium                             // 110
	LilyPad                              // 111
	NetherBrick                          // 112
	NetherBrickFence                     // 113
	StairsNetherBrick                    // 114 // D
	NetherWart                           // 115 // DI
	EnchantmentTable                     // 116 // T
	BrewingStand                         // 117 // DTI
	Cauldron                             // 118 // DI
	PortalEnd                            // 119 // T
	BlockPortalEnd                       // 120 // D
	EndStone                             // 121
	EggDragon                            // 122
	RedstoneLamp                         // 123
	RedstoneLampOn                       // 124
	SlabDoubleWood                       // 125 // DB
	SlabWood                             // 126 // DB
	Cocoa                                // 127 // DI
	StairsSandstone                      // 128 // D
	OreEmerald                           // 129
	ChestEnder                           // 130 // DT
	TripwireHook                         // 131 // D
	Tripwire                             // 132 // DI
	Emerald                              // 133
	StairsSpruceWood                     // 134 // D
	StairsBirchWood                      // 135 // D
	StairsJungleWood                     // 136 // D
	Command                              // 137 // T
	Beacon                               // 138 // T
	CobblestoneWall                      // 139 // BD
	FlowerPot                            // 140 // DI
	Carrot                               // 141 // D
	Potato                               // 142 // D
	ButtonWood                           // 143 // D
	Skull                                // 144 // DTI
	Anvil                                // 145 // D
	ChestTrapped                         // 146 // DT
	WeightedPlateLight                   // 147 // D
	WeightedPlateHeavy                   // 148 // D
	RedstoneComparator                   // 149 // I
	RedstoneComparatorOn                 // 150 // I
	DaylightSensor                       // 151
	Redstone                             // 152
	OreQuartz                            // 153
	Hopper                               // 154 // DT
	Quartz                               // 155 // D
	StairsQuartz                         // 156
	RailActivator                        // 157
	Dropper                              // 158 // DT
)
