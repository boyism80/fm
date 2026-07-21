-- Portal (old/scripts/portal/hontale_Bopen.js): 생명의동굴 미로방 진행

local pq = require("script/lib/party_quest")

function on_enter(me)
	local sm = me:state_machine()
	if sm == nil then
		me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local wz = map:wz()
	if wz == nil then
		return
	end
	local map_id = wz.id
	local progress = tonumber(sm:get_property("stage1progress")) or 0

	if map_id == 240050101 then
		if progress == 0 then
			me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(240050102)
		return
	end
	if map_id == 240050102 then
		if progress <= 1 then
			me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(240050103)
		return
	end
	if map_id == 240050103 then
		if progress <= 2 then
			me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(240050104)
		return
	end
	if map_id == 240050104 then
		if progress <= 3 then
			me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(240050105)
		return
	end
	if map_id == 240050105 then
		if not pq.has_item(me, 4001092, 1) then
			me:notice("알 수 없는 힘으로 포탈이 막혀있어 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		pq.remove_all(4001092, me)
		for _, p in ipairs(sm:players()) do
			if p ~= nil then
				p:notice("붉은 열쇠의 힘으로 이동되었습니다.", Msg.PinkText)
			end
		end
		me:play_portal_sound()
		pq.party_warp(sm, 240050100)
		sm:set_property("stage1progress", tostring(progress + 1))
	end
end
