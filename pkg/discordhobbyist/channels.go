package discordhobbyist

import "strings"

func GetChannelKey(group, channel string) string {
	return "/" + strings.ToLower(group) + "/" + strings.ToLower(channel)
}
