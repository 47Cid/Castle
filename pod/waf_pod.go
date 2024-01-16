package pod

import (
	"context"
	"strconv"

	"github.com/47Cid/Castle/config"
	"github.com/47Cid/Castle/logger"
	"github.com/47Cid/Castle/message"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Pod struct {
	container types.Container
	isBusy    bool
	podType   string
	weight    int
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

		pod := Pod{
			container: container,
			isBusy:    false,
			podType:   pType,
			weight:    weight,
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
	return true
}

func processMessage(pod Pod, message message.Message) bool {
	// TODO: Implement this function to process the message using the given pod
	return true

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
