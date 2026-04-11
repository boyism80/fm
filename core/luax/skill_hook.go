package luax

import "fmt"

func SkillScriptHookName(hook string, skillID uint32) string {
	return fmt.Sprintf("%s_%d", hook, skillID)
}
