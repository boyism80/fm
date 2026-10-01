-- NPC name (String.wz/Npc.img.xml): 샤렌 3세의 영혼

local gq = require("script/lib/guild_quest")

local REST_TEXT = "드디어 샤레니안을 구할 사람을 찾았구나.. 이제.. 편히 잠들 수 있겠어.."

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if sm:get_property("stage4clear") == "true" then
			me:dialog(npc, REST_TEXT)
			return
		end
		if gq.is_leader(me, sm) == false then
			me:dialog(npc, "파티장이 내게 말을 걸어야 한다네...")
			return
		end

		gq.gain_gp_once(me, sm, "stage4clear", 10)
		local gate = me:map():find_reactor_name("ghostgate")
		if gate ~= nil then
			gate:hit(1)
		end
		me:map():show_effect("quest/party/clear")
		me:map():play_sound("Party1/Clear")
		me:dialog(npc, REST_TEXT)
	end
}
