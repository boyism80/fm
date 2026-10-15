package constant

import "time"

// FamilyAnswerDuration is how long an invite or a summon waits for its answer.
const FamilyAnswerDuration = 60 * time.Second

// FamilyMoveDuration is how long a paid teleport or summon may take before it is refunded instead.
const FamilyMoveDuration = 20 * time.Second

const FamilyMaxJuniors = 2
const FamilyMinJuniorLevel = 10
const FamilyMaxLevelGap = 20

// FamilyCosts[feature] is the reputation a family feature costs.
var FamilyCosts = [...]uint32{300, 500, 700, 800, 1000, 1200, 3000}
