package config

var conf = Config{
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
