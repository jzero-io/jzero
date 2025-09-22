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
	Database = "database"
	User     = "user"
	Password = "password"
	Host     = "host"
	Port     = "port"
)

// ParseDSN parses various supported DSN types into a map of key/value pairs which can be used as valid tags.
func ParseDSN(driverName, dsn string) (meta map[string]string, err error) {
	meta = make(map[string]string)
	switch driverName {
	case "mysql":
		meta, err = parseMySQLDSN(dsn)
		if err != nil {
			return meta, err
		}
	case "postgres", "pgx":
		meta, err = parsePostgresDSN(dsn)
		if err != nil {
			return meta, err
		}
	default:
		// Try to parse the DSN and see if the scheme contains a known driver name.
		u, e := url.Parse(dsn)
		if e != nil {
			// dsn is not a valid URL, so just ignore
			return meta, err
		}
		if driverName != u.Scheme {
			// In some cases the driver is registered under a non-official name.
			// For example, "Test" may be the registered name with a DSN of "postgres://postgres:postgres@127.0.0.1:5432/fakepreparedb"
			// for the purposes of testing/mocking.
			// In these cases, we try to parse the DSN based upon the DSN itself, instead of the registered driver name
			return ParseDSN(u.Scheme, dsn)
		}
	}
	return meta, nil
}

// parseMySQLDSN parses a mysql-type dsn into a map.
func parseMySQLDSN(dsn string) (m map[string]string, err error) {
	var cfg *mySQLConfig
	if cfg, err = mySQLConfigFromDSN(dsn); err == nil {
		host, port, _ := net.SplitHostPort(cfg.Addr)
		m = map[string]string{
			"user":     cfg.User,
			"host":     host,
			"port":     port,
			"database": cfg.DBName,
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
