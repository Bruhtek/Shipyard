package actions

import "strings"

func GetDockerCommand(object string, action string, objectId string) []string {
	var empty []string
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
	case "stack":
		permittedActions := []string{"up", "stop", "down", "restart", "pull", "update"}
		permittedActionsJoined := strings.Join(permittedActions, ",")
		if !strings.Contains(permittedActionsJoined, action) {
			return empty
		}

		if action == "update" {
			first := GetDockerCommand(object, "pull", objectId)
			second := GetDockerCommand(object, "up", objectId)
			return append(append(first, "&&"), second...)
		}

		base = append(base, "compose")
		// in this case, the objectId is actually the config file path
		base = append(base, "-f", objectId)
		base = append(base, action)

		if action == "up" {
			// for "up" action, we add --detach and --remove-orphans flags
			base = append(base, "-d", "--remove-orphans")
		}
		return base
	case "TEST":
		//return append(base, "container", "remove", "221f468ab0c3", "700a4d7b2b60")
		return empty
	default:
		return empty
	}

	return append(append(base, object, action), ids...)
}
