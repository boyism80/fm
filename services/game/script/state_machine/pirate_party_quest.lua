-- State machine (old/scripts/event/Pirate.js): 해적 데비존

local pq = require("script/lib/party_quest")

local EXIT_MAP = 925100700
local START_MAP = 925100000
local BOSS_MAP = 925100500
local RANKING_QUEST = 1204
local TIMER_STAGE1_MS = 480000
local TIMER_DEFAULT_MS = 1080000
local TIMER_HD_MS = 180000
local TIMER_BOSS_MS = 480000
local TREASURE_REACTOR = 2512001
local WYANG_REACTOR = 2516000

local STAGE_MAPS = {
	925100000,
	925100100,
	925100200,
	925100201,
	925100202,
	925100300,
	925100301,
	925100302,
	925100400,
	925100500,
}

local function end_run(sm)
	sm:finish(EXIT_MAP)
end

local function treasure_state(sm, map_id)
	local group = sm:group()
	if group == nil then
		return 0
	end
	local map = group:map(map_id)
	if map == nil then
		return 0
	end
	local reactor = map:reactor(TREASURE_REACTOR)
	if reactor == nil then
		return 0
	end
	return reactor:state()
end

local function spawn_deck_pirates(map)
	if map == nil then
		return
	end
	for _ = 1, 5 do
		map:spawn_mob(9300124, 430, 75)
		map:spawn_mob(9300125, 1600, 75)
		map:spawn_mob(9300124, 430, 238)
		map:spawn_mob(9300125, 1600, 238)
	end
end

local function spawn_treasure_guards(map)
	if map == nil then
		return
	end
	for _ = 1, 10 do
		map:spawn_mob(9300112, 0, 238)
		map:spawn_mob(9300113, 1700, 238)
	end
end

local function open_treasure_if_empty(map)
	if map == nil or pq.mob_count(map) > 0 then
		return
	end
	local reactor = map:reactor(TREASURE_REACTOR)
	if reactor ~= nil and reactor:state() == 0 then
		reactor:hit(1)
	end
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
		sm:set_property("stage", "1")
		sm:set_property("stage2", "0")
		sm:set_property("stage2a", "0")
		sm:set_property("stage3a", "0")
		sm:set_property("stage4", "0")
		sm:set_property("clearstage", "")
		sm:set_property("entered_100", "")
		sm:set_property("entered_200", "")
		sm:set_property("entered_300", "")
		sm:set_property("entered_400", "")
		sm:set_property("entered_500", "")
		sm:set_property("hd_202", "")
		sm:set_property("hd_302", "")
		sm:set_property("hd_202_left", "")
		sm:set_property("hd_302_left", "")
		sm:set_property("hd_202_out", "")
		sm:set_property("hd_302_out", "")

		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
			end
		end

		local bow = group:map(925100100)
		if bow ~= nil then
			bow:block_gen(false, 9300114)
			bow:block_gen(false, 9300115)
			bow:block_gen(false, 9300116)
			bow:kill_all_mobs()
		end

		spawn_deck_pirates(group:map(925100200))
		spawn_treasure_guards(group:map(925100201))
		spawn_deck_pirates(group:map(925100300))
		spawn_treasure_guards(group:map(925100301))
	end,

	on_start = function(sm)
		sm:start_timer(TIMER_STAGE1_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(START_MAP)
		player:try_party_quest(RANKING_QUEST)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id < 925100000 or map_id > 925100500 then
			sm:unregister(player)
			return
		end
		if map_id == 925100100 then
			if sm:get_property("entered_100") == "" then
				sm:restart_timer(TIMER_DEFAULT_MS)
				sm:set_property("entered_100", "1")
				sm:notice("구옹이 모험가들에게 할 말이 있는 것 같습니다.")
			end
		elseif map_id == 925100200 then
			if sm:get_property("entered_200") == "" then
				sm:restart_timer(TIMER_DEFAULT_MS)
				sm:set_property("entered_200", "1")
				sm:notice("해적들이 외부인들을 경계합니다. 해적들을 모두 쓰러뜨리십시오!")
			elseif sm:get_property("hd_202_out") == "1" then
				sm:set_property("hd_202_out", "")
				local left = tonumber(sm:get_property("hd_202_left")) or TIMER_DEFAULT_MS
				pq.party_warp(sm, 925100200, player:id())
				sm:restart_timer(left)
			end
		elseif map_id == 925100202 then
			if sm:get_property("hd_202") == "" then
				sm:set_property("hd_202_left", tostring(sm:time_left()))
				sm:restart_timer(TIMER_HD_MS)
				sm:set_property("hd_202", "1")
				sm:notice("해적들이 숨어있는 것 같습니다..")
			end
		elseif map_id == 925100300 then
			if sm:get_property("entered_300") == "" then
				sm:restart_timer(TIMER_DEFAULT_MS)
				sm:set_property("entered_300", "1")
				sm:notice("해적들이 외부인들을 경계합니다. 해적들을 모두 쓰러뜨리십시오!")
			elseif sm:get_property("hd_302_out") == "1" then
				sm:set_property("hd_302_out", "")
				local left = tonumber(sm:get_property("hd_302_left")) or TIMER_DEFAULT_MS
				pq.party_warp(sm, 925100300, player:id())
				sm:restart_timer(left)
			end
		elseif map_id == 925100302 then
			if sm:get_property("hd_302") == "" then
				sm:set_property("hd_302_left", tostring(sm:time_left()))
				sm:restart_timer(TIMER_HD_MS)
				sm:set_property("hd_302", "1")
				sm:notice("해적들이 숨어있는 것 같습니다..")
			end
		elseif map_id == 925100400 then
			if sm:get_property("entered_400") == "" then
				sm:restart_timer(TIMER_DEFAULT_MS)
				sm:set_property("entered_400", "1")
				sm:notice("열쇠를 획득하여 갑판의 문을 모두 잠가버리고, 해적들이 더 나오지 않게 하세요!")
			end
		elseif map_id == 925100500 then
			if sm:get_property("entered_500") == "" then
				sm:restart_timer(TIMER_BOSS_MS)
				sm:set_property("entered_500", "1")
				local t1 = treasure_state(sm, 925100201)
				local t2 = treasure_state(sm, 925100301)
				local group = sm:group()
				local boss_map = group ~= nil and group:map(BOSS_MAP) or nil
				if boss_map ~= nil then
					if t1 == 2 and t2 == 2 then
						sm:notice("데비존이 화가 무척이나 나 있습니다! 주의하세요!")
						boss_map:spawn_mob(9300106, 630, 213)
					elseif t1 == 2 or t2 == 2 then
						sm:notice("데비존이 화가 나 있습니다! 주의하세요!")
						boss_map:spawn_mob(9300105, 630, 213)
					else
						sm:notice("모든 일의 원흉, 해적왕을 물리쳐야 합니다!")
						boss_map:spawn_mob(9300119, 630, 213)
					end
				end
			end
		end
	end,

	on_mob_die = function(sm, mob)
		local group = sm:group()
		if group == nil then
			return
		end
		open_treasure_if_empty(group:map(925100201))
		open_treasure_if_empty(group:map(925100301))
		if mob == nil then
			return
		end
		local id = mob:id()
		if id == 9300119 or id == 9300105 or id == 9300106 then
			local boss_map = group:map(BOSS_MAP)
			if boss_map ~= nil then
				local reactor = boss_map:reactor(WYANG_REACTOR)
				if reactor ~= nil then
					reactor:hit(1)
				end
			end
		end
	end,

	on_left_party = function(sm, player)
		end_run(sm)
	end,

	on_disband_party = function(sm)
		end_run(sm)
	end,

	on_scheduled_timeout = function(sm)
		local players = sm:players()
		local map_id = 0
		if #players > 0 and players[1] ~= nil then
			local map = players[1]:map()
			if map ~= nil and map:wz() ~= nil then
				map_id = map:wz().id
			end
		end
		if sm:get_property("hd_202") == "1" and map_id == 925100202 then
			sm:set_property("hd_202_out", "1")
			pq.party_warp(sm, 925100200)
			local left = tonumber(sm:get_property("hd_202_left")) or TIMER_DEFAULT_MS
			sm:restart_timer(left)
			return
		end
		if sm:get_property("hd_302") == "1" and map_id == 925100302 then
			sm:set_property("hd_302_out", "1")
			pq.party_warp(sm, 925100300)
			local left = tonumber(sm:get_property("hd_302_left")) or TIMER_DEFAULT_MS
			sm:restart_timer(left)
			return
		end
		end_run(sm)
	end,

	on_clear = function(sm)
		end_run(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}
