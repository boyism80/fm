-- NPC name (String.wz/Npc.img.xml): 비밀통로

local function ensure_quest_record(q, default)
	if q == nil then
		return nil
	end
	if not q:started() and not q:completed() then
		q:start(default)
		return default
	end
	local r = q:record()
	if r == nil or r == "" then
		q:record(default)
		return default
	end
	return r
end

return {
	on_click = function(me, npc)
		local q3360 = me:quest(3360)
		if q3360 == nil or (not q3360:started() and not q3360:completed()) then
			return
		end

		local side = 0
		local map = me:map()
		if map ~= nil and map:wz().id == 261020200 then
			side = 1
		end

		local qr7062 = me:quest(7062)
		if qr7062 == nil then
			return
		end
		local access = ensure_quest_record(qr7062, "00")
		local cur_code = access:sub(side + 1, side + 1)

		if cur_code == "0" then
			local text = me:dialog_input(npc, "비밀번호를 입력하시오.")
			if text == nil then
				return
			end
			local qr7061 = me:quest(7061)
			if qr7061 == nil then
				return
			end
			local key = qr7061:record()
			if key == nil then
				key = ""
			end
			if text == key then
				if side == 0 then
					qr7062:record("1" .. access:sub(2, 2))
				else
					qr7062:record(access:sub(1, 1) .. "1")
				end
				me:notice("보안장치가 해제되었습니다. 출입허가명단에 등록되었습니다.")
				local updated = qr7062:record()
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
