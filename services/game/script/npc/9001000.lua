-- NPC name (String.wz/Npc.img.xml): 코-크베어 운영자

local TOWN = 219000000

return {
	on_click = function(me, npc)
		local map = me:map()
		local in_town = map ~= nil and map:template_id() == TOWN
		if in_town then
			local back = me:saved_location("FISHING")
			if back == nil then
				back = 0
			end
			local sel = me:dialog_list(npc, "언제나 코-카 콜라~♪ #b코-크 타운#k에서 즐거운 시간 보냈어? 무엇을 도와줄까?", {
				"이전에 있던 #m" .. back .. "# 마을로 보내주세요.",
			})
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "지금 #b#m" .. back .. "##k 마을로 보내줄까?") then
				me:dialog(npc, "시원한 #b코-크 타운#k에서 즐거운 시간 보내다 가~")
				return
			end
			if id2map(back) == nil then
				me:dialog(npc, "아직 갈 수 없는 곳입니다.")
				return
			end
			me:clear_saved_location("FISHING")
			me:map(back)
			return
		end
		local sel = me:dialog_list(npc, "언제나 코-카 콜라~♪ 새로 오픈된 #b코-크 타운#k으로 놀러가보지 않을래? 시원한 코-카 콜라와 함께해보라구!", {
			"코-크 타운으로 가고 싶어요.",
		})
		if sel == nil then
			return
		end
		if not me:dialog_yes_no(npc, "#b코-크 타운#k으로 놀러가보고 싶어? 언제든지 그곳에 있는 나를 통해 돌아올 수 있어. 지금 당장 출발해볼래?") then
			me:dialog(npc, "흐음. 아쉽구나. 아직 볼일이 남아 있는 모양이지? 나는 언제든지 이곳에서 기다리고 있으니 볼 일이 끝나면 다시 돌아와 줘~")
			return
		end
		me:save_location("FISHING")
		me:map(TOWN)
	end
}
