package pod

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	podAPI "github.com/47Cid/Castle/api"
	"github.com/47Cid/Castle/config"
	"github.com/47Cid/Castle/logger"
	"github.com/47Cid/Castle/message"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Pod struct {
	container     types.Container
	isBusy        bool
	podType       string
	weight        int
	containerIP   string
	containerPort string
}

var pods []Pod
var currentIndex = 0

func GetPods(dockerClient client.Client) {
	// List containers
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// Clear the global pods slice
	pods = []Pod{}

	// Fill up the global Pod variable
	for _, container := range containers {
		// Inspect container details
		inspect, err := dockerClient.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			panic(err)
		}
		// Get the labels
		labels := inspect.Config.Labels
		pType := labels["type"]

		weightStr, ok := labels["weight"]
		if !ok {
			logger.WAFLog.Errorf("weight label not found")
			continue
		}

		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			logger.WAFLog.Errorf("error converting weight to integer: %v", err)
			continue
		}

		// Get the IP address of the Docker container
		containerIP := inspect.NetworkSettings.IPAddress
		// Get the port mapping of the Docker container
		containerPort := ""
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				containerPort = fmt.Sprintf("%d:%d", port.PrivatePort, port.PublicPort)
				break
			}
		}

		pod := Pod{
			container:     container,
			isBusy:        false,
			podType:       pType,
			weight:        weight,
			containerIP:   containerIP,
			containerPort: containerPort,
		}

		pods = append(pods, pod)
	}
}

func logPods(dockerClient client.Client) {
	for _, pod := range pods {
		inspect, err := dockerClient.ContainerInspect(context.Background(), pod.container.ID)
		if err != nil {
			panic(err)
		}
		logger.WAFLog.Infof("Container ID: %s\n", pod.container.ID)
		logger.WAFLog.Infof("Image: %s\n", pod.container.Image)
		logger.WAFLog.Infof("Command: %s\n", pod.container.Command)
		logger.WAFLog.Infof("Status: %s\n", pod.container.Status)
		logger.WAFLog.Infof("IsBusy: %v\n", pod.isBusy)
		logger.WAFLog.Infof("PodType: %s\n", pod.podType)
		logger.WAFLog.Infof("Weight: %d\n", pod.weight)
		logger.WAFLog.Infof("IPAddress: %s\n", inspect.NetworkSettings.IPAddress) // Print the IP address
		// Print the port numbers
		for _, portBindings := range inspect.NetworkSettings.Ports {
			if len(portBindings) > 0 {
				logger.WAFLog.Infof("Port: %s\n", portBindings[0].HostPort)
			}
		}

	}
}

func VerifyMessage(message message.Message) bool {
	// Get the label for the message's destination
	label := config.GetLabel(message.Destination)

	// Pick a pod that's not busy using a weighted round robin algorithm
	for i := 0; i < len(pods); i++ {
		pod := pods[currentIndex]
		if !pod.isBusy && pod.weight > 0 && pod.podType == label {
			// Mark the pod as busy
			pod.isBusy = true

			// Process the message
			isValid := processMessage(pod, message)

			// Mark the pod as not busy
			pod.isBusy = false

			// Update the current index for the next round
			currentIndex = (currentIndex + 1) % len(pods)
			return isValid
		}

		// If the current pod is busy or has a weight of 0, move to the next pod
		currentIndex = (currentIndex + 1) % len(pods)
	}
	logger.WAFLog.Error("No pods available")
	// If no pods are available, return false
	return false
}

func processMessage(pod Pod, message message.Message) bool {
	logger.WAFLog.Infof("Processing message %+v using pod %+v", message.Destination, pod.podType)

	// TODO Make this a post request that sends the client message to the pod
	// Create the URL for the /verify endpoint
	verifyURL := fmt.Sprintf("http://%s:%s/verify", "localhost", "3032")

	// Send an HTTP GET request to the server running inside the pod
	resp, err := http.Get(verifyURL)
	if err != nil {
		logger.WAFLog.Errorf("Error sending GET request to %s: %v", verifyURL, err)
		return false
	}
	defer resp.Body.Close()
	logger.WAFLog.Info("GET request sent, reading response body")

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WAFLog.Errorf("Error reading response body: %v", err)
		return false
	}
	logger.WAFLog.Info("Response body read, logging response")

	// Log the response body
	logger.WAFLog.Infof("Received response: %s", string(body))

	var podResp podAPI.Response
	err = json.Unmarshal(body, &podResp)
	if err != nil {
		logger.WAFLog.Errorf("Error parsing JSON response: %v", err)
		return false
	}
	logger.WAFLog.Infof("Raw response from pod: %+v", podResp.Valid)

	isValid, err := strconv.ParseBool(podResp.Valid)
	if err != nil {
		logger.WAFLog.Errorf("Error converting string to bool: %v", err)
		return false
	}
	logger.WAFLog.Infof("Response from pod: %v", isValid)

	return isValid
}

func Init() {
	// Initialize the logger
	logger.InitWAFProxy()

	logger.WAFLog.Info("Initializing pods")
	// Create a docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	GetPods(*dockerClient)
	logPods(*dockerClient)
}
