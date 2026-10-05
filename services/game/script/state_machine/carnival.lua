-- State machine: 몬스터 카니발

local HUB_MAP = 980000000
local ACCEPT_NPC = 2042001
local RANKING_QUEST = 1301
local WAITING_MS = 180000
local READY_MS = 10000
local BATTLE_MS = 600000
local REWARD_MS = 10000
local GUARDIAN_SPAWN_STATE = 1
local GUARDIAN_DESTROYED_STATE = 5

local GUARDIAN_BUFFS = {
	[140] = MobBuff.WeaponImmunity,
	[141] = MobBuff.MagicImmunity,
	[150] = MobBuff.WeaponAttackUp,
	[151] = MobBuff.MagicAttackUp,
	[152] = MobBuff.WeaponDefenseUp,
	[153] = MobBuff.MagicDefenseUp,
	[154] = MobBuff.Acc,
	[155] = MobBuff.Avoid,
	[156] = MobBuff.Speed,
}

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

local function guardian_reactor_id(carnival_data, team_id)
	if team_id == CARNIVAL_TEAM.RED then
		return carnival_data.reactor_red
	end
	return carnival_data.reactor_blue
end

local function guardian_name(team_id, num)
	return tostring(team_id) .. tostring(num)
end

local function is_guardian_alive(reactor)
	return reactor ~= nil and reactor:state() < GUARDIAN_DESTROYED_STATE
end

local function buff_guardian(mob, guardian)
	local flag = GUARDIAN_BUFFS[guardian.mob_skill_id]
	local skill = guardian.skill
	if flag == nil or skill == nil then
		return
	end
	local effect = skill:effect()
	mob:buff(flag, effect.x, effect.time, skill)
end

local function find_free_mob_gen_pos(field, carnival_data, team_id)
	local spawns = field:summoned_mob_spawns()
	for _, gen in ipairs(carnival_data.mob_gen_pos) do
		if gen.team == team_id or gen.team == CARNIVAL_TEAM.NONE then
			local used = false
			for _, spawn in ipairs(spawns) do
				if spawn.x == gen.x and spawn.y == gen.y and (gen.team == CARNIVAL_TEAM.NONE or spawn.team == gen.team) then
					used = true
					break
				end
			end
			if not used then
				return gen
			end
		end
	end
	return nil
end

local function find_free_guardian_gen_pos(field, carnival_data, team_id)
	for _, gen in ipairs(carnival_data.guardian_gen_pos) do
		if gen.team == team_id or gen.team == CARNIVAL_TEAM.NONE then
			local used = false
			for _, reactor in pairs(field:reactors()) do
				local x, y = reactor:position()
				if x == gen.x and y == gen.y and is_guardian_alive(reactor) then
					used = true
					break
				end
			end
			if not used then
				return gen
			end
		end
	end
	return nil
end

local function summon_mob(player, team, field, carnival_data, num)
	local entry = carnival_data.mobs[num + 1]
	local personal = team:personal_cp(player)
	if entry == nil or personal == nil or personal:available_cp() < entry.spend_cp then
		player:message("CP가 부족합니다.", Msg.PinkText)
		return false
	end
	local gen = find_free_mob_gen_pos(field, carnival_data, team:id())
	if gen == nil or not field:summon_mob(entry.id, gen.x, gen.y, team:id()) then
		player:message("더 이상 소환수를 불러낼 수 없습니다.", Msg.PinkText)
		return false
	end
	return team:use_cp(player, entry.spend_cp)
end

local function use_skill(match, player, team, field, carnival_data, num)
	local skill_id = carnival_data.skills[num + 1]
	if skill_id == nil then
		player:message("오류가 발생했습니다.", Msg.PinkText)
		return false
	end
	local skill = carnival.skill(skill_id)
	local personal = team:personal_cp(player)
	if skill == nil or personal == nil or personal:available_cp() < skill.spend_cp then
		player:message("CP가 부족합니다.", Msg.PinkText)
		return false
	end
	local enemy = match:enemy_team(team:id())
	if enemy == nil or not enemy:debuff(field, skill_id) then
		player:message("오류가 발생했습니다.", Msg.PinkText)
		return false
	end
	return team:use_cp(player, skill.spend_cp)
end

local function summon_guardian(player, team, field, carnival_data, num)
	local guardian = carnival.guardian(num)
	local personal = team:personal_cp(player)
	if guardian == nil or personal == nil or personal:available_cp() < guardian.spend_cp then
		player:message("CP가 부족합니다.", Msg.PinkText)
		return false
	end
	local name = guardian_name(team:id(), num)
	local gen = find_free_guardian_gen_pos(field, carnival_data, team:id())
	if is_guardian_alive(field:find_reactor_name(name))
		or gen == nil
		or field:spawn_reactor(guardian_reactor_id(carnival_data, team:id()), gen.x, gen.y, name, GUARDIAN_SPAWN_STATE) == nil then
		player:message("지금은 더 이상 수호물을 불러낼 수 없습니다.", Msg.PinkText)
		return false
	end
	for _, mob in pairs(field:mobs()) do
		if mob:carnival_team() == team:id() then
			buff_guardian(mob, guardian)
		end
	end
	return team:use_cp(player, guardian.spend_cp)
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
				map:respawn({ include_one_time = true })
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
		local carnival_data = field:carnival()
		if carnival_data == nil then
			return
		end

		local summoned = false
		if tab == CARNIVAL_TAB.MOB then
			summoned = summon_mob(player, team, field, carnival_data, num)
		elseif tab == CARNIVAL_TAB.SKILL then
			summoned = use_skill(match, player, team, field, carnival_data, num)
		elseif tab == CARNIVAL_TAB.GUARDIAN then
			summoned = summon_guardian(player, team, field, carnival_data, num)
		end
		if summoned then
			player:carnival_summon(tab, num)
		end
	end,

	on_mob_spawn = function(sm, mob)
		local team_id = mob:carnival_team()
		local field = mob:map()
		if team_id == CARNIVAL_TEAM.NONE or field == nil then
			return
		end
		local carnival_data = field:carnival()
		if carnival_data == nil then
			return
		end
		local reactor_id = guardian_reactor_id(carnival_data, team_id)
		local prefix = tostring(team_id)
		for _, reactor in pairs(field:reactors()) do
			local name = reactor:name()
			if reactor:id() == reactor_id and is_guardian_alive(reactor) and name:sub(1, #prefix) == prefix then
				local guardian = carnival.guardian(tonumber(name:sub(#prefix + 1)))
				if guardian ~= nil then
					buff_guardian(mob, guardian)
				end
			end
		end
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
