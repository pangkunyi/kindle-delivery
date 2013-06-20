package mail

import (
	"net/smtp"
	"io/ioutil"
	"time"
	message "github.com/sloonz/go-mime-message"
	"bytes"
	"fmt"
)

func Send(attachment []byte, filename string, suffix string) error{
	from := "connect0829@qq.com"
	to := []string{"pangkunyi_72@kindle.com"}
	timeString :=time.Now().Format("-2006-01-02")
	filename=filename+timeString+suffix
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

	auth :=smtp.PlainAuth("", "connect0829@qq.com", "password", "smtp.qq.com")
	err = smtp.SendMail("smtp.qq.com:25", auth, from, to, body)
	return err
}

