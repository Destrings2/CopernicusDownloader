package database

func (c *Connection) Query(query string, args ...interface{}) ([]map[string]interface{}, error, []string) {
	rows, err := c.Database.Query(query, args...)
	if err != nil {
		return nil, err, nil
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err, nil
	}

	var result []map[string]interface{}

	for rows.Next() {
		var row = make(map[string]interface{})
		var values = make([]interface{}, len(columns))
		var scanValues = make([]interface{}, len(columns))

		for i := range values {
			scanValues[i] = &values[i]
		}

		if err := rows.Scan(scanValues...); err != nil {
			return nil, err, nil
		}

		for i, column := range columns {
			row[column] = values[i]
		}

		result = append(result, row)
	}

	return result, nil, columns
}

func (c *Connection) Execute(query string, args ...interface{}) (int64, error) {
	result, err := c.Database.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (c *Connection) SetRequested(uuid string) error {
	_, err := c.Database.Exec("UPDATE `file` SET `requestedDate` = NOW() WHERE `uuid` = ?", uuid)
	return err
}
