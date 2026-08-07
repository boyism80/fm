-- NPC name (String.wz/Npc.img.xml): 조나단

local function start_trial(me, npc, group_name)
	local group = state_machine(group_name)
	if group == nil then
		me:dialog(npc, "오류가 발생해서 시험을 치를 수 없네.")
		return
	end
	if group:get_property("noEntry") == "true" then
		me:dialog(npc, "흠, 이미 누군가가 시험을 치르는 중이네. 나중에 다시 시도하게나.")
		return
	end
	local sm, err = group:start_solo(me)
	if sm == nil then
		me:dialog(npc, "오류가 발생해서 시험을 치를 수 없네.")
		if err ~= nil then
			log(group_name .. " start_solo:", err)
		end
	end
end

return {
	on_click = function(me, npc)
		local quest = me:quest(6400)
		if quest == nil or not quest:started() then
			me:dialog(npc, "크흠, 아무 볼일도 없다면 이만 물러가게나.")
			return
		end
		local completion = me:quest(6401)
		if completion ~= nil and completion:started() then
			me:dialog(npc, "커흠. 이미 나의 모든 시험을 통과했구만? 축하하네~")
			return
		end
		local progress_quest = me:quest(116400)
		if progress_quest == nil then
			return
		end
		if not progress_quest:started() then
			me:dialog(npc, "하하~ 드디어 만나보게 되었군. 난 갈매기들의 제왕 조나단 3세일세. 해적이라면 누구나 거쳐야 할 의식을 위해 자네를 불렀네.")
			me:dialog(npc, "자네 스스로 지혜롭다고 생각되면 나에게 다시 말을 걸게. 알겠나?")
			progress_quest:start("q1")
			return
		end
		local progress = progress_quest:record()
		if progress == "q1" then
			local answer = me:dialog_input(npc, "이 세상에서 가장 강하고 용감한 사람은 누구일까?")
			if answer ~= "조나단" then
				me:dialog(npc, "에긍.. 쯧쯧쯧, 틀렸어. 다시 한번 생각해 보라구.")
				return
			end
			me:dialog(npc, "우하하~! 나를 잠깐 봤는데도 금방 알아보다니, 역시 나를 실망시키지 않는군!")
			progress_quest:record("q2")
			return
		end
		if progress == "q2" then
			me:dialog(npc, "빈 방에 아홉 명의 가짜 바트와 한 명의 진짜 바트가 있을 걸세. 진짜 바트를 찾아보게.")
			if me:dialog_yes_no(npc, "시험을 치를 준비가 되었는가?") then
				start_trial(me, npc, "air_strike")
			end
			return
		end
		if progress == "q22" then
			me:dialog(npc, "바트를 쉽게 찾아냈구만. 다음은 샤레니안 성문의 문지기에게 시험을 받게.")
			start_trial(me, npc, "air_strike_2")
		end
	end
}
