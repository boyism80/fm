-- State machine (old/scripts/event/cpq.js): 몬스터 카니발

local HUB_MAP = 980000000
local RANKING_QUEST = 1301
local WAITING_MS = 180000
local READY_MS = 10000
local BATTLE_MS = 600000
local REWARD_MS = 10000

local function register_slot(slot, waiting, max_members)
	carnival.register(
		slot,
		waiting,
		waiting + 1,
		waiting + 2,
		waiting + 3,
		waiting + 4,
		max_members
	)
end

local function match_for_sm(sm)
	local waiting_id = tonumber(sm:id())
	if waiting_id == nil then
		return nil
	end
	return carnival.match_by_map(waiting_id)
end

local function dispose_all(sm, match)
	if match ~= nil then
		match:finish()
	end
	sm:finish(HUB_MAP)
end

local function warp_out(sm, match)
	if match == nil then
		dispose_all(sm, match)
		return
	end
	local blue = match:blue_team()
	local red = match:red_team()
	if blue ~= nil and blue:is_winner() then
		blue:warp(match:win_map_id(), "sp")
		if red ~= nil then
			red:warp(match:lose_map_id(), "sp")
		end
	elseif red ~= nil and red:is_winner() then
		red:warp(match:win_map_id(), "sp")
		if blue ~= nil then
			blue:warp(match:lose_map_id(), "sp")
		end
	else
		if red ~= nil then
			red:warp(match:win_map_id(), "sp")
		end
		if blue ~= nil then
			blue:warp(match:lose_map_id(), "sp")
		end
	end
	match:finish()
end

return {
	on_init = function(group)
		carnival.bind_group(group)
		carnival.set_skill_hit_chance(4, 100)
		group:declare_exit_map(HUB_MAP)
		register_slot(0, 980000100, 2)
		register_slot(1, 980000200, 2)
		register_slot(2, 980000300, 2)
		register_slot(3, 980000400, 2)
		register_slot(4, 980000500, 3)
		register_slot(5, 980000600, 3)
	end,

	on_create = function(sm)
		local match = match_for_sm(sm)
		if match == nil then
			return { tonumber(sm:id()) }
		end
		local group = sm:group()
		local maps = {
			match:waiting_map_id(),
			match:revive_map_id(),
			match:field_map_id(),
			match:win_map_id(),
			match:lose_map_id(),
		}
		for _, map_id in ipairs(maps) do
			local map = group:map(map_id)
			if map ~= nil and map_id == match:field_map_id() then
				map:reset()
				map:respawn(true)
			end
		end
		return maps
	end,

	on_player_enter = function(sm, player)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		local waiting = sm:group():map(match:waiting_map_id())
		if waiting ~= nil then
			player:map(waiting, 0)
		end
	end,

	on_start = function(sm)
		sm:start_timer(WAITING_MS)
	end,

	on_challenge_accepted = function(sm)
		sm:restart_timer(READY_MS)
	end,

	on_scheduled_timeout = function(sm)
		local match = match_for_sm(sm)
		if match == nil then
			dispose_all(sm, match)
			return
		end
		local state = match:state()
		if state == CARNIVAL_STATE.WAITING then
			dispose_all(sm, match)
		elseif state == CARNIVAL_STATE.READY then
			match:set_state(CARNIVAL_STATE.BATTLE)
			local field_id = match:field_map_id()
			local blue = match:blue_team()
			local red = match:red_team()
			if blue ~= nil then
				blue:warp(field_id, "blue00")
			end
			if red ~= nil then
				red:warp(field_id, "red00")
			end
			for _, p in ipairs(sm:players()) do
				p:try_party_quest(RANKING_QUEST)
			end
			sm:restart_timer(BATTLE_MS)
		elseif state == CARNIVAL_STATE.BATTLE then
			local result = match:result()
			if result == CARNIVAL_RESULT.RED_WIN then
				match:red_team():set_winner(true)
			elseif result == CARNIVAL_RESULT.BLUE_WIN then
				match:blue_team():set_winner(true)
			end
			match:set_state(CARNIVAL_STATE.REWARD)
			sm:restart_timer(REWARD_MS)
		elseif state == CARNIVAL_STATE.REWARD then
			warp_out(sm, match)
		end
	end,

	on_player_revive = function(sm, player)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		local max_hp = player:max_hp()
		local max_mp = player:max_mp()
		player:exchange({}, { hp = math.floor(max_hp / 2), mp = math.floor(max_mp / 2) })
		local revive = sm:group():map(match:revive_map_id())
		if revive ~= nil then
			player:map(revive, 0)
		end
		match:on_player_died(player)
	end,

	on_player_disconnected = function(sm, player)
		dispose_all(sm, match_for_sm(sm))
	end,

	on_left_party = function(sm, player)
		dispose_all(sm, match_for_sm(sm))
	end,

	on_disband_party = function(sm)
		dispose_all(sm, match_for_sm(sm))
	end,
}
