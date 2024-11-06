package goption_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/goption"
	_ "github.com/mattn/go-sqlite3"
)

type Activity struct {
	ID          goption.Optional[int64]     `json:"id"`
	Time        goption.Optional[time.Time] `json:"time"`
	Description goption.Optional[string]    `json:"description"`
}

type WithScan struct {
	data int
}

func (c *WithScan) Scan(src any) error {
	c.data = src.(int)
	return nil
}

func (c WithScan) Value() (driver.Value, error) {
	return c.data, nil
}

var _ = Describe("Sql", func() {

	It("works on db", Label("Integration"), func() {
		var (
			T          = GinkgoT()
			file       = "activities.db"
			create_sql = `
  CREATE TABLE IF NOT EXISTS activities (
  id INTEGER PRIMARY KEY,
  time DATETIME,
  description TEXT
  );`
			createDB = func(db *sql.DB) error {
				if _, err := db.Exec(create_sql); err != nil {
					return err
				}
				return nil
			}

			deleteDB = func() error {
				return os.Remove(file)
			}

			createActivity = func(db *sql.DB, a *Activity) error {
				res, err := db.Exec("INSERT INTO activities VALUES(NULL,?,?);", a.Time, a.Description)
				if err != nil {
					return err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return err
				}

				a.ID = goption.Of(int64(id))

				return nil
			}

			findActivityByID = func(db *sql.DB, id int64) (*Activity, error) {
				// Query DB row based on ID
				row := db.QueryRow("SELECT id, time, description FROM activities WHERE id=?", id)

				// Parse row into Activity struct
				activity := Activity{}
				var err error
				if err = row.Scan(&activity.ID, &activity.Time, &activity.Description); err == sql.ErrNoRows {
					log.Printf("Id not found")
					return nil, fmt.Errorf("%v not found", id)
				}
				return &activity, err
			}
		)
		db, err := sql.Open("sqlite3", file)
		if err != nil {
			T.Log(err)
			T.FailNow()
		}
		if err := createDB(db); err != nil {
			T.Log(err)
			T.FailNow()
		}

		a := Activity{
			Time: goption.Of(time.Now()),
		}

		if err := createActivity(db, &a); err != nil {
			T.Log(err)
			T.FailNow()
		}

		id, _ := a.ID.Get()

		af, err := findActivityByID(db, id)
		if err != nil {
			T.Log(err)
			T.FailNow()
		}

		if err := deleteDB(); err != nil {
			T.Log(err)
			T.FailNow()
		}

		Expect(af.ID.IsPresent()).To(BeTrue())
		Expect(af.Time.IsPresent()).To(BeTrue())
		Expect(af.Description.IsPresent()).To(BeFalse())
	})

	Describe("Scan", func() {
		It("assigns given data to optional", func() {
			opt := goption.Empty[string]()

			Expect(opt.Scan("Hello!")).To(Succeed())
			Expect(opt.IsPresent()).To(BeTrue())
		})

		When("data is not valid", func() {
			It("create a empty optional from empty src", func() {
				opt := goption.Empty[string]()

				Expect(opt.Scan("")).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})

			It("create a empty optional from nil", func() {
				opt := goption.Empty[*string]()

				Expect(opt.Scan(nil)).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})
		})

		When("has a custom type", func() {
			It("assign it by its type", func() {
				type Money int

				opt := goption.Empty[Money]()

				Expect(opt.Scan(400)).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
			})
		})

		When("data implements its own scan method", func() {
			It("calls it to do transform", func() {
				opt := goption.Empty[WithScan]()

				Expect(opt.Scan(400)).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet().data).To(Equal(400))
			})
		})
	})

	Describe("Value", func() {
		When("has no valid data", func() {
			It("returns value as nil without error", func() {
				var opt = goption.Empty[string]()

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(BeNil())
			})
		})

		When("data is valid", func() {
			It("returns value without error", func() {
				var opt = goption.Of("data")

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).ToNot(BeNil())
			})
		})

		When("data implements its own value method", func() {
			It("calls it to do transform", func() {
				var opt = goption.Of(WithScan{data: 300})

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got.(int)).To(Equal(opt.MustGet().data))
			})
		})
	})

})
