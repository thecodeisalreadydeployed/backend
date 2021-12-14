package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("CODEDEPLOY")
	viper.BindEnv("API_URL")
	viper.BindEnv("DEPLOYMENT_ID")
	viper.BindEnv("DEPLOYMENT_GIT_SOURCE")
	viper.BindEnv("DEPLOYMENT_BUILD_CONFIGURATION")

	deploymentID := viper.GetString("CODEDEPLOY_DEPLOYMENT_ID")
	apiURL := fmt.Sprintf("%s/internal/%s/events", strings.TrimSuffix(viper.GetString("CODEDEPLOY_API_URL"), "/"), deploymentID)

	scanner := bufio.NewScanner(os.Stdin)
	queue := NewQueue()
	isExporting := false
	done := make(chan bool, 1)

	for scanner.Scan() {
		text := scanner.Text()
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

		_ = queue.Dequeue()
	}
}