package monitor

import (
	"log"

	"github.com/gocql/gocql"
)

type scyllaConfig struct {
	ConnectionStrings []string
	Keyspace          string
	DefaultTTL        int
	Username          string
	Password          string
}

func InitScylla(cfg scyllaConfig) *gocql.Session {
	cluster := gocql.NewCluster(cfg.ConnectionStrings...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	cluster.Keyspace = cfg.Keyspace
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.LocalOne
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Cannot connect to scylla: ", err)
	}
	return session
}
