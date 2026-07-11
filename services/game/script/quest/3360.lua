-- Quest name (Quest.wz/Quest.img.xml): 비밀번호 인증

local quest_id = 3360

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

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "오! 자네 왔는가? 마침 잘 왔네. 자네를 위해 이 파웬이 비밀통로를 출입할 수 있게 해줄 마스터키를 알아냈다네! 하하하하! 굉장하지 않은가? 어서 굉장하다고 말하게!", false, true)
	if not me:dialog_accept(npc, "자아. 키가 굉장히 길고 복잡하니 잘 기억해 두길 바라겠네. 한 번만 말할 테니 어딘가에 적어 두라고. 준비 되었나?") then
		me:dialog(npc, "빨리 빨리. 외울자신이 없으면 펜이라도 꺼내라고!", false, false)
		return
	end

	local key = shuffle_str("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"):sub(1, 10)
	local q7061 = me:quest(7061)
	if q7061 ~= nil then
		if q7061:wz() == nil then
			q7061:start(key)
		else
			q7061:start(npc, key)
		end
	end
	local q7062 = me:quest(7062)
	if q7062 ~= nil then
		if q7062:wz() == nil then
			q7062:start("00")
		else
			q7062:start(npc, "00")
		end
	end
	q:start(npc, "0")
	me:dialog(npc, "키번호는 #b" .. key .. "#k이네. 잊지 않았겠지? 이 키를 비밀통로 입구에 입력하면 비밀통로를 자유롭게 이용할 수 있을 거야. ", false, false)
end
