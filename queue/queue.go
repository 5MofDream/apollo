package queue

import (
	"github.com/kr/beanstalk"
	"issueCenter/apollo/config"
)

func init()  {

}

func New() *beanstalk.Conn{

	connector :=  config.Read("queue.driver")
	var con *beanstalk.Conn
	switch connector {
	case "beanstalk":
		connect, err := beanstalk.Dial("tcp", config.Read("queue.host"))
		if err != nil {
			panic(err)
		}
		con = connect
	default:

	}
	return con
}
