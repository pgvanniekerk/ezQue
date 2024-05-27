package ezQue

import (
	"context"
	"github.com/pgvanniekerk/ezQue/oraaq"
	"os"
	"strconv"
	"testing"
)

// Oracle
func TestOracleConnect(t *testing.T) {

	server := os.Getenv("DB_SERVER")
	portString := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		t.Fatalf("Invalid port number: %v", err)
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	sid := os.Getenv("DB_SID")

	q, err := Connect(oraaq.OracleAqJms,
		oraaq.Queue("text_msg_queue",
			oraaq.LocatedAt(server, uint16(port)),
			oraaq.AuthenticatedWith(username, password),
			oraaq.UsingSID(sid),
		),
	)
	if err != nil {
		t.Error(err)
	}
	err = q.Disconnect(context.Background())
	if err != nil {
		t.Error(err)
	}

	//q, err = Connect(oraaq.OracleAqJms,
	//	oraaq.Queue("text_msg_queue",
	//		oraaq.UsingJdbcString(fmt.Sprintf("jdbc:oracle:thin:@%s:%d:%s", server, port, sid)),
	//		oraaq.AuthenticatedWith(username, password),
	//	),
	//)
	//if err != nil {
	//	t.Fatalf(err.Error())
	//}
	//err = q.Disconnect(context.Background())
	//if err != nil {
	//	t.Fatalf(err.Error())
	//}

}
