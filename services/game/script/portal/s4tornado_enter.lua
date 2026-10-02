local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		if me:class() ~= 412 then
			return
		end
		local tornado = me:quest(6230)
		local started = tornado:started()
		local completed = tornado:completed()
		local fresh = started == false and completed == false
		local allowed = (fresh and pq.has_item(me, 4001110))
			or started
			or (completed and me:quest(6231):started() == false and me:quest(6231):completed() == false)
		if allowed == false then
			return
		end
		me:play_portal_sound()
		me:map(922020200, 0)
	end
}
