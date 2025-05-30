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
				if meta[User] != "user" {
					t.Errorf("ParseDSN() user = %v, want %v", meta[User], "user")
				}
				if meta[Host] != "localhost" {
					t.Errorf("ParseDSN() host = %v, want %v", meta[Host], "localhost")
				}
				if meta[Port] != "3306" {
					t.Errorf("ParseDSN() port = %v, want %v", meta[Port], "3306")
				}
				if meta[Database] != "dbname" {
					t.Errorf("ParseDSN() dbname = %v, want %v", meta[Database], "dbname")
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
			if meta[User] != "user" {
				t.Errorf("ParseDSN() user = %v, want %v", meta[User], "user")
			}
			if meta[Host] != "localhost" {
				t.Errorf("ParseDSN() host = %v, want %v", meta[Host], "localhost")
			}
			if meta[Port] != "5432" {
				t.Errorf("ParseDSN() port = %v, want %v", meta[Port], "5432")
			}
			if meta[Database] != "dbname" {
				t.Errorf("ParseDSN() dbname = %v, want %v", meta[Database], "dbname")
			}
		})
	})
}
