-- State machine (old/scripts/event/GuildQuest.js): 길드 대항전

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local WAITING_MAP = 990000000
local GATE_MAP = 990000300
local BONUS_MAP = 990001000
local BOARD_NPC = 9040011
local BONUS_MS = 40000

local STAGE_MAPS = {
	990000000,
	990000100,
	990000200,
	990000300,
	990000301,
	990000400,
	990000401,
	990000410,
	990000420,
	990000430,
	990000431,
	990000440,
	990000500,
	990000501,
	990000502,
	990000600,
	990000610,
	990000611,
	990000620,
	990000630,
	990000631,
	990000640,
	990000641,
	990000700,
	990000800,
	990000900,
	990001000,
}

local SHUFFLE_MAPS = { 990000501, 990000502 }
local MAZE_END_MAPS = { 990000611, 990000620, 990000631, 990000641 }

local LEADER_LEFT = {
	disconnect = "길드 대항전의 리더가 접속이 끊겨 자동으로 모두 퇴장됩니다.",
	map = "길드 대항전의 리더가 퇴장하여 자동으로 모두 퇴장됩니다.",
}

local GATE_NOTICES = {
	{ at = 60000, left = "2분" },
	{ at = 120000, left = "1분" },
	{ at = 150000, left = "30초" },
}

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_min_players(1)
		group:declare_exit_map(gq.EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		sm:set_property("state", "waiting")
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn({ include_one_time = true })
			end
		end
		for _, map_id in ipairs(SHUFFLE_MAPS) do
			pq.shuffle_reactors(group:map(map_id))
		end
		for _, map_id in ipairs(MAZE_END_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:set_reactor_respawn(false)
			end
		end
		local waiting = group:map(WAITING_MAP)
		if waiting ~= nil then
			local portal = waiting:portal(5)
			if portal ~= nil then
				portal:script("guildwaitingenter")
			end
		end
		return STAGE_MAPS
	end,

	on_start = function(sm)
		sm:start_timer(gq.WAIT_MS)
		sm:after(GATE_NOTICES[1].at, "on_gate_notice")
	end,

	on_player_enter = function(sm, player)
		local function on_arrive(player)
			player:open_npc(BOARD_NPC)
		end
		player:map(WAITING_MAP, 0, { callback = on_arrive })
	end,

	on_player_leave = function(sm, player, reason)
		local leader_left = LEADER_LEFT[reason]
		if leader_left ~= nil and gq.is_leader(player, sm) then
			sm:message(leader_left)
			sm:finish(gq.EXIT_MAP)
			return
		end
		if reason ~= "disconnect" and reason ~= "exit" then
			return
		end
		if sm:get_property("state") == "waiting" then
			return
		end
		if #sm:players() < gq.MIN_PLAYERS then
			sm:message("길드 대항전을 진행할 인원이 부족하여 자동으로 모두 퇴장됩니다.")
			sm:finish(gq.EXIT_MAP)
		end
	end,

	on_gate_notice = function(sm)
		if sm:get_property("state") ~= "waiting" then
			return
		end
		local index = (tonumber(sm:get_property("gate_notice")) or 0) + 1
		local notice = GATE_NOTICES[index]
		if notice == nil then
			return
		end
		sm:set_property("gate_notice", tostring(index))
		sm:message("샤레니안으로 가는 문이 " .. notice.left .. " 후에 열립니다. 유적 발굴 현장에 최소 " .. gq.MIN_PLAYERS .. "명 이상이 입장해 있어야 합니다.")
		local next_notice = GATE_NOTICES[index + 1]
		if next_notice ~= nil then
			sm:after(next_notice.at - notice.at, "on_gate_notice")
		end
	end,

	on_statue_display = function(sm)
		local oid, rest = sm:get_property("stage1pending"):match("^(%d+),(.*)$")
		if oid == nil then
			return
		end
		sm:set_property("stage1pending", rest)
		local gate = sm:group():map(GATE_MAP)
		if gate ~= nil then
			local statue = gate:reactors()[tonumber(oid)]
			if statue ~= nil then
				statue:hit()
			end
		end
		if rest ~= "" then
			sm:after(3500, "on_statue_display")
		end
	end,

	on_scheduled_timeout = function(sm)
		local state = sm:get_property("state")
		if state == "waiting" then
			if #sm:players() < gq.MIN_PLAYERS then
				sm:message("길드 대항전을 시작하려면 " .. gq.MIN_PLAYERS .. "명의 참가 인원이 필요합니다.")
				sm:finish(gq.EXIT_MAP)
				return
			end
			sm:set_property("state", "started")
			sm:message("샤레니안의 문이 열렸습니다.", Msg.Notice)
			sm:restart_timer(gq.DURATION_MS)
		elseif state == "started" then
			sm:message("제한시간이 다 되었습니다.")
			sm:finish(gq.EXIT_MAP)
		else
			sm:finish(gq.EXIT_MAP)
		end
	end,

	on_clear = function(sm)
		local bonus = sm:group():map(BONUS_MAP)
		if bonus ~= nil then
			bonus:reload_reactors()
		end
		sm:set_property("state", "bonus")
		pq.party_warp(sm, BONUS_MAP)
		sm:restart_timer(BONUS_MS)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}
