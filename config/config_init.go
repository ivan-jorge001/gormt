package config

var conf = Config{
	DBInfo: DBInfo{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "root",
		Database: "",
	},
	PkgName:          "main",
	OutDir:           "./model",
	DbTag:            "gorm",
	IsJsonTag:        true,
	TablePrefix:      "",
	StripTablePrefix: false,
	SelfTypeDef:      make(map[string]string),
	TableNames:       "",
}

func InitConfig(c *Config) {
	conf = *c
}
