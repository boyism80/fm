-- NPC name (String.wz/Npc.img.xml): 조수 블루 - 몬스터 카니발

local HUB_MAP = 980000000

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		me:notice("누군가가 나가기 엔피시를 클릭하여 모두 나가집니다.", Msg.Notice)
		local match = carnival.map_match(tonumber(sm:id()))
		if match ~= nil then
			match:finish()
		end
		sm:finish(HUB_MAP)
	end
}
