-- NPC name (String.wz/Npc.img.xml): 개찰구

local SUBWAY = 103000101
local SQUARE_RIDE = 103000300
local SQUARE = 103000310
local RIDE_SECONDS = 20

local SITES = {
	{ name = "공사장 B1", level = 20, item = 4031036, map = 103000900 },
	{ name = "공사장 B2", level = 30, item = 4031037, map = 103000903 },
	{ name = "공사장 B3", level = 40, item = 4031038, map = 103000906 },
}

return {
	on_click = function(me, npc)
		local level = me:level()
		if level < 20 then
			me:dialog(npc, "입장할 수 있는 구간이 없습니다..")
			return
		end
		local options = {
			"#b커닝시티 지하철 #r(주의 : 스티지, 레이스 등 서식)#b",
			"커닝 스퀘어 백화점(지하철 탑승)",
		}
		local sites = {}
		for _, site in ipairs(SITES) do
			if level >= site.level then
				sites[#sites + 1] = site
				options[#options + 1] = site.name
			end
		end
		local sel = me:dialog_list(npc, "안녕하세요~ 커닝시티 지하철 입니다. 입장하시고 싶은 구간을 선택하세요.\r\n", options)
		if sel == nil then
			return
		end
		if sel == 1 then
			me:play_portal_sound()
			me:map(SUBWAY, 3)
			return
		end
		if sel == 2 then
			me:play_portal_sound()
			local function on_arrive(me)
				me:clock(RIDE_SECONDS, function(me)
					me:map(SQUARE)
				end)
			end
			me:map(SQUARE_RIDE, 0, { callback = on_arrive })
			return
		end
		local site = sites[sel - 2]
		if site == nil then
			return
		end
		if id2map(site.map) == nil then
			me:dialog(npc, "지금은 그 구간으로 입장하실 수 없습니다.")
			return
		end
		if me:exchange({ item = { [site.item] = 1 } }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "입장권이 없으신 것 같습니다. #b지하철 공익요원#k 에게서 입장권을 구매해주세요.")
			return
		end
		me:map(site.map, "sp")
	end
}
