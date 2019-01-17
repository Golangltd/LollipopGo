package models

import (
	"fmt"
	"strings"
)

func init() {

}

//Query232 查询
func Query232(sql11 string) (json string) {
	dsn := "http://root@10.6.104.67:8181?catalog=mysql&schema=small"
	db, err := sql.Open("presto", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query(sql)

	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list := "["

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("log:", err)
			panic(err.Error())
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])

			cell := fmt.Sprintf(`"%v":"%v"`, columName, value)
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	list = list[0 : len(list)-1]
	list += "]"
	fmt.Println(list)
	return list
}
