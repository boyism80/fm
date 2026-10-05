local SCRIPT = "script/quest/3706.lua"
local CRY_MAPS = { 240000000, 240040611 }
local CRY_NOTICE = "나인스피릿 아기용의 힘찬 울음소리를 듣자 신비로운 힘이 솟아오른다."

return {
	on_buff_item = function(me, item_id)
		for _, map_id in ipairs(CRY_MAPS) do
			run_on_map(map_id, SCRIPT, "cry", item_id)
		end
	end,

	cry = function(map, item_id)
		map:message(CRY_NOTICE)
		for _, ch in pairs(map:characters()) do
			ch:use_item(item_id)
		end
	end,
}
