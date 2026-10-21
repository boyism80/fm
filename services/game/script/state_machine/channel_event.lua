-- State machine (old/src/server/events/MapleEvent.java): 채널 이벤트 모집·시작·종료

local GAMES = {
	ox_quiz = require("script/lib/event/ox_quiz"),
}

local WINNER_MAP = 109050000
local LOSER_MAP = 109050001
local RECRUIT_MS = 180000
local COUNTDOWN_MS = 30000
local WINNER_STAY_MS = 600000
local MAX_PLAYERS = 50
local RULES = { "chat", "skill", "chalkboard" }

-- kind and name come from the start call; the machine's properties vanish with it, so a failed create leaves no state.
local function game(sm)
	return GAMES[sm:get_property("kind")]
end

local function set_phase(sm, phase)
	sm:set_property("phase", phase)
end

local function phase(sm)
	return sm:get_property("phase")
end

local function set_rules(sm, on)
	for _, rule in ipairs(RULES) do
		sm:set_property("forbid_" .. rule, tostring(on))
	end
end

local function event_players(sm)
	local players = {}
	for _, map_id in ipairs(game(sm).MAPS) do
		for _, player in pairs(sm:map(map_id):characters()) do
			players[#players + 1] = player
		end
	end
	return players
end

local function close(sm)
	if phase(sm) ~= "recruiting" then
		return
	end
	set_phase(sm, "closing")
	sm:cancel("close")
	world_notice("이벤트 참가 모집이 종료되었습니다.")
	for _, player in ipairs(event_players(sm)) do
		player:message("이벤트가 30초 후 시작됩니다!", Msg.Notice)
		player:clock(COUNTDOWN_MS / 1000)
	end
	sm:after("begin", COUNTDOWN_MS, "on_begin")
end

local hooks = {
	on_init = function(group)
		group:max_machines(1)
	end,

	on_create = function(sm)
		local maps = { WINNER_MAP }
		for _, map_id in ipairs(game(sm).MAPS) do
			maps[#maps + 1] = map_id
		end
		return maps
	end,

	on_prepare = function(sm)
		set_phase(sm, "recruiting")
		set_rules(sm, true)
		world_notice(string.format("GM 이벤트 [%s]가 잠시 후 %d채널에서 시작됩니다! 해당 채널의 리스항구, 사우스페리, 오르비스, 루디브리엄에서 이벤트 NPC를 통해 참가할 수 있습니다!", sm:get_property("name"), channel_id() + 1), Msg.LightBlueText)
		sm:after("close", RECRUIT_MS, "on_close")
	end,

	on_player_enter = function(sm, player)
		if phase(sm) ~= "recruiting" then
			player:clear_saved_location("EVENT")
			sm:unregister(player)
			return
		end
		local function on_arrive(player)
			player:event_instructions()
		end
		player:map(sm:map(game(sm).MAPS[1]), "join00", { callback = on_arrive })
		if #sm:players() >= MAX_PLAYERS then
			close(sm)
		end
	end,

	on_close = function(sm)
		close(sm)
	end,

	on_begin = function(sm)
		set_phase(sm, "running")
		game(sm).start(sm)
	end,

	-- The game marks its winners with the "winner:<id>" property before it calls this.
	on_game_over = function(sm)
		set_phase(sm, "finished")
		set_rules(sm, false)
		local winners = 0
		for _, player in ipairs(event_players(sm)) do
			player:stop_clock()
			if sm:get_property("winner:" .. player:id()) == "true" then
				winners = winners + 1
				player:map(sm:map(WINNER_MAP))
			else
				player:map(LOSER_MAP)
			end
		end
		if winners == 0 then
			sm:finish(LOSER_MAP)
			return
		end
		sm:after("winner_stay", WINNER_STAY_MS, "on_winner_timeout")
	end,

	on_winner_timeout = function(sm)
		sm:finish(LOSER_MAP)
	end,

	on_cancel = function(sm)
		world_notice("이벤트가 취소되었습니다.")
		sm:finish(LOSER_MAP)
	end,

	on_player_leave = function(sm, player, reason)
		if phase(sm) == "finished" and #sm:players() == 0 then
			sm:finish(0)
		end
	end,
}

for _, g in pairs(GAMES) do
	for name, fn in pairs(g.hooks) do
		hooks[name] = fn
	end
end

return hooks
