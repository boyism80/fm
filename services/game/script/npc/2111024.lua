-- NPC name (String.wz/Npc.img.xml): 비밀통로

return {
	on_click = function(me, npc)
		local q3360 = me:quest(3360)
		if q3360 == nil or (not q3360:started() and not q3360:completed()) then
			return
		end

		local side = 0
		local map = me:map()
		if map ~= nil and map:wz():id() == 261020200 then
			side = 1
		end

		local access = me:records():text("magatia.secret_access")
		if access == "" then
			access = "00"
		end
		local cur_code = access:sub(side + 1, side + 1)

		if cur_code == "0" then
			local text = me:dialog_input(npc, "비밀번호를 입력하시오.")
			if text == nil then
				return
			end
			if text == me:records():text("magatia.secret_key") then
				local updated = access:sub(1, 1) .. "1"
				if side == 0 then
					updated = "1" .. access:sub(2, 2)
				end
				me:records():set_text("magatia.secret_access", updated)
				me:message("보안장치가 해제되었습니다. 출입허가명단에 등록되었습니다.")
				if updated == "11" then
					q3360:record("1")
					q3360:sync_progress()
					me:show_quest_completion(3360)
				end
			else
				me:dialog(npc, "... 비밀번호가 틀렸습니다.", false, false)
			end
		else
			me:play_portal_sound()
			if side == 0 then
				me:map(261030000, 2)
			else
				me:map(261030000, 1)
			end
		end
	end
}
