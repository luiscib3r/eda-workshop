package nats

import "strings"

func StreamName(channel string) string {
	name := strings.ReplaceAll(channel, ".", "_")
	return strings.ToUpper(name)
}
