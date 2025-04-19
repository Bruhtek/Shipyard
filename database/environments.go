package database

type EnvData struct {
	Name    string
	EnvType string
}

func LoadEnvironments() []EnvData {
	data, err := DB.Query("SELECT name, env_type FROM environments")
	if err != nil {
		panic(err)
	}

	defer data.Close()

	var envList []EnvData = make([]EnvData, 0)

	for data.Next() {
		var env EnvData
		err = data.Scan(&env.Name, &env.EnvType)
		if err != nil {
			panic(err)
		}
		envList = append(envList, env)
	}

	return envList
}
