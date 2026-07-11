-- NPC name (String.wz/Npc.img.xml): 첫번째 파이프손잡이

local shuffle_key_quest = 164201
local pipe_progress_quest = 164202
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

function on_click(me, npc)
	local npc_id = npc
	if type(npc) ~= "number" then
		npc_id = npc:id()
	end
	local qr = me:quest(shuffle_key_quest)
	local qr2 = me:quest(pipe_progress_quest)
	if qr == nil or qr2 == nil then
		return
	end

	local key = ensure_quest_record(qr, shuffle_str("123"))
	ensure_quest_record(qr2, "")

	local progress = qr2:record()
	if #progress < 3 then
		qr2:record(progress .. (npc_id - pipe_base_npc))
		progress = qr2:record()
	end

	local len = #progress
	if len == 1 then
		if key:sub(1, 1) == progress:sub(1, 1) then
			me:dialog(npc, "파이프가 날카로운 쇳소리와 함께 오른쪽으로 조금 돌아갔다.", false, false)
		else
			qr2:record("")
			me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
		end
	elseif len == 2 then
		if key:sub(2, 2) == progress:sub(2, 2) then
			me:dialog(npc, "파이프가 날카로운 쇳소리와 함께 왼쪽으로 조금 돌아갔다.", false, false)
		else
			qr2:record("")
			me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
		end
	else
		if key:sub(3, 3) == progress:sub(3, 3) then
			local text = me:dialog_input(npc, "파이프가 아래로 움직이면서 보안장치가 나타났다. 암호를 입력해야 한다.")
			if text == nil then
				return
			end
			if text == "필리아는 내 사랑" then
				qr2:record("")
				qr:record(shuffle_str("123"))
				me:play_portal_sound()
				me:map(261000001, 1)
			else
				qr2:record("")
				qr:record(shuffle_str("123"))
				me:dialog(npc, "비밀번호가 틀린 것 같다. 보안장치가 사라지고 파이프가 원상태로 돌아갔다.", false, false)
			end
		else
			qr2:record("")
			me:dialog(npc, "... 파이프가 날카로운 쇳소리와 함께 돌아가려다가 멈추고, 원래대로 돌아갔다. 다시 처음부터 돌려보자.", false, false)
		end
	end
end
