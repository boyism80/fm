-- NPC name (String.wz/Npc.img.xml): 숨겨진 문서들

local DOCUMENT = 4031797

return {
	on_click = function(me, npc)
		local mx, my = me:position()
		local nx, ny = npc:position()
		local dx = mx - nx
		local dy = my - ny
		if dx * dx + dy * dy > 8000 then
			me:message("너무 멀어 조사할 수 없다.", Msg.PinkText)
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local key = tostring(npc:oid())
		if sm:get_property(key) ~= "" then
			return
		end
		me:mkitem(DOCUMENT, 1)
		sm:set_property(key, "1")
	end
}
