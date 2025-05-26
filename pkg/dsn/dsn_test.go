package dsn

import (
	"fmt"
	"testing"
)

func TestParseDSN(t *testing.T) {
	// Add tests for ParseDSN function
	t.Run("Test ParseDSN", func(t *testing.T) {
		// Add test cases for ParseDSN function
		t.Run("Test ParseDSN with mysql", func(t *testing.T) {
			// Add test cases for ParseDSN function with mysql
			t.Run("Test ParseDSN with mysql and valid dsn", func(t *testing.T) {
				// Add test cases for ParseDSN function with mysql and valid dsn
				meta, err := ParseDSN("mysql", "user:password@tcp(localhost:3306)/dbname")
				if err != nil {
					t.Errorf("ParseDSN() error = %v", err)
				}
				if meta[DBUser] != "user" {
					t.Errorf("ParseDSN() user = %v, want %v", meta[DBUser], "user")
				}
				if meta[TargetHost] != "localhost" {
					t.Errorf("ParseDSN() host = %v, want %v", meta[TargetHost], "localhost")
				}
				if meta[TargetPort] != "3306" {
					t.Errorf("ParseDSN() port = %v, want %v", meta[TargetPort], "3306")
				}
				if meta[DBName] != "dbname" {
					t.Errorf("ParseDSN() dbname = %v, want %v", meta[DBName], "dbname")
				}
			})
		})
	})

	t.Run("Test ParseDSN with postgres", func(t *testing.T) {
		// Add test cases for ParseDSN function with postgres
		t.Run("Test ParseDSN with postgres and valid dsn", func(t *testing.T) {
			// Add test cases for ParseDSN function with postgres and valid dsn
			meta, err := ParseDSN("postgres", "postgres://user:password@localhost:5432/dbname")
			if err != nil {
				t.Errorf("ParseDSN() error = %v", err)
			}
			fmt.Println(meta)
			if meta[DBUser] != "user" {
				t.Errorf("ParseDSN() user = %v, want %v", meta[DBUser], "user")
			}
			if meta[TargetHost] != "localhost" {
				t.Errorf("ParseDSN() host = %v, want %v", meta[TargetHost], "localhost")
			}
			if meta[TargetPort] != "5432" {
				t.Errorf("ParseDSN() port = %v, want %v", meta[TargetPort], "5432")
			}
			if meta[DBName] != "dbname" {
				t.Errorf("ParseDSN() dbname = %v, want %v", meta[DBName], "dbname")
			}
		})
	})
}
