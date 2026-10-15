local pq = require("script/integration/lib/party_quest")

local DIALOGS = { resp.dialog, resp.dialog_yes_no, resp.dialog_list, resp.warp }
local CHIEFS = 211000001
local SHRINE = 211040401
local HOLY_STONE = 2030006
local BLACK_CHARM = 4031059
local DARK_CRYSTAL = 4005004
local STRENGTH_NECKLACE = 4031057
local WISDOM_NECKLACE = 4031058
local TRUTH = "진실이다"

local ANSWERS = {
	["레벨1에서 2로 갈때 필요한 경험치량"] = { "15" },
	["1차전직에 필요한 조건이 아닌것"] = { "도적 럭 20" },
	["상태이상과 약화 효과가 올바르게 짝지어 지지 않은것"] = { "허약 - 이동속도가 느려짐" },
	["몬스터를 공격할때 틀린것"] = { "독 - 보스 몬스터에게 강한 데미지" },
	["적중률에 가장 많이 의존하는 직업"] = { "궁수" },
	["몬스터가 드롭하는 아이템으로 틀린것"] = { "레이스 - 식탁보", "엑스텀프 - 나뭇가지" },
	["포션의 효과로 알맞은것"] = { "피자 - HP 400 회복" },
	["HP와 MP를 50% 회복하는 아이템"] = { "엘릭서" },
	["포션의효과로 틀린것"] = { "새벽의 이슬 - MP 3000 회복" },
	["가장 레벨이 높은 몬스터"] = { "엑스텀프" },
	["메이플 아일랜드에서 볼 수 없는 몬스터"] = { "아이스 센티넬", "파이어보어" },
	["오르비스로 이동하는 배에서 볼 수 있는 몬스터"] = { "크림슨 발록" },
	["빅토리아아일랜드에서 볼수없는 몬스터"] = { "헥터" },
	["엘나스에서 볼 수 없는 몬스터"] = { "리게이터" },
	["비행하는 몬스터"] = { "맬러디" },
	["오시리아 대륙에서 볼 수 없는 몬스터"] = { "리게이터" },
	["스텀프 50마리를 잡는 퀘스트"] = { "스텀프가 무서워요" },
	["반복수행이 가능한 퀘스트"] = { "아르웬의 유리구두" },
	["2차전직이 아닌것"] = { "메이지" },
	["알을 모아오는 퀘스트를 주는 NPC"] = { "네미" },
	["인기도를 주는 퀘스트를 가진 NPC"] = { "슈미" },
	["'하인즈'가 있는 마을"] = { "엘리니아" },
	["'알케스터'가 있는 마을"] = { "엘나스" },
	["'엘나스'에서 신발을 만드는 NPC"] = { "고든" },
	["헤네시스에서 화살을 만드는 NPC"] = { "비셔스" },
	["'장난감공장'은 어느대륙"] = { "루디브리엄" },
	["'네미'는 어떤 마을"] = { "루디브리엄" },
}

local function solve(p)
	for question, answers in pairs(ANSWERS) do
		if p.text:find(question, 1, true) ~= nil then
			for i, selection in ipairs(p.selections) do
				for _, answer in ipairs(answers) do
					if selection == answer then
						return i - 1
					end
				end
			end
		end
	end
	return nil
end

local function talk(ctx, bot, template, pick)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local p, name = bot:request(DIALOGS, req.npc_click { oid = oid })
	local last = nil
	for _ = 1, 30 do
		if p == false and last ~= nil then
			return last
		end
		if p == false then
			return ctx:fail(string.format("NPC %d 대화 중 응답 없음 (맵 %d)", template, bot:map()))
		end
		last = nil
		if name == resp.warp then
			return true
		end
		if name == resp.dialog and p.next == false and p.text:find(TRUTH, 1, true) == nil then
			bot:dialog(false)
			return p.text
		end
		local reply = req.dialog { dialog_type = 0, next = true }
		if name == resp.dialog_yes_no then
			reply = req.dialog { dialog_type = 1, next = true }
		elseif name == resp.dialog_list then
			local selected = pick(p)
			if selected == nil then
				bot:dialog(false)
				return ctx:fail(string.format("NPC %d 선택지를 고르지 못함: %s", template, p.text))
			end
			reply = req.dialog { dialog_type = 4, next = true, selected = selected }
		end
		if name == resp.dialog then
			last = p.text
			p, name = bot:request(DIALOGS, reply, nil, 1500)
		else
			p, name = bot:request(DIALOGS, reply)
		end
	end
	return ctx:fail(string.format("NPC %d 대화가 끝나지 않음", template))
end

local function first(p)
	return 0
end

local function visit(ctx, bot, map_id, template, pick)
	if bot:map_move(map_id) == false then
		return ctx:fail("맵 이동 실패: " .. map_id)
	end
	return talk(ctx, bot, template, pick or first)
end

local function advance(ctx, path)
	local bot = ctx:bot(0)
	local items = string.format("%d:1,%d:1", BLACK_CHARM, DARK_CRYSTAL)
	if pq.command(bot, "/봇초기화 70 " .. path.second .. " 0 " .. items .. " -", "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패")
	end
	if visit(ctx, bot, CHIEFS, path.chief) == false then
		return false
	end
	if visit(ctx, bot, path.mentor_map, path.mentor) == false then
		return false
	end
	if visit(ctx, bot, path.mentor_map, path.mentor) == false then
		return false
	end
	if (bot:items()[STRENGTH_NECKLACE] or 0) ~= 1 then
		return ctx:fail("강인함의 목걸이를 받지 못함")
	end
	if visit(ctx, bot, CHIEFS, path.chief) == false then
		return false
	end
	if (bot:items()[STRENGTH_NECKLACE] or 0) ~= 0 then
		return ctx:fail("강인함의 목걸이가 회수되지 않음")
	end
	if visit(ctx, bot, SHRINE, HOLY_STONE, solve) == false then
		return false
	end
	if (bot:items()[WISDOM_NECKLACE] or 0) ~= 1 then
		return ctx:fail("지혜의 목걸이를 받지 못함")
	end
	local said = visit(ctx, bot, CHIEFS, path.chief)
	if said == false then
		return false
	end
	if bot:class() ~= path.third then
		return ctx:fail(string.format("3차 전직 실패: %d (기대 %d) %s", bot:class(), path.third, tostring(said)))
	end
	return true
end

test_suite {
	name = "Job: 3차 전직",
	bot_count = 1,

	on_initialize = function(ctx)
		if pq.command(ctx:bot(0), "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return true
	end,

	scenarios = {
		function(ctx)
			return advance(ctx, { second = 110, third = 111, chief = 2020008, mentor = 1022000, mentor_map = 102000003 })
		end,
		function(ctx)
			return advance(ctx, { second = 210, third = 211, chief = 2020009, mentor = 1032001, mentor_map = 101000003 })
		end,
		function(ctx)
			return advance(ctx, { second = 310, third = 311, chief = 2020010, mentor = 1012100, mentor_map = 100000201 })
		end,
		function(ctx)
			return advance(ctx, { second = 410, third = 411, chief = 2020011, mentor = 1052001, mentor_map = 103000003 })
		end,
		function(ctx)
			return advance(ctx, { second = 510, third = 511, chief = 2020013, mentor = 1090000, mentor_map = 120000101 })
		end,
	},
}
