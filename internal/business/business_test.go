package business

// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/stretchr/testify/assert"
// )

// func NewTempDatabase(t *testing.T) *pgxpool.Pool {
// 	t.Helper()

// 	conn, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/dbname")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create database
// 	_, err = conn.Exec(context.Background(), "CREATE DATABASE test_db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Connect to database
// 	conn.Close()

// 	conn, err = pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/test_db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create tables
// 	os.ReadDir(name)
// }

// func TestBusiness_GetCrawlerSeedDomains(t *testing.T) {
// 	conn, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/dbname")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer conn.Close()

// 	b := New(conn)

// 	// Test with count = 0
// 	domains, err := b.GetCrawlerSeedDomains(context.Background(), 0)
// 	assert.NoError(t, err)
// 	assert.Equal(t, initialSeedDomains, domains)

// 	// Test with count > 0
// 	domains, err = b.GetCrawlerSeedDomains(context.Background(), 2)
// 	assert.NoError(t, err)
// 	assert.Len(t, domains, 2)

// 	// Test with count > number of seed domains
// 	domains, err = b.GetCrawlerSeedDomains(context.Background(), 10)
// 	assert.NoError(t, err)
// 	assert.Len(t, domains, 10)
// }
