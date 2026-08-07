-- State machine (old/scripts/event/OrbisPQ.js): 여신의 탑

local pq = require("script/lib/party_quest")

local EXIT_MAP = 920011200
local REWARD_MAP = 920011300
local START_MAP = 920010000
local BONUS_MAP = 920011100
local GARDEN_MAP = 920010800
local STAGE2_MAP = 920010300
local LOBBY_MAP = 920010400
local DURATION_MS = 3600000
local BONUS_MS = 60000
local RANKING_QUEST = 1203
local GUIDE_NPC = 2013001
local STAGE2_MOB = 9300040
local DARK_NEPENTHES = 9300049
local PAPA_PIXIE = 9300039
local STAGE2_PIECE = 4001045
local FIXED_REACTOR = 2001016

local STAGE_MAPS = {
	920010000,
	920010100,
	920010200,
	920010300,
	920010400,
	920010500,
	920010600,
	920010601,
	920010602,
	920010603,
	920010604,
	920010700,
	920010800,
	920010900,
	920010910,
	920010911,
	920010912,
	920010920,
	920010921,
	920010922,
	920010930,
	920010931,
	920010932,
	920011000,
	920011100,
}

local CX = { 200, -300, -300, -300, 200, 200, 200, -300, -300, 200, 200, -300, -300, 200 }
local CY = { -2321, -2114, -2910, -2510, -1526, -2716, -717, -1310, -3357, -1912, -1122, -1736, -915, -3116 }

local STAGE3_MUSIC = {
	[1] = "Bgm08/ForTheGlory",
	[2] = "Bgm11/Aquarium",
	[3] = "Bgm06/WelcomeToTheHell",
	[4] = "Bgm06/FantasticThinking",
	[5] = "Bgm02/EvilEyes",
	[6] = "Bgm10/TheWayGrotesque",
	[7] = "Bgm01/MoonlightShadow",
}

local function end_run(sm, map_id)
	if map_id == nil then
		map_id = EXIT_MAP
	end
	sm:finish(map_id)
end

local function clear_props(sm)
	sm:set_property("allfinish", "")
	sm:set_property("prestage", "")
	sm:set_property("prestageRewarded", "")
	sm:set_property("status", "0")
	sm:set_property("stage0clear", "")
	sm:set_property("stage1clear", "")
	sm:set_property("stage2clear", "")
	sm:set_property("stage2mob", "")
	sm:set_property("stage3_music", "")
	sm:set_property("stage3_item", "")
	sm:set_property("stage3clear", "")
	sm:set_property("stage4rand", "")
	sm:set_property("stage4try", "")
	sm:set_property("stage4clear", "")
	sm:set_property("stage4_rand", "")
	sm:set_property("stage4_r4way1", "")
	sm:set_property("stage4_r4way2", "")
	sm:set_property("stage5clear", "")
	sm:set_property("stage6clear", "")
	sm:set_property("stage6_ans", "")
	sm:set_property("stage6_way", "")
	sm:set_property("stagebossclear", "")
	sm:set_property("clearall", "")
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps(STAGE_MAPS)
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		clear_props(sm)
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
			end
		end
		local lobby = group:map(LOBBY_MAP)
		if lobby ~= nil then
			pq.shuffle_reactors(lobby, 2002004, 2002010)
		end
		local garden = group:map(GARDEN_MAP)
		if garden ~= nil then
			pq.shuffle_reactors(garden, nil, nil, FIXED_REACTOR)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(START_MAP)
		player:try_party_quest(RANKING_QUEST)
		player:open_npc(GUIDE_NPC)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id < 920010000 or map_id > 920011100 then
			sm:unregister(player)
			return
		end
		if map_id == BONUS_MAP then
			sm:set_property("allfinish", "1")
			sm:restart_timer(BONUS_MS)
			return
		end
		if map_id == LOBBY_MAP then
			local music_id = tonumber(sm:get_property("stage3_music")) or 0
			local path = STAGE3_MUSIC[music_id]
			if path ~= nil then
				local map = player:map()
				if map ~= nil then
					map:music(path)
				end
			end
		end
	end,

	on_mob_die = function(sm, mob)
		if mob == nil then
			return
		end
		local group = sm:group()
		if group == nil then
			return
		end
		local id = mob:id()
		if id == DARK_NEPENTHES then
			sm:notice("Boss Spawned.")
			local garden = group:map(GARDEN_MAP)
			if garden ~= nil then
				garden:spawn_mob(PAPA_PIXIE, -830, 563)
			end
			return
		end
		if id ~= STAGE2_MOB then
			return
		end
		local st = sm:get_property("stage2mob")
		if st == "" then
			sm:set_property("stage2mob", "0")
			st = "0"
		end
		local count = tonumber(st) or 0
		local map = group:map(STAGE2_MAP)
		if map == nil then
			return
		end
		if count < 14 then
			local idx = count + 1
			map:spawn_mob(STAGE2_MOB, CX[idx], CY[idx])
			sm:notice("Celion Spawned.", Msg.PinkText)
			sm:set_property("stage2mob", tostring(count + 1))
		else
			map:spawn_item(STAGE2_PIECE, 1, { CX[14], CY[14] })
		end
	end,

	on_left_party = function(sm, player)
		end_run(sm, EXIT_MAP)
	end,

	on_disband_party = function(sm)
		end_run(sm, EXIT_MAP)
	end,

	on_scheduled_timeout = function(sm)
		if sm:get_property("allfinish") == "1" then
			end_run(sm, REWARD_MAP)
		else
			end_run(sm, EXIT_MAP)
		end
	end,

	on_clear = function(sm)
		end_run(sm, EXIT_MAP)
	end,

	on_finish = function(sm)
		local group = sm:group()
		group:set_property("state", "0")
		clear_props(sm)
	end
}
