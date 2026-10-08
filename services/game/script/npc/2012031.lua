-- NPC name (String.wz/Npc.img.xml): 하프줄<솔>

local MELODY = "004455433221104433221443322100445543322110"

return {
	on_click = function(me, npc)
		local npc_id = npc:id()
		me:play_sound("orbis/sol", true)
		local q = me:quest(3114)
		if q == nil or not q:started() then
			return
		end
		if q:record() == "42" then
			return
		end
		local data = me:records():text("orbis.melody") .. (npc_id - 2012027)
		me:records():set_text("orbis.melody", data)
		if data == MELODY then
			me:message("노래를 정확하게 연주하여 엘리쟈가 잠에 빠져듭니다.", Msg.PinkText)
			q:record("42")
			q:sync_progress()
			me:show_quest_completion(3114)
		elseif string.sub(MELODY, 1, #data) ~= data then
			me:message("연주가 틀렸습니다. 처음부터 다시 연주해 주세요.")
			me:records():set_text("orbis.melody", "")
		end
	end
}
