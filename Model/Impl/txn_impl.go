package Impl

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	entity "waas/Model/entity"
)

func GenerateCSV() {
	// currTime := time.Now()
	currTime, _ := time.Parse(time.RFC822, "06 May 21 09:00 IST") // For testing
	if currTime.Hour() != 9 {
		log.Println("Report can be generated ony at 9 AM")
		return
	}

	prevDate := currTime.AddDate(0, 0, -1)

	var transactions []entity.Transaction
	db.Debug().Where("YEAR(`time`) = ?", prevDate.Year()).
		Where("MONTH(`time`) = ?", int(prevDate.Month())).
		Where("DAY(`time`) = ?", prevDate.Day()).
		Preload("Wallet.User").
		Find(&transactions)

	fileName := "Data/" + strings.ReplaceAll(prevDate.Format(time.RFC850), " ", "") + ".csv"
	fileName = strings.ReplaceAll(fileName, ",", "")
	fileName = strings.ReplaceAll(fileName, "-", "")
	fileName = strings.ReplaceAll(fileName, ":", "")

	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Cannot create CSV file", err)
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"userid", "username", "txnid", "txntype", "amount"})
	if err != nil {
		log.Println("Cannot write headers to file", err)
	}
	var row []string
	var txTypeMapper = map[bool]string{true: "Credit", false: "Debit"}
	for _, transaction := range transactions {
		row = []string{}
		row = append(row, strconv.Itoa(transaction.Wallet.UserId))
		row = append(row, transaction.Wallet.User.UserName)
		row = append(row, strconv.Itoa(transaction.ID))
		row = append(row, txTypeMapper[transaction.Type])
		row = append(row, strconv.FormatFloat(transaction.Amount, 'f', 6, 64))
		err = writer.Write(row)
		if err != nil {
			log.Println("Cannot write row to file", err)
		}
	}
}
