package websocket

import "strings"

func GetDockerCommand(object string, action string, objectId string) []string {
	var empty = []string{}
	var base = []string{"docker"}

	switch object {
	case "container":
		permittedActions := []string{"start", "stop", "restart", "remove"}
		permittedActionsJoined := strings.Join(permittedActions, ",")

		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}

		return append(base, object, action, objectId)
	case "image":
		permittedActions := []string{"pull", "rm"}
		permittedActionsJoined := strings.Join(permittedActions, ",")

		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}
		return append(base, object, action, objectId)
	case "TEST":
		return []string{"docker", "run", "ubuntu", "bash", "-c", "while true; do sleep 1 && echo Slept; done"}
	default:
		return empty
	}
}
