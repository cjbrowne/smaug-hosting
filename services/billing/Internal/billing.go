package billingint

import (
	"bitbucket.org/smaug-hosting/services/billing/pricing"
	"bitbucket.org/smaug-hosting/services/container-service/containers"
	"bitbucket.org/smaug-hosting/services/idp/users"
	"github.com/sirupsen/logrus"
)

func BillAllUsers() {
	// basically everything in this file will be logged with severity CRITICAL
	criticalLogger := logrus.WithField("severity", "CRITICAL")

	allUsers, err := users.UserRepository{}.FindAll()
	if err != nil {
		criticalLogger.Errorf("Could not fetch users for billing: %s", err)
		return
	}

	for _, user := range allUsers {
		go billUser(user)
	}

}

func billUser(user users.User) {
	// basically everything in this file will be logged with severity CRITICAL
	criticalLogger := logrus.WithField("severity", "CRITICAL")

	allContainers, err := containers.ContainerRepository{}.GetContainersForUser(user.Id)
	if err != nil {
		criticalLogger.Errorf("Could not fetch containers for user %d for billing purposes: %s", user.Id, err)
		return
	}

	for _, container := range allContainers {
		containerStatus, err := containers.GetStatusForContainer(container)
		if err != nil {
			criticalLogger.Errorf("Could not get status for container %d: %s", container.Id, err)
			continue
		}
		if containerStatus.Up {
			price, err := pricing.PricingRepository{}.FindPriceBySoftwareAndTier(container.Software, container.Tier)
			if err != nil {
				criticalLogger.Errorf("Could not find price for software %s / tier %d: %s", container.Software, container.Tier, err)
				continue
			}
			user.Balance -= price.Amount
			if user.Balance <= 0 {
				go func() {
					allContainers, err := containers.ContainerRepository{}.GetContainersForUser(user.Id)
					if err != nil {
						logrus.Errorf("Could not shut down containers for out-of-balance user %d: %s", user.Id, err)
					}
					for _, container := range allContainers {
						err := containers.StopContainer(container)
						if err != nil {
							logrus.Errorf("Could not shut down container %d for out-of-balance user %d: %s", container.Id, user.Id, err)
							continue
						}
					}
				}()
			}
		}
	}

	// save the modified user back to the repo after billing
	err = users.UserRepository{}.Save(user)
	if err != nil {
		criticalLogger.Errorf("Could not save billed user back to the user repository: %s", err)
	}
}
