-- NPC name (String.wz/Npc.img.xml): 하프줄<시>

local MELODY = "004455433221104433221443322100445543322110"

function on_click(me, npc)
	local npc_id = npc
	if type(npc) ~= "number" then
		npc_id = npc:id()
	end
	me:play_sound("orbis/si", true)
	local q = me:quest(3114)
	if q == nil or not q:started() then
		return
	end
	if q:record() == "42" then
		return
	end
	local seq = me:quest(103114)
	if seq == nil then
		return
	end
	if not seq:started() and not seq:completed() then
		seq:start("")
	end
	local data = seq:record()
	if data == nil then
		data = ""
	end
	data = data .. (npc_id - 2012027)
	seq:record(data)
	seq:sync_progress()
	if data == MELODY then
		me:notice("노래를 정확하게 연주하여 엘리쟈가 잠에 빠져듭니다.", Msg.PinkText)
		q:record("42")
		q:sync_progress()
		me:show_quest_completion(3114)
	elseif string.sub(MELODY, 1, #data) ~= data then
		me:notice("연주가 틀렸습니다. 처음부터 다시 연주해 주세요.")
		seq:record("")
		seq:sync_progress()
	end
end
