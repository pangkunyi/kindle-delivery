package mail

import (
	"net/smtp"
	"io/ioutil"
	message "github.com/sloonz/go-mime-message"
	"bytes"
	"os"
	"strings"
	"fmt"
)

type Account struct{
	from string
	to string
	username string
	password string
	smtpHost string
	port string 
}

func (this *Account) load(){
	data, err := ioutil.ReadFile(os.Getenv("HOME")+"/.kindle-delivery/account.conf")
	if err !=nil{
		panic(err)
	}
	fields :=strings.Fields(string(data))
	this.from=fields[0]
	this.to=fields[1]
	this.username=fields[2]
	this.password=fields[3]
	this.smtpHost=fields[4]
	this.port=fields[5]
	fmt.Printf("%v", this)
}

var acc Account
func init(){
	acc.load()
}
func Send(attachment []byte, filename string) error{
	from := acc.from
	to := []string{acc.to}
	subject :="convert file "+filename

	msg :=message.NewMultipartMessage("mixed","")
	att :=message.NewBinaryMessage(bytes.NewBuffer(attachment))
	att.SetHeader("Content-Type","application/octet-stream; name="+filename)
	att.SetHeader("Content-Disposition","attachment; filename="+filename)
	msg.AddPart(att)
	msg.SetHeader("From",from)
	msg.SetHeader("Subject",subject)
	msg.SetHeader("To", to[0])
	body,err := ioutil.ReadAll(msg)
	if err!=nil {
		return err
	}
//	fmt.Println(string(body))

	auth :=smtp.PlainAuth("", acc.username, acc.password, acc.smtpHost)
	err = smtp.SendMail(acc.smtpHost+":"+acc.port, auth, from, to, body)
	return err
}

