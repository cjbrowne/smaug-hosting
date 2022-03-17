package containers

import (
	"bitbucket.org/smaug-hosting/services/database"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type ContainerRepository struct{}

const tableName = "containers"

func (cr ContainerRepository) FindById(id int64) (*Container, error) {
	sql, params, err := squirrel.Select("name", "tier", "software", "user_id", "id").From(tableName).Where("id = ?", id).ToSql()
	if err != nil {
		return nil, err
	}

	row := database.Connection.QueryRowx(sql, params...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	c := new(Container)

	err = row.StructScan(c)

	return c, err
}

func (cr ContainerRepository) Save(container Container) (Container, error) {
	var result Container // only used if we fail

	containerMap := map[string]interface{}{
		"name":     container.Name,
		"tier":     container.Tier,
		"software": container.Software,
		"user_id":  container.UserId,
	}

	if container.Id > 0 {
		containerMap["id"] = container.Id
	}

	sql, params, err := squirrel.
		Insert(tableName).
		SetMap(containerMap).
		ToSql()

	if err != nil {
		return result, err
	}

	res, err := database.Connection.Exec(sql, params...)
	if err != nil {
		return result, err
	}

	container.Id, err = res.LastInsertId()

	return container, err
}

func (cr ContainerRepository) GetContainersForUser(userId int64) ([]Container, error) {
	sql, params, err := squirrel.
		Select("name", "tier", "software", "id", "user_id").
		From(tableName).
		Where("user_id=?", userId).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := database.Connection.Queryx(sql, params...)
	if err != nil {
		return nil, err
	}

	containers := make([]Container, 0)

	for rows.Next() {
		container := new(Container)
		err = rows.StructScan(container)
		if err != nil {
			logrus.Errorf("Could not scan row into container: %s", err)
			continue
		}
		containers = append(containers, *container)
	}

	return containers, err
}
