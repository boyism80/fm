local M = {}

M.GROUP = "mu_lung_dojo"
M.TUTORIAL_GROUP = "mu_lung_dojo_tutorial"
M.TUTORIAL = 925020010
M.SO_GONG = 9300269
M.LOBBY = 925020001
M.EXIT = 925020002
M.ROOF = 925020003
M.ROOF_PORTAL = 1
M.LAST_FLOOR = 38
M.DOOR_PORTAL = 6
M.DOOR_REACTOR = 2508000
M.POINT_QUEST = 150100
M.BELT_QUEST = 150101
M.REST_QUEST = 150000
M.BEST_FLOOR_QUEST = 150102
M.BEST_TIME_QUEST = 150103
M.BOSS_SPAWNS = {
	{ 140, 0 },
	{ -193, 0 },
	{ 355, 0 },
}
M.TAUNT_EFFECT = 5120024
M.TAUNT_MS = 30000
M.TAUNTS = {
	"무릉도장에 도전한 것을 후회하게 해주겠다! 어서 들어와봐!",
	"기다리고 있었다! 용기가 남았다면 들어와 보시지!",
	"배짱 하나는 두둑하군! 현명함과 무모함을 혼동하지말라고!",
	"무릉도장에 도전하다니 용기가 가상하군!",
	"패배의 길을 걷고싶다면 들어오라고!",
}
M.BELTS = {
	{ item = 1132000, level = 25, points = 20 },
	{ item = 1132001, level = 35, points = 300 },
	{ item = 1132002, level = 45, points = 800 },
	{ item = 1132003, level = 60, points = 900 },
	{ item = 1132004, level = 75, points = 1200 },
}

function M.map_id(floor)
	return 925020000 + floor * 100
end

function M.floor_of(map_id)
	local floor = math.floor((map_id - 925020000) / 100)
	if floor < 1 or floor > M.LAST_FLOOR then
		return nil
	end
	return floor
end

function M.is_rest(floor)
	return floor % 6 == 0
end

function M.tier(floor)
	return math.floor((floor - 1) / 6) + 1
end

function M.stage(floor)
	return floor - M.tier(floor) + 1
end

function M.boss(floor)
	return 9300183 + M.stage(floor)
end

function M.time_limit_ms(floor)
	local tier = M.tier(floor)
	if tier == 7 then
		return 15 * 60000
	end
	return (4 + tier) * 60000
end

function M.record(me, quest_id)
	local q = me:quest(quest_id)
	if not q:started() then
		return 0
	end
	return tonumber(q:record()) or 0
end

function M.set_record(me, quest_id, value)
	local q = me:quest(quest_id)
	if not q:started() then
		q:start(tostring(value))
		return
	end
	q:record(tostring(value))
end

function M.taunt(map)
	if map:property("taunt") then
		return
	end
	map:property("taunt", true)
	map:weather(M.TAUNT_EFFECT, M.TAUNTS[math.random(1, #M.TAUNTS)])
	sleep(M.TAUNT_MS)
	map:weather(0)
	map:property("taunt", false)
end

function M.current_floor(me)
	local map = me:map()
	if map == nil then
		return nil
	end
	return M.floor_of(map:template_id())
end

function M.cleared(sm, floor)
	return sm:get_property("cleared:" .. floor) == "1"
end

function M.players_on(sm, floor)
	local out = {}
	for _, player in ipairs(sm:players()) do
		if M.current_floor(player) == floor then
			table.insert(out, player)
		end
	end
	return out
end

function M.grant_points(players, floor)
	local points = M.tier(floor) * 4
	if #players > 1 then
		points = M.tier(floor) * 3
	end
	for _, player in ipairs(players) do
		local total = M.record(player, M.POINT_QUEST) + points
		M.set_record(player, M.POINT_QUEST, total)
		player:message(string.format("수련점수를 %d점 받았습니다. 총 수련점수가 %d점이 되었습니다.", points, total), Msg.PinkText)
	end
end

function M.record_floor(players, floor)
	for _, player in ipairs(players) do
		if floor > M.record(player, M.BEST_FLOOR_QUEST) then
			M.set_record(player, M.BEST_FLOOR_QUEST, floor)
		end
	end
end

function M.record_time(players, seconds)
	for _, player in ipairs(players) do
		local best = M.record(player, M.BEST_TIME_QUEST)
		if best == 0 or seconds < best then
			M.set_record(player, M.BEST_TIME_QUEST, seconds)
			player:message(string.format("무릉도장 최단 완주 기록을 세웠습니다: %s", M.format_time(seconds)), Msg.PinkText)
		end
	end
end

function M.format_time(seconds)
	return string.format("%d분 %02d초", math.floor(seconds / 60), seconds % 60)
end

function M.advance(sm, floor)
	local key = "advanced:" .. floor
	if sm:get_property(key) == "1" then
		return
	end
	sm:set_property(key, "1")
	local players = M.players_on(sm, floor)
	if not M.is_rest(floor) then
		M.grant_points(players, floor)
		M.record_floor(players, floor)
	end
	if floor >= M.LAST_FLOOR then
		if sm:get_property("start") == "1" then
			M.record_time(players, os.time() - tonumber(sm:get_property("started_at")))
		end
		sm:finish(M.ROOF, M.ROOF_PORTAL)
		return
	end
	local next_map = sm:map(M.map_id(floor + 1))
	for _, player in ipairs(players) do
		player:map(next_map)
	end
end

return M
