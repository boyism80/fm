-- NPC name (String.wz/Npc.img.xml): 요정 마르

local pet_revive = require("script/lib/pet_revive")

return {
	on_click = function(me, npc)
		pet_revive.talk(me, npc, 2049,
			"저는 펫을 되살려 주거나, 기존 펫의 능력치를 새로운 펫으로 옮겨주는 일을 맡고 있는 요정 마르랍니다.",
			"준비해 오셔야 할 재료는 #b생명의 물#k 과 #r생명의 주문서#k 입니다. 생명의 물은.. 구하기 힘들지만.. 어떤 상점에서 현재 판매하고 있다고 하더군요.  생명의 주문서가 문제인데.. 이것은 헤네시스의 #b조련사 바르토스#k에게 찾아가 보시면 뭔가 알 수 있으실거에요.. 후훗")
	end
}
