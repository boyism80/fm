-- NPC name (String.wz/Npc.img.xml): 아몬

local ENTRY_MAP = 211042300

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		local msg
		if wz ~= nil and wz.id == 280030000 then
			msg = "이 곳에서 나가겠는가? 혹여나 소환 엔피시가 없으면 운영자를 호출해주게나."
		else
			msg = "이 곳에서 나가겠는가? 다음번에 들어올 때는 처음부터 다시 시도해야 한다네."
		end
		if not me:dialog_yes_no(npc, msg) then
			return
		end
		me:map(ENTRY_MAP)
	end
}
