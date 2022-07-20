package psql

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_psql_GetConnection(t *testing.T) {
	log := logrus.New()
	p, err := NewPsqlStore(log)
	if err != nil {
		t.Fatalf("Error setting up db")
	}
	connection := p.GetConnection()
	assert.NotNil(t, connection)
}
