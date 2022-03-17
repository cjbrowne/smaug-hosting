package containers

import (
	"bitbucket.org/smaug-hosting/services/cache"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

func getImageForSoftware(software string) string {
	switch software {
	case "minecraft":
		// todo: we will eventually be using our own minecraft server images kept on a private registry
		// but for the POC/MVP we can just pull in someone else's
		return "itzg/minecraft-server:20190824"
	default:
		logrus.Errorf("Invalid software (docker service create is about to break...): %s", software)
		return ""
	}
}

func removeContainer(container Container) {
	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return
	}

	err = dockerClient.ServiceRemove(context.Background(), getServiceIdForContainer(container))
	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not remove container, docker responded with an error: %s", err)
		return
	}
}

func getPortForSoftware(software string) uint32 {
	switch software {
	case "minecraft":
		return 25565
	}
	// fail
	logrus.Errorf("Unrecognised software: %s", software)
	return 0
}

func getDataDirForSoftware(software string) string {
	switch software {
	case "minecraft":
		return "/data"
	}
	// fail
	logrus.Errorf("Unrecognised software: %s", software)
	return ""
}

func spinUpContainer(c Container) {
	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return
	}

	srvcCreateResponse, err := dockerClient.ServiceCreate(context.Background(), swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: getServiceIdForContainer(c),
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: getImageForSoftware(c.Software),
				Env:   getEnvForSoftware(c.Software),
				Mounts: []mount.Mount{
					mount.Mount{
						Type:   "volume",
						Source: getServiceIdForContainer(c),
						Target: getDataDirForSoftware(c.Software), // todo: multiple mounts?
					},
				},
			},
		},
		Mode:         swarm.ServiceMode{},
		UpdateConfig: nil,
		Networks:     nil,
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				{
					Name:          "",
					Protocol:      "tcp",
					TargetPort:    getPortForSoftware(c.Software),
					PublishedPort: getUnusedPortForContainer(c),
					PublishMode:   "",
				},
			},
		},
	}, types.ServiceCreateOptions{})

	if err != nil {
		c.LastError = "could not create docker service"
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker service for customer: %s", err)
		return
	}

	for _, warning := range srvcCreateResponse.Warnings {
		logrus.Warnf("Warning while creating docker service: %s", warning)
	}
}

func getUnusedPortForContainer(c Container) uint32 {
	port, err := cache.Client.Get("ports." + getServiceIdForContainer(c)).Result()
	if err == redis.Nil {
		nextPort, err := cache.Client.Get("next_port").Result()
		if err == redis.Nil {
			cache.Client.Set("next_port", "50000", 0)
			nextPort = "50000"
		}
		port = nextPort
	} else if err != nil {
		logrus.Errorf("Could not get free port from redis, returning zero so client gets something at least (err=%s)", err)
		return 0
	}
	res, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		logrus.Errorf("Invalid port: %s (%s)", port, err)
		return 0
	}
	return uint32(res)
}

func getEnvForSoftware(software string) []string {
	switch software {
	case "minecraft":
		return []string{"EULA=TRUE"}
	default:
		return []string{}
	}
}

func getServiceIdForContainer(c Container) string {
	return fmt.Sprintf("whelp-%s-%d-%d-%d", c.Software, c.UserId, c.Tier, c.Id)
}

func GetStatusForContainer(container Container) (ContainerStatus, error) {
	var containerStatus ContainerStatus

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return containerStatus, err
	}

	args, err := filters.ParseFlag(fmt.Sprintf("service=%s", getServiceIdForContainer(container)), filters.NewArgs())
	if err != nil {
		logrus.Errorf("Could not parse args: %s", err)
		return containerStatus, err
	}

	tasks, err := dockerClient.TaskList(context.Background(), types.TaskListOptions{
		Filters: args,
	})

	if err != nil {
		// HACK: the only way to know if the error was "not found"
		if strings.Contains(err.Error(), "not found") {
			// try to spin up the container again, we probably are recovering from some kind of crash
			logrus.Warnf("Container %s (id=%d) not found, trying to spin it up again", container.Name, container.Id)
			spinUpContainer(container)
			return containerStatus, err
		} else {
			logrus.Errorf("Could not get service info for container %s: %s", getServiceIdForContainer(container), err)
			return containerStatus, err
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})

	if len(tasks) == 0 {
		containerStatus.Up = false
		containerStatus.State = "stopped"
	} else {
		task := tasks[0]
		containerStatus.State = string(task.Status.State)
		if task.Status.State == swarm.TaskStateRunning {
			containerStatus.Up = true
		} else if task.Status.State == swarm.TaskStateShutdown {
			containerStatus.Up = false
			containerStatus.State = "stopped"
		} else if task.Status.State != swarm.TaskStateFailed {
			// break on non-failed status because we only want to report "failed" if _all_ tasks failed.
			containerStatus.Up = false
		}
	}

	return containerStatus, nil
}

func StopContainer(c Container) error {

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return err
	}

	service, _, err := dockerClient.ServiceInspectWithRaw(context.Background(), getServiceIdForContainer(c))

	if err != nil {
		logrus.Errorf("Could not inspect service: %s", err)
		return err
	}

	zero := uint64(0)

	serviceUpdateResponse, err := dockerClient.ServiceUpdate(
		context.Background(),
		getServiceIdForContainer(c),
		swarm.Version{
			Index: service.Version.Index,
		},
		swarm.ServiceSpec{
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: getImageForSoftware(c.Software),
					Env:   getEnvForSoftware(c.Software),
					Mounts: []mount.Mount{
						mount.Mount{
							Type:   "volume",
							Source: getServiceIdForContainer(c),
							Target: getDataDirForSoftware(c.Software), // todo: multiple mounts?
						},
					},
				},
			},
			Mode: swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{
					Replicas: &zero,
				},
			},
			Annotations: swarm.Annotations{
				Name: getServiceIdForContainer(c),
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: []swarm.PortConfig{
					{
						Name:          "",
						Protocol:      "tcp",
						TargetPort:    getPortForSoftware(c.Software),
						PublishedPort: getUnusedPortForContainer(c),
						PublishMode:   "",
					},
				},
			},
		},
		types.ServiceUpdateOptions{},
	)
	if err != nil {
		logrus.Errorf("Could not update service: %s", err)
		return err
	}

	for _, warning := range serviceUpdateResponse.Warnings {
		logrus.Warnf("Service update warning: %s", warning)
	}

	return nil
}

func startContainer(c Container) error {

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return err
	}

	service, _, err := dockerClient.ServiceInspectWithRaw(context.Background(), getServiceIdForContainer(c))

	if err != nil {
		logrus.Errorf("Could not inspect service: %s", err)
		return err
	}

	one := uint64(1)

	serviceUpdateResponse, err := dockerClient.ServiceUpdate(
		context.Background(),
		getServiceIdForContainer(c),
		swarm.Version{
			Index: service.Version.Index,
		},
		swarm.ServiceSpec{
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: getImageForSoftware(c.Software),
					Env:   getEnvForSoftware(c.Software),
					Mounts: []mount.Mount{
						{
							Type:   "volume",
							Source: getServiceIdForContainer(c),
							Target: getDataDirForSoftware(c.Software), // todo: multiple mounts?
						},
					},
				},
			},
			Mode: swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{
					Replicas: &one,
				},
			},
			Annotations: swarm.Annotations{
				Name: getServiceIdForContainer(c),
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: []swarm.PortConfig{
					{
						Name:          "",
						Protocol:      "tcp",
						TargetPort:    getPortForSoftware(c.Software),
						PublishedPort: getUnusedPortForContainer(c),
						PublishMode:   "",
					},
				},
			},
		},
		types.ServiceUpdateOptions{},
	)
	if err != nil {
		logrus.Errorf("Could not update service: %s", err)
		return err
	}

	for _, warning := range serviceUpdateResponse.Warnings {
		logrus.Warnf("Service update warning: %s", warning)
	}

	return nil
}

func GetIpAndPortForContainer(container Container) (string, uint32, error) {
	var ip string
	var port uint32

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not create docker client: %s", err)
		return ip, port, err
	}

	srvc, _, err := dockerClient.ServiceInspectWithRaw(context.Background(), getServiceIdForContainer(container))

	if err != nil {
		logrus.WithField("severity", "CRITICAL").Errorf("Could not get connection information for service: %s", err)
		return ip, port, err
	}

	port = srvc.Endpoint.Ports[0].PublishedPort
	info, _ := dockerClient.Info(context.Background())
	ip = info.Swarm.NodeAddr

	return ip, port, err
}
