package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IntegrationTestSuite struct {
	suite.Suite
	botReader *tgbotapi.BotAPI
	appCancel context.CancelFunc
	srvPort   int
	c         http.Client
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.Require().NoError(godotenv.Load())
	testToken := os.Getenv("HC_TEST_TOKEN")
	suite.Require().NotEmpty(testToken)

	bot, err := tgbotapi.NewBotAPI(testToken)
	suite.Require().NoError(err)
	bot.Debug = true

	suite.botReader = bot
	suite.c = http.Client{}
	suite.c.Timeout = 200 * time.Millisecond
	suite.setupEnv()
}

func (suite *IntegrationTestSuite) SetupTest() {
	a := New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	suite.appCancel = cancel
	go a.Run(ctx)
	suite.waitUntilServerReady()
}

func (suite *IntegrationTestSuite) TearDownTest() {
	suite.appCancel()
}

func (suite *IntegrationTestSuite) waitUntilServerReady() {
	var lastErr error
	var lastStatusCode int
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			{
				resp, err := suite.c.Get(fmt.Sprintf("http://localhost:%d/is_alive", suite.srvPort))
				lastErr = err
				if resp != nil {
					lastStatusCode = resp.StatusCode
				}

				if lastErr == nil && lastStatusCode == http.StatusOK {
					return
				}
			}
		case <-time.After(3 * time.Second):
			{
				var errStr string
				if lastErr != nil {
					errStr = lastErr.Error()
				}
				suite.Assert().FailNowf("waited too long for server to be ready err=%s status_code=%d", errStr, lastStatusCode)
			}
		}
	}
}

func (suite *IntegrationTestSuite) TestHealthy() {
	passphrase := os.Getenv("HC_PASSPHRASE")
	resp, err := suite.c.Post(fmt.Sprintf("http://localhost:%d/ping", suite.srvPort), "application/json", strings.NewReader(fmt.Sprintf("{\"passphrase\": \"%s\"}", passphrase)))
	suite.Require().NoError(err)
	suite.Assert().Equal(http.StatusOK, resp.StatusCode)

	actual := suite.awaitMessage()
	suite.Assert().Equal("State changed to 'Healthy'", actual)
}

func (suite *IntegrationTestSuite) setupEnv() {
	listener, err := net.Listen("tcp", ":0")
	suite.Require().NoError(err)
	suite.srvPort = listener.Addr().(*net.TCPAddr).Port
	suite.Require().NoError(listener.Close())
	os.Setenv("HC_PORT", strconv.Itoa(suite.srvPort))
	os.Setenv("HC_SAMPLE_RATE", "100ms")
	os.Setenv("HC_GRACE_PERIOD", "10s")
}

func (suite *IntegrationTestSuite) awaitMessage() string {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for {
		select {
		case update := <-suite.botReader.GetUpdatesChan(u):
			if update.Message != nil {
				logrus.WithField("message", update.Message).Info("received message")
				return update.Message.Text
			}
		case <-time.After(3 * time.Second):
			suite.FailNow("waited too long for message")
		}
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
