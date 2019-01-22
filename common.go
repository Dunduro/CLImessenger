package main

import "strconv"

const (
	port                int = 7182
	commandCodeLogin        = "CX10001"
	commandCodeLogout       = "CX10002"
	commandCodeSay          = "CX20002"
	commandCodeStatus       = "CX30001"
	commandCodeUserList     = "CX30002"

	chatCommandLogout   = "logout"
	chatCommandSay      = "say"
	chatCommandWisper   = "wisper"
	chatCommandStatus   = "status"
	chatCommandUserList = "users"

	chatShortCommandSay    = "s"
	chatShortCommandWisper = "w"

	appModeServer = "server"
	appModeClient = "client"
	appModeTest   = "test"
)

func genAddress(url string) string {
	return url + ":" + strconv.Itoa(port)
}