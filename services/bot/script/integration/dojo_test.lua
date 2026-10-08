local pq = require("script/integration/lib/party_quest")

local LOBBY = 925020001
local EXIT = 925020002
local FLOOR_1 = 925020100
local FLOOR_2 = 925020200
local SO_GONG = 2091005
local BOSS = 9300184
local SNAIL = 100100
local WHITE_BELT = 1132000
local POINT_QUEST = 150100

local DIALOG_YES_NO = 1
local DIALOG_LIST = 4

test_suite {
	name = "Mu Lung Dojo",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 30 0 0 - - " .. POINT_QUEST .. "=16", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(LOBBY)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, SO_GONG)
			if oid == false then
				return false
			end
			local dlg = bot:npc_click(oid)
			if dlg == false or dlg.selections == nil then
				return ctx:fail("소공 선택지가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
					return p.character.map == FLOOR_1
				end, 10000) == false then
				return ctx:fail("1층으로 이동하지 않음")
			end
			if bot:request(resp.tremble, nil, nil, 5000) == false then
				return ctx:fail("층 시작 연출이 오지 않음")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local mob = bot:mobs(BOSS)[1]
			if mob == nil then
				local spawn = bot:request(resp.spawn_mob, nil, function(p)
					return p.mob.mob_id == BOSS
				end, 10000)
				if spawn == false then
					return ctx:fail("1층 보스가 나타나지 않음")
				end
				mob = { oid = spawn.mob.oid }
			end
			if pq.command(bot, "/몬스터생성 " .. SNAIL, "몬스터 생성:") == false then
				return ctx:fail("일반 몬스터 생성 실패")
			end
			local snail = bot:mobs(SNAIL)[1]
			if snail == nil then
				return ctx:fail("일반 몬스터가 보이지 않음")
			end
			bot:attack(snail.oid, 1)
			local energy = bot:request(resp.session_value, nil, function(p)
				return p.key == "energy"
			end, 5000)
			if energy == false or tonumber(energy.value) <= 0 then
				return ctx:fail("도장 에너지가 차지 않음")
			end
			if bot:kill(mob.oid) == false then
				return ctx:fail("1층 보스 처치 실패")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
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
				return p.message:find("수련점수를 4점", 1, true) ~= nil
			end, 5000)
			if notice == false then
				return ctx:fail("수련점수를 받지 못함")
			end
			return pq.wait_map(ctx, bot, FLOOR_2, 5000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, SO_GONG)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("2층 소공 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return p.character.map == EXIT
				end, 10000) == false then
				return ctx:fail("포기 후 퇴장하지 않음")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(LOBBY) == false then
				return ctx:fail("로비로 돌아가지 못함")
			end
			local oid = pq.npc(ctx, bot, SO_GONG)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("소공 대화가 오지 않음")
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
	},
}
