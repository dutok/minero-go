package server

import (
	"github.com/toqueteos/minero/config"
)

func ConfigCreate() *config.Config {
	c := config.NewFrom(ConfigDefault())
	c.Save("./server.conf")
	return c
}

// ConfigDefault
func ConfigDefault() config.Map {
	return config.Map{
		"server.difficulty":          "1",
		"server.gamemode":            "0",
		"server.hardcore":            "false",
		"server.ip":                  "",
		"server.max_players":         "20",
		"server.motd":                "A minero server",
		"server.online_mode":         "true",
		"server.port":                "25565",
		"server.pvp":                 "true",
		"server.query.enable":        "false",
		"server.query.port":          "25565",
		"server.rcon.enable":         "false",
		"server.rcon.password":       "",
		"server.rcon.port":           "25575",
		"server.texture_pack":        "",
		"server.view_distance":       "10",
		"server.white_list":          "false",
		"spawn.animals":              "true",
		"spawn.monsters":             "true",
		"spawn.npcs":                 "true",
		"spawn.protection.enabled":   "true",
		"spawn.protection.shape":     "square",
		"spawn.protection.size":      "10",
		"worlds.allow_end":           "false",
		"worlds.allow_nether":        "false",
		"worlds.flight":              "false",
		"worlds.generate_structures": "true",
		"worlds.generator_settings":  "",
		"worlds.level_name":          "minero_test",
		"worlds.level_seed":          "0",
		"worlds.level_type":          "DEFAULT",
		"worlds.max_build_height":    "256",
		"worlds.snooper_enabled":     "true",
	}
}
