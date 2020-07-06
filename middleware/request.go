package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const HeaderRequestId = "X-Request-Id"

var (
	prefix string
	reqId  uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func RequestStart(c *gin.Context) {
	requestId := c.Request.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = genRequestId()
	}
	c.Set(HeaderRequestId, requestId)
	c.Header(HeaderRequestId, requestId)

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// 读取后写回
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	entry := logrus.WithFields(logrus.Fields{
		"url":       c.Request.RequestURI,
		"params":    string(bodyBytes),
		"header":    c.Request.Header,
		"clientIp":  c.ClientIP(),
		"requestId": requestId,
	})

	entry.Info("request_in")

	//fmt.Printf("%v url=%s||method=%s||param=%v \n", time.Now(), c.Request.URL, c.Request.Method, c.Request.Form)
}

func genRequestId() string {
	id := atomic.AddUint64(&reqId, 1)
	return fmt.Sprintf("%s-%06d", prefix, id)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestOut(c *gin.Context) {
	startTime := time.Now()

	// replace c.Write to hook response
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	// process request
	c.Next()

	// log out
	latency := time.Since(startTime)
	requestId, _ := c.Get(HeaderRequestId)
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// 读取后写回
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	entry := logrus.WithFields(logrus.Fields{
		"url":       c.Request.RequestURI,
		"params":    string(bodyBytes),
		"header":    c.Request.Header,
		"clientIp":  c.ClientIP(),
		"requestId": requestId,
		"httpCode":  c.Writer.Status(),
		"response":  blw.body.String(),
		"startTime": startTime.String(),
		"latency":   latency.String(),
	})

	entry.Info("request_out")

	//fmt.Printf("%v url=%s||method=%s||resp=%v latency= %d \n", time.Now(), c.Request.URL, c.Request.Method, "", latency)
}
