package model

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// query, total_time, calls, rows, mean_time
type QeryState struct {
	Query     string `json:"query"`
	TotalTime string `json:"total_time"`
	Calls     string `json:"calls"`
	Rows      string `json:"rows"`
	MeanTime  string `json:"mean_time"`
}

func (q *QeryState) GetQueryState(db *gorm.DB, pageSize int, offset int, filter string, sort string) ([]map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New("database not connected !")
	}
	var err error
	var f []map[string]interface{}
	var sql string
	if sort == "DESC" && filter == "*" {
		sql = `
		SELECT dbid,query,rows,calls,
		max_exec_time,mean_exec_time,total_exec_time 
		FROM pg_stat_statements
		ORDER BY max_exec_time DESC 
		LIMIT ?
		OFFSET ?;
		`
		err = db.Debug().Raw(sql, pageSize, offset).Find(&f).Error
	} else if sort == "ASC" && filter == "*" {
		sql = `
		SELECT dbid,query,rows,calls,
		max_exec_time,mean_exec_time,total_exec_time 
		FROM pg_stat_statements
		ORDER BY max_exec_time ASC 
		LIMIT ?
		OFFSET ?;
		`
		err = db.Debug().Raw(sql, pageSize, offset).Find(&f).Error
	} else if sort == "ASC" {
		sql = `
		SELECT dbid,query,rows,calls,
		max_exec_time,mean_exec_time,total_exec_time 
		FROM pg_stat_statements
		WHERE query LIKE ? OR query LIKE ?
		ORDER BY max_exec_time ASC 
		LIMIT ?
		OFFSET ?;
		`
		err = db.Debug().Raw(sql, filter, strings.ToLower(filter), pageSize, offset).Find(&f).Error
	} else {
		sql = `
		SELECT dbid,query,rows,calls,
		max_exec_time,mean_exec_time,total_exec_time 
		FROM pg_stat_statements
		WHERE query LIKE ? OR query LIKE ?
		ORDER BY max_exec_time DESC
		LIMIT ?
		OFFSET ?;
		`
		err = db.Debug().Raw(sql, filter, strings.ToLower(filter), pageSize, offset).Find(&f).Error
	}

	if err != nil {
		return nil, err
	}
	return f, nil
}
