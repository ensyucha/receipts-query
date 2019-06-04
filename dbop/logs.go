package dbop

import "time"

func WriteLog(logType string, logInfo string, username string) {

	if enableLogG == true {

		logTime := time.Now().Format("2006-1-2 15:04:05")

		stmt, _ := db.Prepare("INSERT INTO logs(logtime, logtype, loginfo, username) VALUES (?,?,?,?);")

		defer stmt.Close()

		_, _ = stmt.Exec(logTime, logType, logInfo, username)
	}
}
