-- NPC name (String.wz/Npc.img.xml): 가정부 엘마

return {
	on_click = function(me, npc)
		if me:map():wz():id() ~= 200000200 then
			return
		end

		local juice = me:quest(3048)
		if juice:started() or juice:completed() or not me:quest(3047):completed() then
			me:dialog(npc, "주인님이 갑자기 집을 나가버리신 지 벌써 여러달이 지났군요. 저야 할 일이 줄어들어 기분은 좋지만 덕분에 월급이 몇달째 밀려 있어서... 이러다가 월급도 못받고 이 집의 주인이 다른 사람으로 바뀌어 버리는건 아닌지 모르겠어요.\r\n\r\n[엘마의 네펜데스 주스] 퀘스트를 받은 상태에서 퀘스트가 소실된 분들을 위한 버그 픽스용 기능 입니다. 음.. 당신은 이 기능이 필요하지 않을 것 같군요.")
			return
		end

		if next(me:item(4031200)) ~= nil then
			me:exchange({ item = { [4031200] = 1 } }, nil)
		end
		local npc_id = npc
		if type(npc) ~= "number" then
			npc_id = npc:id()
		end
		juice:start(npc_id, true)
	end
}
