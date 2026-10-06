-- State machine: 몬스터 카니발

local HUB_MAP = 980000000
local ACCEPT_NPC = 2042001
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
	return carnival.map_match(waiting_id)
end

local function red_leader(sm)
	local match = match_for_sm(sm)
	local red = match ~= nil and match:red_team() or nil
	if red == nil then
		return nil, match
	end
	return red:leader(), match
end

local function offer_challenge(leader, match)
	if leader == nil or leader:in_dialog() or not match:has_pending_challenge() then
		return
	end
	leader:open_npc(ACCEPT_NPC)
end

local function warp_out(sm, match)
	if match == nil then
		sm:finish(HUB_MAP)
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
	sm:finish(0)
end

return {
	on_init = function(group)
		carnival.bind_group(group)
		carnival.set_skill_hit_chance(4, 100)
		group:exit_map(HUB_MAP)
		group:max_machines(nil)
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
		return {
			match:waiting_map_id(),
			match:revive_map_id(),
			match:field_map_id(),
		}
	end,

	on_prepare = function(sm)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		sm:map(match:field_map_id()):respawn({ include_one_time = true })
	end,

	on_player_enter = function(sm, player)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		player:map(sm:map(match:waiting_map_id()), 0)
	end,

	on_start = function(sm)
		sm:start_timer(WAITING_MS)
	end,

	on_challenge = function(sm)
		local leader, match = red_leader(sm)
		offer_challenge(leader, match)
	end,

	on_challenge_failed = function(sm)
		local leader, match = red_leader(sm)
		if leader == nil then
			return
		end
		leader:message("도전을 수락하는데 실패하였네.", Msg.PinkText)
		offer_challenge(leader, match)
	end,

	on_challenge_accepted = function(sm)
		sm:restart_timer(READY_MS)
	end,

	on_scheduled_timeout = function(sm)
		local match = match_for_sm(sm)
		if match == nil then
			sm:finish(HUB_MAP)
			return
		end
		local state = match:state()
		if state == CARNIVAL_STATE.WAITING then
			sm:finish(HUB_MAP)
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
			match:show_result()
			local field = sm:map(match:field_map_id())
			field:set_respawn(false)
			field:kill_all_mobs()
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
		player:hp(math.floor(player:max_hp() / 2))
		player:mp(math.floor(player:max_mp() / 2))
		match:on_player_died(player)
	end,

	on_carnival_summon = function(sm, player, tab, num)
		local match = match_for_sm(sm)
		if match == nil or match:state() ~= CARNIVAL_STATE.BATTLE then
			return
		end
		local team = match:find_team(player)
		local field = player:map()
		if team == nil or field == nil then
			return
		end

		local result
		if tab == CARNIVAL_TAB.MOB then
			result = team:summon_mob(player, field, num)
		elseif tab == CARNIVAL_TAB.SKILL then
			result = team:use_skill(player, field, num)
		elseif tab == CARNIVAL_TAB.GUARDIAN then
			result = team:summon_guardian(player, field, num)
		else
			return
		end

		if result == CARNIVAL_SUMMON_RESULT.SUCCESS then
			player:carnival_summon(tab, num)
		elseif result == CARNIVAL_SUMMON_RESULT.LACK_CP then
			player:message("CP가 부족합니다.", Msg.PinkText)
		elseif result == CARNIVAL_SUMMON_RESULT.NO_SLOT and tab == CARNIVAL_TAB.GUARDIAN then
			player:message("지금은 더 이상 수호물을 불러낼 수 없습니다.", Msg.PinkText)
		elseif result == CARNIVAL_SUMMON_RESULT.NO_SLOT then
			player:message("더 이상 소환수를 불러낼 수 없습니다.", Msg.PinkText)
		else
			player:message("오류가 발생했습니다.", Msg.PinkText)
		end
	end,

	on_mob_spawn = function(sm, mob)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		match:on_mob_spawn(mob)
	end,

	on_player_leave = function(sm, player, reason)
		if reason == "disconnect" then
			sm:finish(HUB_MAP)
		end
	end,

	on_left_party = function(sm, player)
		sm:finish(HUB_MAP)
	end,

	on_disband_party = function(sm)
		sm:finish(HUB_MAP)
	end,

	on_finish = function(sm)
		local match = match_for_sm(sm)
		if match == nil then
			return
		end
		if match:state() == CARNIVAL_STATE.REWARD then
			match:conclude()
		else
			match:finish()
		end
	end,
}
