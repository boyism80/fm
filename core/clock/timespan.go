package clock

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var timespanPattern = regexp.MustCompile(`^(?P<neg>-)?((?P<day>\d+)\.)?(?P<hour>\d+):(?P<min>\d{2}):(?P<sec>\d{2})(?:\.(?P<ms>\d+))?$`)

func ParseTimespan(raw string) (time.Duration, error) {
	matches := timespanPattern.FindStringSubmatch(raw)
	if matches == nil {
		return 0, fmt.Errorf("cannot parse timespan: %q", raw)
	}
	group := make(map[string]string, len(matches))
	for i, name := range timespanPattern.SubexpNames() {
		if name == "" || i >= len(matches) {
			continue
		}
		group[name] = matches[i]
	}

	day := 0
	if group["day"] != "" {
		v, err := strconv.Atoi(group["day"])
		if err != nil {
			return 0, fmt.Errorf("cannot parse timespan day: %w", err)
		}
		day = v
	}
	hour, err := strconv.Atoi(group["hour"])
	if err != nil {
		return 0, fmt.Errorf("cannot parse timespan hour: %w", err)
	}
	minute, err := strconv.Atoi(group["min"])
	if err != nil {
		return 0, fmt.Errorf("cannot parse timespan minute: %w", err)
	}
	second, err := strconv.Atoi(group["sec"])
	if err != nil {
		return 0, fmt.Errorf("cannot parse timespan second: %w", err)
	}

	ms := 0
	if group["ms"] != "" {
		in := group["ms"]
		step := 100
		for i := 0; i < len(in) && i < 3; i++ {
			digit := int(in[i] - '0')
			ms += digit * step
			step /= 10
		}
	}

	d := time.Duration(ms)*time.Millisecond +
		time.Duration(second)*time.Second +
		time.Duration(minute)*time.Minute +
		time.Duration(hour)*time.Hour +
		time.Duration(day)*24*time.Hour
	if group["neg"] == "-" {
		d = -d
	}
	return d, nil
}
