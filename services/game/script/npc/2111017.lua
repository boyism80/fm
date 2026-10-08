-- NPC name (String.wz/Npc.img.xml): 첫번째 파이프손잡이

local KEY_RECORD = "magatia.pipe_key"
local PROGRESS_RECORD = "magatia.pipe_progress"
local pipe_base_npc = 2111016

local function shuffle_str(s)
	local t = {}
	for i = 1, #s do
		t[#t + 1] = s:sub(i, i)
	end
	for i = #t, 2, -1 do
		local j = math.random(i)
		t[i], t[j] = t[j], t[i]
	end
	return table.concat(t)
end

return {
	on_click = function(me, npc)
		local npc_id = npc:id()
		local key = me:records():text(KEY_RECORD)
		if key == "" then
			key = shuffle_str("123")
			me:records():set_text(KEY_RECORD, key)
		end

		local progress = me:records():text(PROGRESS_RECORD)
		if #progress < 3 then
			progress = progress .. (npc_id - pipe_base_npc)
			me:records():set_text(PROGRESS_RECORD, progress)
		end

		local len = #progress
		if len == 1 then
			if key:sub(1, 1) == progress:sub(1, 1) then
				me:dialog(npc, "파이프가 날카로운 쇳소리와 함께 오른쪽으로 조금 돌아갔다.", false, false)
			else
				me:records():remove(PROGRESS_RECORD)
				me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
			end
		elseif len == 2 then
			if key:sub(2, 2) == progress:sub(2, 2) then
				me:dialog(npc, "파이프가 날카로운 쇳소리와 함께 왼쪽으로 조금 돌아갔다.", false, false)
			else
				me:records():remove(PROGRESS_RECORD)
				me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
			end
		else
			if key:sub(3, 3) == progress:sub(3, 3) then
				local text = me:dialog_input(npc, "파이프가 아래로 움직이면서 보안장치가 나타났다. 암호를 입력해야 한다.")
				if text == nil then
					return
				end
				if text == "필리아는 내 사랑" then
					me:records():remove(PROGRESS_RECORD)
					me:records():set_text(KEY_RECORD, shuffle_str("123"))
					me:play_portal_sound()
					me:map(261000001, 1)
				else
					me:records():remove(PROGRESS_RECORD)
					me:records():set_text(KEY_RECORD, shuffle_str("123"))
					me:dialog(npc, "비밀번호가 틀린 것 같다. 보안장치가 사라지고 파이프가 원상태로 돌아갔다.", false, false)
				end
			else
				me:records():remove(PROGRESS_RECORD)
				me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
			end
		end
	end
}
