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

type ActivityZero struct {
	ID          goption.Zero[int64]     `json:"id"`
	Time        goption.Zero[time.Time] `json:"time"`
	Description goption.Zero[string]    `json:"description"`
}

type WithScan struct {
	data int
}

func (c *WithScan) Scan(src any) error {
	c.data = int(src.(int64))
	return nil
}

func (c WithScan) Value() (driver.Value, error) {
	return int64(c.data), nil
}

var _ = Describe("Sql", func() {
	It("preserves defined string SQL behavior", func() {
		var opt goption.Optional[definedTextString]
		var zero goption.Zero[definedTextString]
		Expect(opt.Scan("hello")).To(Succeed())
		Expect(zero.Scan("hello")).To(Succeed())
		Expect(opt.MustGet()).To(Equal(definedTextString("hello")))
		Expect(zero.MustGet()).To(Equal(definedTextString("hello")))
		optValue, err := opt.Value()
		Expect(err).NotTo(HaveOccurred())
		Expect(optValue).To(Equal("hello"))
		zeroValue, err := zero.Value()
		Expect(err).NotTo(HaveOccurred())
		Expect(zeroValue).To(Equal("hello"))
		Expect(opt.Scan("")).To(Succeed())
		Expect(zero.Scan("")).To(Succeed())
		Expect(opt.IsPresent()).To(BeTrue())
		Expect(zero.IsPresent()).To(BeFalse())
		Expect(opt.Scan(nil)).To(Succeed())
		Expect(zero.Scan(nil)).To(Succeed())
		Expect(opt.IsPresent()).To(BeFalse())
		Expect(zero.IsPresent()).To(BeFalse())
	})

	It("works on db for Optional[T]", Label("Integration"), func() {
		var (
			T          = GinkgoT()
			file       = "activities_optional.db"
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

				a.ID = goption.Some(id)

				return nil
			}

			findActivityByID = func(db *sql.DB, id int64) (*Activity, error) {
				row := db.QueryRow("SELECT id, time, description FROM activities WHERE id=?", id)

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
			Time: goption.Some(time.Now()),
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

	It("works on db for Zero[T]", Label("Integration"), func() {
		var (
			T          = GinkgoT()
			file       = "activities_zero.db"
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

		// Insert empty description
		a := ActivityZero{
			Time:        goption.SomeZero(time.Now()),
			Description: goption.SomeZero(""),
		}

		res, err := db.Exec("INSERT INTO activities VALUES(NULL,?,?);", a.Time, a.Description)
		Expect(err).ToNot(HaveOccurred())
		id, err := res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		row := db.QueryRow("SELECT id, time, description FROM activities WHERE id=?", id)
		var af ActivityZero
		Expect(row.Scan(&af.ID, &af.Time, &af.Description)).To(Succeed())

		Expect(deleteDB()).To(Succeed())

		Expect(af.ID.IsPresent()).To(BeTrue())
		Expect(af.Time.IsPresent()).To(BeTrue())
		Expect(af.Description.IsPresent()).To(BeFalse())
	})

	Describe("Scan for Optional[T] and Zero[T]", func() {
		It("assigns given data to Optional[string]", func() {
			var opt goption.Optional[string]

			Expect(opt.Scan("Hello!")).To(Succeed())
			Expect(opt.IsPresent()).To(BeTrue())
			Expect(opt.MustGet()).To(Equal("Hello!"))
		})

		When("scanning zero values vs nil into Optional[string]", func() {
			It("Optional[string] treats empty string as valid present", func() {
				var opt goption.Optional[string]

				Expect(opt.Scan("")).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet()).To(Equal(""))
			})

			It("Optional[string] treats nil as invalid/null", func() {
				opt := goption.Some("initial")

				Expect(opt.Scan(nil)).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})
		})

		When("scanning zero values vs nil into Zero[string]", func() {
			It("Zero[string] treats empty string as absent/empty", func() {
				var opt goption.Zero[string]

				Expect(opt.Scan("")).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})

			It("Zero[string] treats non-empty string as present", func() {
				var opt goption.Zero[string]

				Expect(opt.Scan("hello")).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet()).To(Equal("hello"))
			})

			It("Zero[string] treats nil as absent/empty", func() {
				opt := goption.SomeZero("initial")

				Expect(opt.Scan(nil)).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})
		})

		When("has a custom type", func() {
			It("assigns it by its convertible type", func() {
				type Money int

				var opt goption.Optional[Money]

				Expect(opt.Scan(400)).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet()).To(Equal(Money(400)))
			})
		})

		When("data implements its own scan method", func() {
			It("calls it to do transform for Optional[T]", func() {
				var opt goption.Optional[WithScan]

				Expect(opt.Scan(int64(400))).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet().data).To(Equal(400))
			})

			It("calls it to do transform for Zero[T]", func() {
				var opt goption.Zero[WithScan]

				Expect(opt.Scan(int64(400))).To(Succeed())
				Expect(opt.IsPresent()).To(BeTrue())
				Expect(opt.MustGet().data).To(Equal(400))
			})
		})
	})

	Describe("Value for Optional[T] and Zero[T]", func() {
		When("has no valid data", func() {
			It("returns value as nil without error for Optional[T]", func() {
				var opt goption.Optional[string]

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(BeNil())
			})

			It("returns value as nil without error for Zero[T]", func() {
				var opt goption.Zero[string]

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(BeNil())
			})
		})

		When("data is valid", func() {
			It("returns value without error for Optional[T]", func() {
				opt := goption.Some("data")

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal("data"))
			})

			It("returns value without error for Zero[T]", func() {
				opt := goption.SomeZero("data")

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal("data"))
			})
		})

		When("data implements its own value method", func() {
			It("calls it to do transform", func() {
				opt := goption.Some(WithScan{data: 300})

				got, err := opt.Value()

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(int64(300)))
			})
		})
	})
})
