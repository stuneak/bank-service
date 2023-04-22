package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stuneak/simplebank/db/sqlc_custom"
	"github.com/stuneak/simplebank/util"
)

func newTestServer(t *testing.T, store sqlc_custom.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)

	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
