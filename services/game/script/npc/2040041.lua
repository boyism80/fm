-- NPC name (String.wz/Npc.img.xml): 아쿠아 벌룬

local pq = require("script/lib/party_quest")

local HINT = "#b하나, 3, 3, 2, 가운데, 1, 셋, 3, 3, 왼쪽, 둘, 3, 1, 일, ?"

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:dialog(npc, "오류가 발생했어요.")
			return
		end
		if not pq.is_leader(me) then
			me:dialog(npc, "안녕하세요. 여섯번째 스테이지에 오신 것을 환영합니다. 이 곳에는 숫자가 적힌 상자가 있는데 올바른 상자 위로 올라가 ↑ 키를 누르면 다음 상자들로 이동됩니다. 제가 파티장에게 단 두번만 통과에 관한 힌트를 드릴텐데 파티장은 그 힌트를 잘 기억하고 있다가 힌트대로 모든 파티원이 차례 차례 위로 올라가면 됩니다.\r\n맨 위까지 올라가면 다음 스테이지로 갈 수 있는 포탈이 있을 겁니다. 포탈을 통해 파티원 전원이 다음 스테이지로 이동하면 클리어 한 것이 됩니다. 그럼 힘내주세요!")
			return
		end
		local val = sm:get_property("stage6")
		if val == "" then
			sm:set_property("stage6", "1")
			me:dialog(npc, "안녕하세요. 여섯번째 스테이지에 오신 것을 환영합니다. 이 곳에는 숫자가 적힌 상자가 있는데 올바른 상자 위로 올라가 ↑ 키를 누르면 다음 상자들로 이동됩니다. 제가 파티장에게 단 두번만 통과에 관한 힌트를 드릴텐데 파티장은 그 힌트를 잘 기억하고 있다가 힌트대로 모든 파티원이 차례 차례 위로 올라가면 됩니다.\r\n맨 위까지 올라가면 다음 스테이지로 갈 수 있는 포탈이 있을 겁니다. 포탈을 통해 파티원 전원이 다음 스테이지로 이동하면 클리어 한 것이 됩니다. 올바른 상자를 기억하는 것이 관건이겠군요. 그럼 힌트를 드릴테니 잘 외우도록 하세요~!\r\n\r\n"
				.. HINT)
		elseif val == "1" then
			sm:set_property("stage6", "2")
			me:dialog(npc, "다시 한번 힌트를 드릴테니 잘 보고 외우세요~! \r\n\r\n" .. HINT)
		else
			me:dialog(npc, "힌트를 두번 모두 드렸답니다. 더 이상 힌트를 드릴 수 없어요! 파티원들과 함께 힘을 합쳐 문제를 해결해보세요.")
		end
	end
}
