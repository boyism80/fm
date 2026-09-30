-- NPC name (String.wz/Npc.img.xml): 가위바위보 운영자

local function input_number(me, npc, text, default, min, max)
	local n = tonumber(me:dialog_input(npc, text))
	if n == nil then
		n = default
	end
	return math.max(min, math.min(max, math.floor(n)))
end

return {
	on_click = function(me, npc)
		if me:role() < ROLE.Admin then
			me:dialog(npc, "여행은 즐거우세요?")
			return
		end

		local want = input_number(me, npc, "제작할 아이템의 #r코드#k를 입력해주세요.\r\n", 1112400, 0, 2000000)
		local allstat = input_number(me, npc, "#i" .. want .. "#의 #b올스텟#k 옵션을 선택해주세요.", 100, 0, 32767)
		local damage = input_number(me, npc, "#i" .. want .. "#의 #b공격력#k과 #b마력#k을 선택해주세요.\r\n(최대 32,767 공격력은 최대 적용 2000)", 100, 0, 32767)

		local reward = {
			item = { [want] = 1 },
			bonus = { str = allstat, dex = allstat, int = allstat, luk = allstat, watk = damage, matk = damage },
		}
		if me:exchange(nil, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "존재하지 않는 아이템코드네요.")
			return
		end
		me:dialog(npc, "#fUI/UIWindow.img/QuestIcon/4/0#\r\n#i" .. want .. "# #b#z" .. want .. "##k\r\nSTR : " .. allstat .. "\r\nDEX : " .. allstat .. "\r\nINT : " .. allstat .. "\r\nLUK : " .. allstat .. "\r\n공격력 : " .. damage .. "\r\n마　력 : " .. damage)
	end
}
