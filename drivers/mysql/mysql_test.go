//go:build mysql

package mysql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k1LoW/tbls/schema"
	"github.com/xo/dburl"
)

var s *schema.Schema
var db *sql.DB

func TestMain(m *testing.M) {
	s = &schema.Schema{
		Name: "testdb",
	}
	db, _ = dburl.Open("my://root:mypass@localhost:33308/testdb")
	defer db.Close()
	exit := m.Run()
	if exit != 0 {
		os.Exit(exit)
	}
}

func TestAnalyzeView(t *testing.T) {
	driver, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	err = driver.Analyze(s)
	if err != nil {
		t.Fatal(err)
	}
	view, _ := s.FindTableByName("post_comments")
	if got := view.Def; got == "" {
		t.Errorf("got not empty string.")
	}
}

func TestExtraDef(t *testing.T) {
	driver, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := driver.Analyze(s); err != nil {
		t.Fatal(err)
	}
	tbl, _ := s.FindTableByName("comments")
	{
		c, _ := tbl.FindColumnByName("id")
		got := c.ExtraDef
		if want := "auto_increment"; got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
	}
	{
		c, _ := tbl.FindColumnByName("post_id_desc")
		got := c.ExtraDef
		if want := "GENERATED ALWAYS AS (`post_id` * -(1)) VIRTUAL"; got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
	}
}

func TestInfo(t *testing.T) {
	driver, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	d, err := driver.Info()
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "mysql" {
		t.Errorf("got %v\nwant %v", d.Name, "mysql")
	}
	if d.DatabaseVersion == "" {
		t.Errorf("got not empty string.")
	}
}

func TestCheckConstraintsOrder(t *testing.T) {
	if _, err := db.Exec(`DROP TABLE IF EXISTS tbls_check_order`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`
CREATE TABLE tbls_check_order (
  a int,
  b int,
  CONSTRAINT chk_tbls_check_order_b CHECK (b > 0),
  CONSTRAINT chk_tbls_check_order_a CHECK (a > 0)
)`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = db.Exec(`DROP TABLE IF EXISTS tbls_check_order`)
	}()

	driver, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	sc := &schema.Schema{Name: "testdb"}
	if err := driver.Analyze(sc); err != nil {
		t.Fatal(err)
	}

	tbl, err := sc.FindTableByName("tbls_check_order")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range tbl.Constraints {
		if c.Type == "CHECK" {
			got = append(got, c.Name)
		}
	}
	want := []string{"chk_tbls_check_order_a", "chk_tbls_check_order_b"}
	if len(got) != len(want) {
		t.Fatalf("got %v\nwant %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v\nwant %v", got, want)
			break
		}
	}
}
