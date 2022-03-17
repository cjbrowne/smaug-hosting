package containers

import (
	"bitbucket.org/smaug-hosting/services/billing/pricing"
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/idp/users"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func HandleGetContainers(response http.ResponseWriter, request *http.Request) {
	claims := request.Context().Value("token_claims").(tokens.TokenClaims)

	containers, err := ContainerRepository{}.GetContainersForUser(claims.UserId)
	if err != nil {
		logrus.Errorf("Could not get containers for user: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not fetch containers for user", response)
		return
	}

	for i, c := range containers {
		c.Status, err = GetStatusForContainer(c)
		if err != nil {
			// hard-code status to be predictable if any error occurred while fetching the status
			c.Status.Up = false
			c.Status.State = "unknown"
		}
		if c.Status.Up {
			c.IP, c.Port, err = GetIpAndPortForContainer(c)
		}
		containers[i] = c
	}

	libhttp.SendJson(containers, response)
}

func HandleStopContainer(response http.ResponseWriter, request *http.Request) {
	claims := request.Context().Value("token_claims").(tokens.TokenClaims)
	containerId := request.Context().Value("containerId").(string)

	sql, params, err := squirrel.Select("user_id", "id", "software", "tier").From(tableName).Where("id = ?", containerId).ToSql()

	if err != nil {
		logrus.Debugf("Could not fetch container from db: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	row := database.Connection.QueryRowx(sql, params...)
	if row.Err() != nil {
		logrus.Errorf("Could not fetch container from db: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	container := Container{}

	err = row.StructScan(&container)
	if err != nil {
		logrus.Errorf("Could not scan container into struct: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	if container.UserId != claims.UserId {
		logrus.Warnf("Tried to delete a container that doesn't belong to him: %d (target container=%d)", claims.UserId, containerId)
		libhttp.SendError(http.StatusUnauthorized, "You can only stop your own containers", response)
		return
	}

	err = StopContainer(container)
	if err != nil {
		logrus.Errorf("Could not stop container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	// empty 200 response if all went well
	libhttp.SendJson(struct{}{}, response)
}

type CreateContainerRequest struct {
	Name     string
	Software string
	Tier     int
}

func HandlePostContainer(response http.ResponseWriter, request *http.Request) {
	body := CreateContainerRequest{}

	err := ParseBody(request.Body, &body)
	if err != nil {
		logrus.Debugf("Could not parse body: %s", err)
		libhttp.SendError(http.StatusBadRequest, "Could not parse container request", response)
		return
	}

	claims := request.Context().Value("token_claims").(tokens.TokenClaims)

	price, err := pricing.PricingRepository{}.FindPriceBySoftwareAndTier(body.Software, body.Tier)
	if err != nil {
		logrus.Errorf("Could not fetch price from database: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not fetch price for new container", response)
		return
	}

	user, err := users.UserRepository{}.Find(claims.UserId)
	if err != nil {
		logrus.Errorf("Could not fetch user: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not fetch user details for checking balance", response)
		return
	}

	if !user.Verified {
		libhttp.SendError(http.StatusForbidden, "You must verify your email address before you can spin up whelps", response)
		return
	}

	if user.Balance <= price.Amount {
		libhttp.SendError(http.StatusPaymentRequired, "You do not have sufficient funds to complete this request", response)
		return
	}

	container, err := ContainerRepository{}.Save(Container{
		Name:     body.Name,
		Tier:     body.Tier,
		Software: body.Software,
		UserId:   claims.UserId,
	})

	// asynchronously start spinning up the container to bring it live
	go spinUpContainer(container)

	if err != nil {
		logrus.Errorf("Could not save new container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not save new container to database", response)
		return
	}

	libhttp.SendJsonWithStatus(http.StatusCreated, container, response)
}

func ParseBody(body io.ReadCloser, target *CreateContainerRequest) error {
	res, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, target)
}

func HandleStartContainer(response http.ResponseWriter, request *http.Request) {
	claims := request.Context().Value("token_claims").(tokens.TokenClaims)
	containerId := request.Context().Value("containerId").(string)

	sql, params, err := squirrel.Select("user_id", "id", "software", "tier").From(tableName).Where("id = ?", containerId).ToSql()

	if err != nil {
		logrus.Debugf("Could not fetch container from db: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	row := database.Connection.QueryRowx(sql, params...)
	if row.Err() != nil {
		logrus.Errorf("Could not fetch container from db: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	container := Container{}

	err = row.StructScan(&container)
	if err != nil {
		logrus.Errorf("Could not scan container into struct: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	if container.UserId != claims.UserId {
		logrus.Warnf("Tried to delete a container that doesn't belong to him: %d (target container=%d)", claims.UserId, containerId)
		libhttp.SendError(http.StatusUnauthorized, "You can only stop your own containers", response)
		return
	}

	err = startContainer(container)
	if err != nil {
		logrus.Errorf("Could not stop container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not stop container", response)
		return
	}

	// empty 200 response if all went well
	libhttp.SendJson(struct{}{}, response)
}

func HandleDeleteContainer(response http.ResponseWriter, request *http.Request) {
	claims := request.Context().Value("token_claims").(tokens.TokenClaims)
	containerId := request.Context().Value("containerId").(string)

	containerIdInt64, err := strconv.ParseInt(containerId, 10, 64)
	if err != nil {
		logrus.Errorf("Could not parse container id: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "could not delete container, invalid id", response)
		return
	}

	container, err := ContainerRepository{}.FindById(containerIdInt64)
	if err != nil {
		logrus.Errorf("Could not fetch whelp for deletion: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "could not fetch whelp for deletion", response)
		return
	}

	status, err := GetStatusForContainer(*container)

	if err != nil {
		logrus.Errorf("Could not determine status of container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "could not determine if whelp is ready for deletion", response)
		return
	}

	if status.Up {
		libhttp.SendError(http.StatusInternalServerError, "whelps must be stopped before they can be removed", response)
		return
	}

	sql, params, err := squirrel.Delete(tableName).Where("user_id = ? AND id = ?", claims.UserId, containerId).ToSql()
	if err != nil {
		logrus.Errorf("Could not delete container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "could not delete container", response)
		return
	}

	result, err := database.Connection.Exec(sql, params...)
	if err != nil {
		logrus.Errorf("Could not delete container: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "could not delete container", response)
		return
	}

	if rows, err := result.RowsAffected(); rows < 1 {
		if err != nil {
			logrus.Errorf("Could not verify container was deleted: %s", err)
			libhttp.SendError(http.StatusInternalServerError, "could not verify container was deleted (please try again)", response)
			return
		}
		logrus.Warnf("Duplicate delete on container: %d", containerId)
	}

	go removeContainer(*container)

	// empty 200 response if everything went well
	libhttp.SendJson(struct{}{}, response)
}
