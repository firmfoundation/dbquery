package model

import "gorm.io/gorm"

// query, total_time, calls, rows, mean_time
type QeryState struct {
	Query     string `json:"query"`
	TotalTime string `json:"total_time"`
	Calls     string `json:"calls"`
	Rows      string `json:"rows"`
	MeanTime  string `json:"mean_time"`
}

func (q *QeryState) GetQueryState(db *gorm.DB, pageSize int, offset int) ([]map[string]interface{}, error) {
	var f []map[string]interface{}
	sql := `
	SELECT dbid,query,rows,calls,
	max_exec_time,mean_exec_time,total_exec_time 
	FROM pg_stat_statements
	ORDER BY max_exec_time DESC 
	LIMIT ?
	OFFSET ?;
	`
	err := db.Debug().Raw(sql, pageSize, offset).Find(&f).Error
	if err != nil {
		return nil, err
	}
	return f, nil
}
