// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package dsn

import (
	"net"
	"net/url"
	"strings"
)

const (
	// DBApplication indicates the application using the database.
	DBApplication = "db.application"

	// DBName indicates the database name.
	DBName = "db.name"

	// DBUser indicates the user name of Database, e.g. "readonly_user" or "reporting_user".
	DBUser     = "db.user"
	TargetHost = "out.host"
	TargetPort = "out.port"
)

// ParseDSN parses various supported DSN types into a map of key/value pairs which can be used as valid tags.
func ParseDSN(driverName, dsn string) (meta map[string]string, err error) {
	meta = make(map[string]string)
	switch driverName {
	case "mysql":
		meta, err = parseMySQLDSN(dsn)
		if err != nil {
			return
		}
	case "postgres", "pgx":
		meta, err = parsePostgresDSN(dsn)
		if err != nil {
			return
		}
	default:
		// Try to parse the DSN and see if the scheme contains a known driver name.
		u, e := url.Parse(dsn)
		if e != nil {
			// dsn is not a valid URL, so just ignore
			return
		}
		if driverName != u.Scheme {
			// In some cases the driver is registered under a non-official name.
			// For example, "Test" may be the registered name with a DSN of "postgres://postgres:postgres@127.0.0.1:5432/fakepreparedb"
			// for the purposes of testing/mocking.
			// In these cases, we try to parse the DSN based upon the DSN itself, instead of the registered driver name
			return ParseDSN(u.Scheme, dsn)
		}
	}
	return reduceKeys(meta), nil
}

// reduceKeys takes a map containing parsed DSN information and returns a new
// map containing only the keys relevant as tracing tags, if any.
func reduceKeys(meta map[string]string) map[string]string {
	keysOfInterest := map[string]string{
		"user":             DBUser,
		"application_name": DBApplication,
		"dbname":           DBName,
		"host":             TargetHost,
		"port":             TargetPort,
	}
	m := make(map[string]string)
	for k, v := range meta {
		if nk, ok := keysOfInterest[k]; ok {
			m[nk] = v
		}
	}
	return m
}

// parseMySQLDSN parses a mysql-type dsn into a map.
func parseMySQLDSN(dsn string) (m map[string]string, err error) {
	var cfg *mySQLConfig
	if cfg, err = mySQLConfigFromDSN(dsn); err == nil {
		host, port, _ := net.SplitHostPort(cfg.Addr)
		m = map[string]string{
			"user":   cfg.User,
			"host":   host,
			"port":   port,
			"dbname": cfg.DBName,
		}
		return m, nil
	}
	return nil, err
}

// parsePostgresDSN parses a postgres-type dsn into a map.
func parsePostgresDSN(dsn string) (map[string]string, error) {
	var err error
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		// url form, convert to opts
		dsn, err = parseURL(dsn)
		if err != nil {
			return nil, err
		}
	}
	meta := make(map[string]string)
	if err := parseOpts(dsn, meta); err != nil {
		return nil, err
	}
	// remove sensitive information
	delete(meta, "password")
	return meta, nil
}
