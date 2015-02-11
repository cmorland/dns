package domains

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tlds"
)

const (
	SELECT string = "SELECT * FROM domain "
)

type Result struct {
	One  func() (Domain, error)
	List func() ([]Domain, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Domain, error) {
			return Get(query, args...)
		},
		func() ([]Domain, error) {
			return GetList(query, args...)
		},
	}
}

func GetByJoinWhoisEmails(email string) Result {
	return newResult(
		"SELECT DISTINCT d.* FROM domain AS d JOIN whois w ON d.uuid = w.domain WHERE w.emails ? $1",
		email,
	)
}

func GetByNameAndTLD(name string, tld int32) Result {
	return newResult(
		SELECT+"WHERE name = $1 AND tld = $2",
		name, tld,
	)
}

func GetByUUID(uuid string) Result {
	return newResult(
		SELECT+"WHERE uuid = $1",
		uuid,
	)
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetAllLimitOffset(limit, offset int) Result {
	return newResult(
		SELECT+"LIMIT $1 OFFSET $2",
		limit, offset,
	)
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var dUUID string
	var tldID int32
	err := row.Scan(&dUUID, &d.Name, &tldID)
	if err != nil {
		return d, err
	}
	d.UUID = uuid.Parse(dUUID)
	d.TLD, err = tlds.GetByID(tldID).One()
	return d, err
}

func Get(query string, args ...interface{}) (Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return Domain{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Domain{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Domain{}, err
	}
	var list []Domain
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}