local SHOW = 223030210

return {

	on_enter = function(me)
		local quest = me:quest(31326)
		local ok, count = run_on_map(SHOW, "script/lib/party_quest.lua", "character_count")
		if (quest:started() or quest:completed()) and ok and count == 0 then
			me:open_npc(2190002)
			return
		end
		me:message("길이 막혀있거나 이미 누군가 공연장에 입장해 있습니다.", Msg.PinkText)
	end
}
