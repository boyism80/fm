local pq = require("script/integration/lib/party_quest")

local LOBBY = 925020001
local EXIT = 925020002
local FLOOR_1 = 925020100
local FLOOR_2 = 925020200
local FLOOR_6 = 925020600
local FLOOR_7 = 925020700
local SO_GONG = 2091005
local SO_GONG_MOB = 9300269
local TUTORIAL = 925020010
local BOSS = 9300184
local SNAIL = 100100
local WHITE_BELT = 1132000
local POINT_RECORD = "dojo.points"
local REST_RECORD = "dojo.rest"
local BAMBOO_RAIN = 1009
local INVINCIBILITY = 1010
local KILL_DAMAGE = 2147483647

local DIALOG_YES_NO = 1
local DIALOG_LIST = 4

local function talk(ctx, bot)
	local oid = pq.npc(ctx, bot, SO_GONG)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail(bot:name() .. " 소공 대화가 오지 않음")
	end
	return dlg
end

local function choose_warp(ctx, bot, selected, map_id, what)
	local dlg = talk(ctx, bot)
	if dlg == false then
		return false
	end
	if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = selected }, function(p)
			return p.character.map == map_id
		end, 10000) == false then
		return ctx:fail(bot:name() .. " " .. what)
	end
	return true
end

local function give_up(ctx, bot)
	if talk(ctx, bot) == false then
		return false
	end
	if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
			return p.character.map == EXIT
		end, 10000) == false then
		return ctx:fail(bot:name() .. " 포기 후 퇴장하지 않음")
	end
	return true
end

local function find_boss(ctx, bot, boss)
	local mob = bot:mobs(boss)[1]
	if mob ~= nil then
		return mob
	end
	local spawn = bot:request(resp.spawn_mob, nil, function(p)
		return p.mob.mob_id == boss
	end, 10000)
	if spawn == false then
		return ctx:fail("보스가 나타나지 않음: " .. boss)
	end
	return { oid = spawn.mob.oid }
end

local function energy(bot, value)
	return bot:request(resp.session_value, nil, function(p)
		return p.key == "energy" and p.value == value
	end, 5000)
end

local function charge_full(ctx, bot)
	if pq.command(bot, "/몬스터생성 " .. SNAIL .. " 6", "몬스터 생성:") == false then
		return ctx:fail("일반 몬스터 생성 실패")
	end
	local snails = bot:mobs(SNAIL)
	if #snails < 6 then
		return ctx:fail("일반 몬스터가 모두 보이지 않음")
	end
	for i = 1, 5 do
		if bot:kill(snails[i].oid) == false then
			return ctx:fail("일반 몬스터 처치 실패")
		end
	end
	bot:attack(snails[6].oid, KILL_DAMAGE)
	if energy(bot, "300") == false then
		return ctx:fail("도장 에너지가 가득 차지 않음")
	end
	return true
end

local function climb(ctx, bot, next_map, points)
	if pq.move(bot, 429, 7) == false then
		return ctx:fail("문 포탈로 이동 실패")
	end
	if bot:request(resp.field_relocate, req.warp { target = 0xFFFFFFFF, portal_name = "out00" }, nil, 5000) == false then
		return ctx:fail("문 포탈로 올라가지 않음")
	end
	if pq.move(bot, 6, -367) == false then
		return ctx:fail("다음 층 포탈로 이동 실패")
	end
	local notice = bot:request(resp.notice, req.warp { target = 0xFFFFFFFF, portal_name = "out001" }, function(p)
		return p.message:find("수련점수를 " .. points .. "점", 1, true) ~= nil
	end, 5000)
	if notice == false then
		return ctx:fail("수련점수를 받지 못함")
	end
	return pq.wait_map(ctx, bot, next_map, 5000)
end

test_suite {
	name = "Mu Lung Dojo",
	bot_count = 2,

	on_finished = function(ctx)
		pq.command(ctx:bot(0), "/타이머 0", "타이머 제한: 0")
	end,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/봇초기화 30 0 0 - - - " .. POINT_RECORD .. "=16", "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(LOBBY) == false then
				return false
			end
		end
		return true
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if choose_warp(ctx, bot, 0, FLOOR_1, "1층으로 이동하지 않음") == false then
				return false
			end
			if bot:request(resp.tremble, nil, nil, 5000) == false then
				return ctx:fail("층 시작 연출이 오지 않음")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if charge_full(ctx, bot) == false then
				return false
			end
			local cast = bot:request(resp.session_value, req.active_skill { skill_id = INVINCIBILITY, skill_level = 1 }, function(p)
				return p.key == "energy" and p.value == "0"
			end, 5000)
			if cast == false then
				return ctx:fail("금강불괴를 쓰지 못함")
			end			if charge_full(ctx, bot) == false then
				return false
			end
			if pq.command(bot, "/몬스터생성 " .. SNAIL, "몬스터 생성:") == false then
				return ctx:fail("일반 몬스터 생성 실패")
			end
			local snail = bot:mobs(SNAIL)[1]
			if snail == nil then
				return ctx:fail("일반 몬스터가 보이지 않음")
			end
			bot:attack(snail.oid, KILL_DAMAGE, 1, BAMBOO_RAIN)
			if energy(bot, "0") == false then
				return ctx:fail("죽간천격을 쓰지 못함")
			end			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local boss = find_boss(ctx, bot, BOSS)
			if boss == false then
				return false
			end
			if bot:kill(boss.oid) == false then
				return ctx:fail("1층 보스 처치 실패")
			end
			if climb(ctx, bot, FLOOR_2, 4) == false then
				return false
			end
			if give_up(ctx, bot) == false then
				return false
			end
			if bot:map_move(LOBBY) == false then
				return ctx:fail("로비로 돌아가지 못함")
			end
			if talk(ctx, bot) == false then
				return false
			end
			local record = bot:dialog(true, 4)
			if record == nil or record.text:find("1층#k까지", 1, true) == nil then
				return ctx:fail("최고 층 기록이 남지 않음")
			end
			bot:dialog(false)
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if pq.command(bot, "/봇초기화 30 0 0 - - - " .. POINT_RECORD .. "=20," .. REST_RECORD .. "=1", "봇초기화 완료") == false then
				return ctx:fail("휴식층 저장 상태 설정 실패")
			end
			if bot:map_move(LOBBY) == false then
				return ctx:fail("로비로 돌아가지 못함")
			end
			if choose_warp(ctx, bot, 0, FLOOR_6, "저장된 6층으로 이동하지 않음") == false then
				return false
			end
			if talk(ctx, bot) == false then
				return false
			end
			local saved = bot:dialog(true, 2)
			if saved == nil or saved.text:find("진행상황이 저장되었어", 1, true) == nil then
				return ctx:fail("휴식층 진행상황을 저장하지 못함")
			end
			bot:dialog(false)
			if pq.command(bot, "/타이머 3", "타이머 제한: 3") == false then
				return ctx:fail("타이머 제한 설정 실패")
			end
			return choose_warp(ctx, bot, 0, FLOOR_7, "7층으로 이동하지 않음")
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if pq.wait_map(ctx, bot, EXIT, 15000) == false then
				return ctx:fail("시간 초과 후 퇴장하지 않음")
			end
			if pq.command(bot, "/타이머 0", "타이머 제한: 0") == false then
				return ctx:fail("타이머 제한 해제 실패")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(LOBBY) == false then
				return ctx:fail("로비로 돌아가지 못함")
			end
			if talk(ctx, bot) == false then
				return false
			end
			local belts = bot:dialog(true, 2)
			if belts == nil or belts.selections == nil then
				return ctx:fail("허리띠 목록이 오지 않음")
			end
			if bot:request(resp.inventory_operation, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, nil, 5000) == false then
				return ctx:fail("하얀 띠를 받지 못함")
			end
			if bot:items()[WHITE_BELT] == nil then
				return ctx:fail("인벤토리에 하얀 띠가 없음")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if talk(ctx, bot) == false then
				return false
			end
			if bot:request(resp.weather, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 5 }, function(p)
					return p.item_id ~= 0
				end, 10000) == false then
				return ctx:fail("소공의 방 도발 메시지가 오지 않음")
			end
			if bot:map() ~= TUTORIAL then
				return ctx:fail("소공의 방으로 이동하지 않음")
			end
			local so_gong = find_boss(ctx, bot, SO_GONG_MOB)
			if so_gong == false then
				return false
			end
			for _, mob in ipairs(bot:mobs(SO_GONG_MOB)) do
				if bot:kill(mob.oid) == false then
					return ctx:fail("소공 처치 실패")
				end
			end
			if pq.move(bot, 2, -366) == false then
				return ctx:fail("소공의 방 출구로 이동 실패")
			end
			if bot:request(resp.warp, req.warp { target = 0xFFFFFFFF, portal_name = "out000" }, function(p)
					return p.character.map == LOBBY
				end, 5000) == false then
				return ctx:fail("소공의 방에서 로비로 나가지 못함")
			end
			return true
		end,
		function(ctx)
			if pq.form_party(ctx) == false then
				return false
			end
			local leader = ctx:bot(0)
			if choose_warp(ctx, leader, 1, FLOOR_1, "파티 1층으로 이동하지 않음") == false then
				return false
			end
			return pq.wait_all(ctx, FLOOR_1, 10000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local member = ctx:bot(1)
			local boss = find_boss(ctx, leader, BOSS)
			if boss == false then
				return false
			end
			if leader:kill(boss.oid) == false then
				return ctx:fail("파티 1층 보스 처치 실패")
			end
			if pq.move(leader, 429, 7) == false then
				return ctx:fail("문 포탈로 이동 실패")
			end
			if leader:request(resp.field_relocate, req.warp { target = 0xFFFFFFFF, portal_name = "out00" }, nil, 5000) == false then
				return ctx:fail("문 포탈로 올라가지 않음")
			end
			if pq.move(leader, 6, -367) == false then
				return ctx:fail("다음 층 포탈로 이동 실패")
			end
			local notice = leader:request_on(member, resp.notice, req.warp { target = 0xFFFFFFFF, portal_name = "out001" }, function(p)
				return p.message:find("수련점수를 3점", 1, true) ~= nil
			end, 5000)
			if notice == false then
				return ctx:fail("파티원이 수련점수를 받지 못함")
			end
			return pq.wait_all(ctx, FLOOR_2, 10000)
		end,
		function(ctx)
			if give_up(ctx, ctx:bot(0)) == false then
				return false
			end
			return pq.wait_all(ctx, EXIT, 10000)
		end,
	},
}
