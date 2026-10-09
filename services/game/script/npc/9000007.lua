-- NPC name (String.wz/Npc.img.xml): 천지

local ANCIENT_SCROLL = 4031019
local MATERIAL = 4000008
local MATERIAL_COUNT = 50
local MAGIC_BOX = 4031017
local MAGIC_BOX_MINUTES = 24 * 60

local function count(me, item_id)
	local n = 0
	for _, it in pairs(me:item(item_id)) do
		n = n + it:count()
	end
	return n
end

return {
	on_click = function(me, npc)
		if count(me, ANCIENT_SCROLL) < 1 then
			me:dialog(npc, "건드리지 마라~ 다친다~")
			return
		end
		if me:dialog(npc, "호오, #b#t4031019##k를 구해오다니. 너 보통내기가 아닌걸? 흐음.. 너, 그걸 나에게 주지 않을래? 나에게 준다면 #b#t4031017##k를 주도록 하지.", false, true) == false then
			return
		end
		if me:dialog(npc, "하지만, 이 문서를 해석하려면.. 아무래도 재료가 약간 필요할 것 같아. 재료는.. 중앙 던전에서 #b#t4000008# 50개#k 정도만 구해오면 될 것 같은데..", true, true) == false then
			return
		end
		if count(me, MATERIAL) < MATERIAL_COUNT then
			me:dialog(npc, "재료를 구하면 다시 나에게 찾아오도록 해.")
			return
		end
		if me:dialog(npc, "엇? 그건 #b#t4000008# 50개#k! 좋아. 재료도 모였겠다, 이제 너에게 보상을 주도록 하지.", true, true) == false then
			return
		end
		if count(me, MAGIC_BOX) > 0 then
			me:dialog(npc, "어이? 이미 #t4031017#를 갖고 있는거 아니야? 이미 갖고 있으면 줄 수 없다구.")
			return
		end

		local code = me:exchange(
			{ item = { [MATERIAL] = MATERIAL_COUNT, [ANCIENT_SCROLL] = 1 } },
			{ item = { [MAGIC_BOX] = 1 }, period = { [MAGIC_BOX] = MAGIC_BOX_MINUTES } }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간이 부족한거 아니야? 다시 한번 확인해 봐.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "이 매직박스는 1일 이내에 열어야 해. 그렇지 않으면 사라지고 말지. 지금 받은 즉시 열어 보도록 해. #b커닝시티#k의 #r몽땅따#k 아저씨에게 가 보면 이 상자를 열 수 있을거야.")
	end
}
