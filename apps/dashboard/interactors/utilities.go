package interactors

import (
	"golang.org/x/crypto/bcrypt"
	. "iugo.fleet/common/logger"
	"github.com/robfig/cron"
	"fmt"
	"stock/common/projectArch/interactors"
	"time"
)


func comparePasswords(hashedPwd []byte, plainPwd []byte) bool {

	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func hashAndSalt(pwd []byte) []byte {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		LogError(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return hash
}

func StartReceivingCheckCronJob(){
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		fmt.Println("StartReceivingCheckCronJob ...")
		receivings,err := interactors.ReceivingRepo.SelectReceivings("","","","",0,0)
		if err != nil {
			LogError(err)
			return
		}
		timeNow := int(time.Now().Unix())
		for _,item := range receivings.Items{
			if item.ExpectedDate < timeNow{
				interactors.ReceivingRepo.SetStatus("Gecikmiş",item.Id)
			}
		}

	})
	c.Start()
}