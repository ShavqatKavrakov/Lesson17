package main

import (
	"Lesson17_ful_export/pkg/wallet"
	"log"
)

func main() {
	srv := &wallet.Service{}
	_, err := srv.RegisterAccount("918925874")
	if err != nil {
		log.Print(err)
	}
	_, err = srv.Deposit(1, 10_00000)
	if err != nil {
		log.Print(err)
	}
	_, err = srv.RegisterAccount("987026424")
	if err != nil {
		log.Print(err)
	}
	_, err = srv.Deposit(2, 20_00000)
	if err != nil {
		log.Print(err)
	}
	
}
