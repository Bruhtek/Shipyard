package websocket

import "strings"

func GetDockerCommand(object string, action string, objectId string) []string {
	var empty = []string{}
	var base = []string{"docker"}

	var ids = strings.Split(objectId, ",")

	switch object {
	case "container":
		permittedActions := []string{"start", "stop", "restart", "remove"}
		permittedActionsJoined := strings.Join(permittedActions, ",")

		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}
	case "image":
		permittedActions := []string{"pull", "rm"}
		permittedActionsJoined := strings.Join(permittedActions, ",")

		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}
	case "network":
		permittedActions := []string{"remove"}
		permittedActionsJoined := strings.Join(permittedActions, ",")

		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}
	case "TEST":
		//return append(base, "container", "remove", "221f468ab0c3", "700a4d7b2b60")
		return empty
	default:
		return empty
	}

	return append(append(base, object, action), ids...)
}
