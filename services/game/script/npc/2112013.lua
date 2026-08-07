-- NPC name (String.wz/Npc.img.xml): 조사 결과

local pq = require("script/lib/party_quest")
local rj = require("script/lib/romeo_juliet")

local LETTER = 4001130
local STAGE1_EXP = 5000

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil or map:wz().id ~= 926110000 then
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local mx, my = me:position()
		local nx, ny = npc:position()
		local dx = mx - nx
		local dy = my - ny
		if dx * dx + dy * dy > 5000 then
			me:dialog(npc, "조사하기에는 너무 멀리 있다.")
			return
		end
		local key = "stage1_" .. tostring(npc:oid())
		local prop = sm:get_property(key)
		local stage1 = sm:get_property("stage1")
		if stage1 == "" or stage1 == "0" then
			if prop == "" then
				sm:set_property(key, "1")
				local rand = math.random()
				if rand < 0.2 then
					me:dialog(npc, "경험치를 획득했지만 아무것도 찾지는 못했다.")
					me:exchange({}, { exp = 500 })
				elseif rand < 0.5 then
					me:dialog(npc, "500메소를 발견했다.")
					me:exchange({}, { meso = 500 })
				else
					me:dialog(npc, "아무것도 없는 것 같다.")
				end
				return
			end
			if prop == "1" then
				me:dialog(npc, "이미 조사한 곳이다.")
				return
			end
			if prop == "2" then
				rj.hit_all_reactors(map)
				map:message(me:name() .. "님이 스위치를 누르자 특수한 포탈이 나타났다.")
				rj.clear_fx(map)
				sm:set_property("stage1", "1")
				pq.party_exp(sm, STAGE1_EXP)
				return
			end
			if prop == "3" then
				me:mkitem(LETTER, 1)
				me:dialog(npc, "무언가 편지를 발견했다.")
				sm:set_property(key, "1")
				return
			end
			return
		end
		if prop == "3" then
			me:mkitem(LETTER, 1)
			me:dialog(npc, "무언가 편지를 발견했다.")
			sm:set_property(key, "1")
			return
		end
		me:dialog(npc, "이미 다음 스테이지로 가는 포탈이 열려있다.")
	end
}
