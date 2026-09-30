-- NPC name (String.wz/Npc.img.xml): 조수 레드 - 몬스터 카니발

local HUB_MAP = 980000000

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		me:message("누군가가 나가기 엔피시를 클릭하여 모두 나가집니다.", Msg.Notice)
		local match = carnival.map_match(tonumber(sm:id()))
		sm:finish(HUB_MAP)
		if match ~= nil then
			match:finish()
		end
	end
}
