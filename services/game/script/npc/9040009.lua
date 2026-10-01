-- NPC name (String.wz/Npc.img.xml): 성문 문지기

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local LAST_PHASE = 3
local DISPLAY_DELAY_MS = 5000

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:map(gq.EXIT_MAP)
			return
		end
		local map = me:map()
		local gate = map:find_reactor_name("statuegate")
		if gq.is_leader(me, sm) == false or gate == nil or gate:state() > 0 then
			me:dialog(npc, "...")
			return
		end

		local status = sm:get_property("stage1status")
		local phase = tonumber(sm:get_property("stage1phase")) or 1
		if status == "active" then
			sm:set_property("stage1status", "waiting")
			if sm:get_property("stage1combo") ~= sm:get_property("stage1guess") and pq.is_gm(me) == false then
				sm:message("성문 문지기의 시험에 실패하였습니다.")
				sm:set_property("stage1phase", "1")
				return
			end
			if phase < LAST_PHASE then
				sm:set_property("stage1phase", tostring(phase + 1))
				map:message("성문 문지기의 시험 중 " .. phase .. "번째 시험을 통과하였습니다.")
				me:dialog(npc, "... 다음 행적을 쫓을 준비를 하라. 준비가 되면 다시 내게 말을 걸도록...")
				return
			end
			gate:hit(1)
			sm:message("성문 문지기의 시험을 전부 통과하였습니다.")
			map:show_effect("quest/party/clear")
			map:play_sound("Party1/Clear")
			gq.gain_gp_once(me, sm, "stage1clear", 15)
			return
		end
		if status ~= "" and status ~= "waiting" then
			me:dialog(npc, "...")
			return
		end

		sm:set_property("stage1phase", tostring(phase))
		local text = "... 다음 행적을 쫓을 준비를 하라."
		if phase == 1 then
			text = "... 도전을 하겠다..? 내가 이 석상들을 움직일 것이다. 그렇다면 내가 움직인 석상들의 빛나는 행적을 알아 맞추도록.."
		end
		if me:dialog(npc, text) == false then
			return
		end

		local statues = {}
		for oid, reactor in pairs(map:reactors()) do
			if reactor:name() ~= "statuegate" then
				statues[#statues + 1] = oid
			end
		end
		for i = #statues, 2, -1 do
			local j = math.random(i)
			statues[i], statues[j] = statues[j], statues[i]
		end
		local pending = ""
		for i = 1, math.min(gq.combo_length(sm), #statues) do
			pending = pending .. statues[i] .. ","
		end
		sm:set_property("stage1status", "display")
		sm:set_property("stage1combo", "")
		sm:set_property("stage1pending", pending)
		sm:after(DISPLAY_DELAY_MS, "on_statue_display")
	end
}
