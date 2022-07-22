//go:build integration

package user

import (
	"USATUKirill96/gridgo/users/pkg/pagination"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var fixtures = []*User{
	&User{0, "ServiceTest", 4.3191, 3.4815},
	&User{0, "Close1", 5.000, 3.000}, // 96.63 km
	&User{0, "Close2", 5.000, 4.000}, // 95.05 km
	&User{0, "Far1", 6.000, 3.000},   // 194.4 km
	&User{0, "Far2", 6.000, 6.000},   // 335.7 km
}

func NewTestDatabase() TestDatabase {

	err := godotenv.Load("./../../../.env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", os.Getenv("POSTGRES_DB_TEST_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return TestDatabase{db}
}

type TestDatabase struct {
	db *sql.DB
}

func (td TestDatabase) NewRepository() Repository {
	return Repository{td.db}
}

func (td TestDatabase) Setup() {

	for _, u := range fixtures {

		stmt := `
	   INSERT INTO app_user (username, longitude, latitude) 
	   VALUES ($1, $2, $3) 
	`
		_, err := td.db.Exec(stmt, u.Username, u.Longitude, u.Latitude)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func (td TestDatabase) TearDown() {
	stmt := `TRUNCATE TABLE app_user`
	_, err := td.db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func assertUsers(u1, u2 User) (string, bool) {
	formatError := func(f1, f2 interface{}) string {
		return fmt.Sprintf("Asserting users: expected %v, got %v", f1, f2)
	}
	if u1.Username != u2.Username {
		return formatError(u1.Username, u2.Username), false
	}
	if u1.Longitude != u2.Longitude {
		return formatError(u1.Longitude, u2.Longitude), false
	}
	if u1.Latitude != u2.Latitude {
		return formatError(u1.Latitude, u2.Latitude), false
	}
	return "", true
}

func TestRepository_Insert(t *testing.T) {

	td := NewTestDatabase()
	r := td.NewRepository()
	defer td.TearDown()

	testUser := User{
		Username:  "TestInsert",
		Longitude: 0.123,
		Latitude:  -0.123,
	}

	// Perform tested logic
	u, err := r.Insert(testUser)
	if err != nil {
		log.Fatal(err)
	}

	// Assert returned object is correct
	message, ok := assertUsers(testUser, *u)
	if !ok {
		t.Error(message)
	}

	// Assert data is saved
	stmt := `
	SELECT id, username, longitude, latitude
	  FROM app_user
	 WHERE username = $1
	`

	su := &User{}
	err = r.db.QueryRow(stmt, testUser.Username).Scan(&su.ID, &su.Username, &su.Longitude, &su.Latitude)
	if err != nil {
		t.Error(err)

	}
	message, ok = assertUsers(testUser, *su)
	if !ok {
		t.Error(message)
	}
}

func TestRepository_ByDistance(t *testing.T) {

	const (
		target = iota
		close1
		close2
		far1
		far2
	)

	td := NewTestDatabase()
	r := td.NewRepository()
	defer td.TearDown()

	td.Setup()

	cases := []struct {
		distance int
		expected []*User
	}{
		{50, []*User{}},
		{100, []*User{fixtures[close1], fixtures[close2]}},
		{200, []*User{fixtures[close1], fixtures[close2], fixtures[far1]}},
		{400, []*User{fixtures[close1], fixtures[close2], fixtures[far1], fixtures[far2]}},
	}

	for _, tc := range cases {
		// perform logic
		result, err := r.ByDistance(*fixtures[target], tc.distance, pagination.Pagination{})
		if err != nil {
			t.Error(err)
		}
		// assert number of results
		// TODO: assert exact users match
		if len(result) != len(tc.expected) {
			t.Errorf("Incorrect length of result: expected %v, got %v", len(tc.expected), len(result))
		}
	}
}
func TestRepository_ByDistance_pagination(t *testing.T) {

	td := NewTestDatabase()
	r := td.NewRepository()
	defer td.TearDown()

	td.Setup()

	for i := 0; i < 15; i++ {
		r.Insert(User{0, fmt.Sprintf("ByDistance%v", i), 4.3191, 3.4815})
	}

	cases := []struct {
		pg       pagination.Pagination
		expected int
	}{
		{pagination.Pagination{Limit: 5, Offset: 0}, 5},  // Limit
		{pagination.Pagination{Limit: 0, Offset: 10}, 5}, // Offset
		{pagination.Pagination{Limit: 0, Offset: 35}, 0}, // Offset is too large
		{pagination.Pagination{}, 15},                    // Empty pagination
	}
	for _, tc := range cases {
		res, err := r.ByDistance(*fixtures[0], 5, tc.pg)
		if err != nil {
			t.Error(err)
		}
		if len(res) != tc.expected {
			t.Errorf("Unexpected number of results. Expected: %v, got: %v", tc.expected, len(res))
		}
	}

}

func TestRepository_ByUsername(t *testing.T) {

	td := NewTestDatabase()
	r := td.NewRepository()
	defer td.TearDown()

	td.Setup()

	for _, f := range fixtures {
		// Perform logic
		u, err := r.ByUsername(f.Username)
		if err != nil {
			t.Error(err)
		}
		// Assert data
		message, ok := assertUsers(*f, *u)
		if !ok {
			t.Error(message)
		}
	}
}

func TestRepository_Update(t *testing.T) {

	const testedUser = 1
	const newLongitude = 12.123321

	td := NewTestDatabase()
	r := td.NewRepository()
	defer td.TearDown()
	td.Setup()

	u, err := r.ByUsername(fixtures[testedUser].Username)
	if err != nil {
		t.Error(t)
	}

	u.Longitude = newLongitude
	// Perform logic
	u, err = r.Update(*u)
	if err != nil {
		t.Error(err)
	}
	// Assert function return
	if u.Longitude != newLongitude {
		t.Errorf("Incorrect return from r.Update. Expected: %v, got: %v", newLongitude, u.Longitude)
	}
	// Assert database save
	u, err = r.ByUsername(fixtures[testedUser].Username)
	if err != nil {
		t.Error(err)
	}
	if u.Longitude != newLongitude {
		t.Errorf("Incorrect data saved. Expected longitude: %v, got: %v", newLongitude, u.Longitude)
	}
}
