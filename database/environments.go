package database

type EnvData struct {
	Name    string
	EnvType string
	EnvKey  string
}

func LoadEnvironments() []EnvData {
	data, err := DB.Query("SELECT name, env_type, coalesce(env_key, '') FROM environments")
	if err != nil {
		panic(err)
	}

	defer data.Close()

	var envList []EnvData = make([]EnvData, 0)

	for data.Next() {
		var env EnvData
		err = data.Scan(&env.Name, &env.EnvType, &env.EnvKey)
		if err != nil {
			panic(err)
		}
		envList = append(envList, env)
	}

	return envList
}
