package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
)

type CustomerLoggedInEvent struct {
	UserID          int    `json:"customer_id"`
	UserUUID        string `json:"customer_uuid"`
	DeviceID        string `json:"device_id"`
	ApplicationType string `json:"application_type"`
	LoginType       string `json:"login_type"`
}

// var logBuffer []CustomerLoggedInEvent
// var logTmpBuffer []CustomerLoggedInEvent
// var dbClient *sql.DB

func WrapHandleCustomerLoggedIn(event rabbitmq.Event) error {
	var data CustomerLoggedInEvent
	if event.Content == "" {
		fmt.Println("Empty content, skipping")
		return nil
	}
	err := json.Unmarshal([]byte(event.Content), &data)
	if err != nil {
		return err
	}

	return HandleCustomerLoggedIn(data)
}

func HandleCustomerLoggedIn(payload CustomerLoggedInEvent) error {
	fmt.Println("HandleCustomerAccountConfirmEmail", payload)
	return nil
	// logBuffer = append(logBuffer, payload)
	// if len(logBuffer) >= 500 {
	//	var mutex = &sync.Mutex{}
	//	mutex.Lock()
	//	logTmpBuffer = logBuffer
	//	logBuffer = nil
	//	mutex.Unlock()
	//	go SaveToDb()
	//}
	// return nil
}

// func SaveToDb() {
//	dbClient, _ = tools.SetupWriteDB()
//	defer dbClient.Disconnect()
//
//	sqlStr := "INSERT INTO VALUES "
//	var vals []interface{}
//	currentTime := time.Now()
//
//	for _, event := range logTmpBuffer {
//		sqlStr += "(?, ?, ?, ?),"
//		vals = append(vals, event.UserID, currentTime.Format("2006-01-02"), event.LoginType, event.ApplicationType)
//	}
//
//	sqlStr = sqlStr[0 : len(sqlStr)-1]
//	stmt, err := dbClient.Prepare(sqlStr)
//	if err != nil {
//		tools.PrintError(err)
//	}
//	_, err2 := stmt.Exec(vals...)
//	if err2 != nil {
//		tools.PrintError(err2)
//	}
//}
