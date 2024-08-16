## Golang Database to Struct

### MySQL Database to Golang Struct Conversion Tool Based on GORM v2

This tool allows you to automatically generate Golang structs from a MySQL database using GORM (version 2). It supports the big Camel-Case naming convention and JSON tags.

### Improvements

This project is derived and improved from `https://github.com/xxjwxc/gormt`. Changes include:
- Removal of the GUI for a simpler command-line interface.
- Addition of table prefix configuration options.
- Added support for `GEOMETRY` data type.
- Removal of some unnecessary code for a more streamlined experience.

## Support for GORM Attributes

- **Database Tables and Column Field Annotation**: Support for detailed annotations in generated structs.
- **JSON Tag**: Automatic JSON tag output in structs.
- **GORM.Model**: Inclusion of GORM's built-in model features.
- **PRIMARY_KEY**: Specifies a column as the primary key.
- **UNIQUE**: Marks a column as unique.
- **NOT NULL**: Marks a column as NOT NULL.
- **INDEX**: Allows the creation of indexes, with or without a name; using the same name creates composite indexes.
- **UNIQUE_INDEX**: Similar to INDEX, but creates a unique index.
- **GEOMETRY**: Support for geometry data type.

### Installation

Install the tool using the following Go command:

 `go get github.com/ivan-jorge001/gormt@latest`
 
### Usage

```go
package main

import (
	"github.com/ivan-jorge001/gormt"
	"github.com/ivan-jorge001/gormt/config"
)

func main() {
	dbInfo := config.DBInfo{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "test",
		Type:     0,
	}

	conf := config.Config{
		DBInfo:           dbInfo,
		PkgName:          "schema", 
		OutDir:           "./examples/model", 
		DbTag:            "gorm",
		IsJsonTag:        true,
		IsNullToSqlNull:  true,
		TablePrefix:      "q_",
		StripTablePrefix: true,
		OutFileName:      "schema",
	}

	gormt.ExecuteConfig(&conf)
}
```
