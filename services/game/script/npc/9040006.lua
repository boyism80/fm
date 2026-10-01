-- NPC name (String.wz/Npc.img.xml): 분수 조각상

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local FIRST_OFFERING = 4001027
local LAST_OFFERING = 4001030
local STATUES = 4
local ATTEMPTS = 7

local INTRO = "분수 앞에 옛 샤레니안의 충신들의 석상이 있다. 이 충신들은 살아 생전 왕이 수여한 귀중한 보물들을 소유하고 있었다. 이러한 보물들을 제물로 바쳐야 하며, 제대로 바쳐질 경우 비밀 문이 열리게 된다. 제물로 바쳐야 할 물건들은 다음과 같다.\r\r #v4001027# #t4001027#\r #v4001028# #t4001028#\r #v4001029# #t4001029#\r #v4001030# #t4001030#\r. 바로 이것들이오. 그러나 누가 어느 제물을 원하는지는 알 수 없다오. 또, 개중에는 아예 필요없는 제물도 있을 수 있소. 기회는 총 7번이라오."

local function make_combo()
	local combo = ""
	for _ = 1, STATUES do
		combo = combo .. (math.random(0, 3))
	end
	return combo
end

local function ground_guess(me)
	local map = me:map()
	local wz = map:wz()
	local items = {}
	for _, item in pairs(map:items()) do
		items[#items + 1] = item
	end
	if #items ~= STATUES then
		return nil
	end

	local placed = {}
	for _, item in ipairs(items) do
		local id = 0
		if item.wz ~= nil then
			id = item:wz():id()
		end
		if id < FIRST_OFFERING or id > LAST_OFFERING then
			me:message("필요한 아이템이 아닌 아이템이 있습니다.", Msg.PinkText)
			return nil
		end
		local x, y = item:field_placement():position()
		for i = 1, STATUES do
			local area = wz:area(i)
			if area ~= nil and x >= area.left and x < area.right and y >= area.top and y < area.bottom then
				placed[i] = id - FIRST_OFFERING
				break
			end
		end
	end

	local missing = {}
	for i = 1, STATUES do
		if placed[i] == nil then
			missing[#missing + 1] = i .. "번 석상"
		end
	end
	if #missing > 0 then
		me:message("다음 석상 앞에 제물을 놓아주세요: " .. table.concat(missing, ", "), Msg.PinkText)
		return nil
	end
	return table.concat(placed)
end

local function hint(combo, guess)
	local correct = 0
	local answer_items = { 0, 0, 0, 0 }
	local guess_items = { 0, 0, 0, 0 }
	for i = 1, STATUES do
		local a = tonumber(combo:sub(i, i))
		local g = tonumber(guess:sub(i, i))
		if a == g then
			correct = correct + 1
		else
			answer_items[a + 1] = answer_items[a + 1] + 1
			guess_items[g + 1] = guess_items[g + 1] + 1
		end
	end
	local incorrect = 0
	for k = 1, 4 do
		incorrect = incorrect + math.min(answer_items[k], guess_items[k])
	end
	local unknown = STATUES - correct - incorrect

	local text = ""
	if correct > 0 then
		text = text .. correct .. " 명은 올바른 제물이라 하오.\r\n"
	end
	if incorrect > 0 then
		text = text .. incorrect .. " 명은 틀린 제물이라 하오.\r\n"
	end
	if unknown > 0 then
		text = text .. unknown .. " 명은 모르는 제물이라 하오.\r\n"
	end
	return text
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:map(gq.EXIT_MAP)
			return
		end
		local map = me:map()
		local gate = map:find_reactor_name("watergate")
		if gate ~= nil and gate:state() > 0 then
			me:dialog(npc, "...")
			return
		end

		local combo = sm:get_property("stage3combo")
		if combo == "" then
			sm:set_property("stage3combo", make_combo())
			sm:set_property("stage3attempt", "1")
			me:dialog(npc, INTRO)
			return
		end

		local guess = ground_guess(me)
		if guess == nil then
			me:dialog(npc, "충신 앞에 제대로 놓였는지 확인하고 다시 말을 걸어주세요.")
			return
		end
		if guess == combo or pq.is_gm(me) then
			if gate ~= nil then
				gate:hit(1)
			end
			sm:message("충신의 제물을 통과하였습니다. 비밀문이 열립니다!")
			map:show_effect("quest/party/clear")
			map:play_sound("Party1/Clear")
			gq.gain_gp_once(me, sm, "stage3clear", 25)
			return
		end

		local attempt = tonumber(sm:get_property("stage3attempt")) or 1
		if attempt >= ATTEMPTS then
			sm:set_property("stage3combo", "")
			sm:message("충신의 제물을 바르게 놓는데 실패하여 몬스터가 소환되었습니다.")
			for _ = 1, 5 do
				map:spawn_mob(9300036, math.random(-350, 399), 150)
				map:spawn_mob(9300037, math.random(-350, 399), 150)
			end
			return
		end
		sm:set_property("stage3attempt", tostring(attempt + 1))
		me:dialog(npc, hint(combo, guess) .. "이번이 " .. attempt .. "번째 시도요.")
	end
}
