-- NPC name (String.wz/Npc.img.xml): 알리

local pq = require("script/lib/party_quest")

local ENTRY_MAP = 211042300
local APPROVAL_RECORD = "zakum.approval"
local STAGE1_RECORD = "zakum.stage1"
local STAGE_COMPLETED = 2
local PAPER = 4001015
local KEY = 4001016
local FIRE_ORE = 4001018

return {
	on_click = function(me, npc)
		if me:records():get(STAGE1_RECORD) == STAGE_COMPLETED then
			if not me:dialog(npc, "1단계를 훌륭히 해냈군 그래. 좋아... 자네를 #b아도비스#k가 있는 밖으로 내보내 주겠네. 그 전에! 이곳에서 얻은 특수한 아이템들은 밖으로 가지고 나갈 수 없게 되어 있네. 내보내 주면서 강제로 그 물건을 빼앗을 수도 있으니 참고해 주게나. 그럼 잘가게!", true, true) then
				return
			end
			me:records():set_text(APPROVAL_RECORD, "Zakum1Clear")
		else
			if not me:dialog(npc, "도중에 포기한 모양이로군. 좋아... 자네를 지금 당장 밖으로 내보내 주겠네. 하지만 그 전에! 이곳에서 얻은 특수한 아이템들은 밖으로 가지고 나갈 수 없게 되어 있네. 내보내 주면서 강제로 다 빼앗을 수도 있으니 참고해 주게나. 그럼 잘가게!", true, true) then
				return
			end
		end
		pq.remove_all(PAPER, me)
		pq.remove_all(KEY, me)
		pq.remove_all(FIRE_ORE, me)
		me:map(ENTRY_MAP)
	end
}
