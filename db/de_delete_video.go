package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 用户请求 apiserver 把要删除的视频id 写入 `video_del_info` 表
// Dispatcher 批量读取待删除视频id 减少数据库压力
// Executor 真正删除了待删除的视频后,要把  `video_del_info` 表相应数据删除

func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_id FROM video_del_info LIMIT ?`)
	var ids []string
	if err != nil {
		return ids, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare(`DELETE FROM video_del_info WHERE video_id=?`)
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}
	defer stmtDel.Close()

	return nil
}

func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO video_del_rec (video_id) VALUES(?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}