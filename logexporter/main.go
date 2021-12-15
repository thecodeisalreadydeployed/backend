package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viper.SetEnvPrefix("CODEDEPLOY")
	viper.BindEnv("API_URL")
	viper.BindEnv("DEPLOYMENT_ID")
	viper.BindEnv("DEPLOYMENT_GIT_SOURCE")
	viper.BindEnv("DEPLOYMENT_BUILD_CONFIGURATION")

	deploymentID := viper.GetString("CODEDEPLOY_DEPLOYMENT_ID")
	apiURL := fmt.Sprintf("%s/%s/events", strings.TrimSuffix(viper.GetString("CODEDEPLOY_API_URL"), "/"), deploymentID)

	scanner := bufio.NewScanner(os.Stdin)
	queue := NewQueue()
	isExporting := false
	done := make(chan bool, 1)

	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("In: %s\n", text)
		queue.Enqueue(text)

		if !isExporting {
			isExporting = true
			go export(apiURL, queue, done)
		}
	}

	queue.End()
	<-done
}

func export(apiURL string, queue Queue, done chan bool) {
	for {
		if queue.N() == 0 && queue.IsEnded() {
			done <- true
			break
		}

		text := queue.Dequeue()
		data := map[string]string{
			"text":       text,
			"exportedAt": time.Now().String(),
			"type":       "DEBUG",
		}
		dataJSON, _ := json.Marshal(data)
		requestBody := bytes.NewBuffer(dataJSON)

		req, err := http.NewRequest("POST", apiURL, requestBody)
		req.Header.Add("X-CodeDeploy-Internal-Request", "True")

		if err != nil {
			zap.L().Error("failed to create an HTTP request")
			continue
		}

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			zap.L().Error("failed to send an HTTP request")
			continue
		}
	}
}
