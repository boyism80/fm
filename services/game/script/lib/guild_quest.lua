local M = {}

M.GROUP = "guild_quest"
M.ID = "guild_quest"
M.EXIT_MAP = 990001100
M.WAIT_MS = 180000
M.DURATION_MS = 5400000
M.MIN_PLAYERS = 6

local MAZE_QUEST = 7600
local GATE_MAP = 990000300

function M.is_leader(me, sm)
	local leader = sm:leader()
	return leader ~= nil and leader:id() == me:id()
end

function M.gain_gp_once(me, sm, key, amount)
	if sm:get_property(key) ~= "" then
		return
	end
	sm:set_property(key, "true")
	local guild = me:guild()
	if guild ~= nil then
		guild:gain_gp(amount)
	end
end

function M.pass_gate(me, name, clear_key, map_id, spawn, closed)
	local gate = me:map():find_reactor_name(name)
	local opened = gate ~= nil and gate:state() == 1
	if opened == false and clear_key ~= nil then
		local sm = me:state_machine()
		opened = sm ~= nil and sm:get_property(clear_key) == "true"
	end
	if opened == false then
		me:message(closed or "지금은 포탈이 닫혀있습니다.", Msg.PinkText)
		return
	end
	me:play_portal_sound()
	me:map(map_id, spawn or 0)
end

function M.maze(me)
	return me:quest(MAZE_QUEST):record()
end

function M.set_maze(me, value)
	local q = me:quest(MAZE_QUEST)
	if q:started() then
		q:record(value)
		return
	end
	q:start(value)
end

function M.combo_length(sm)
	return (tonumber(sm:get_property("stage1phase")) or 1) + 3
end

function M.record_statue(reactor)
	local map = reactor:map()
	if map == nil or map:wz():id() ~= GATE_MAP then
		return
	end
	local sm = state_machine(M.GROUP):get(M.ID)
	if sm == nil then
		return
	end

	local status = sm:get_property("stage1status")
	if status == "display" then
		local combo = sm:get_property("stage1combo") .. reactor:oid() .. ","
		sm:set_property("stage1combo", combo)
		if select(2, combo:gsub(",", "")) == M.combo_length(sm) then
			sm:set_property("stage1status", "active")
			sm:set_property("stage1guess", "")
		end
		return
	end
	if status == "active" then
		local guess = sm:get_property("stage1guess")
		if select(2, guess:gsub(",", "")) < M.combo_length(sm) then
			sm:set_property("stage1guess", guess .. reactor:oid() .. ",")
		end
	end
end

return M
