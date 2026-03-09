-- Common validation for all active skills. Runs before per-skill on_preactivated.
-- Extend here when adding: event map check (with event system), etc.
-- Return true to allow, false to reject (e.g. "이곳에서 스킬을 사용할 수 없습니다.").

function on_preactivated_common(me, skill)
	return true
end
