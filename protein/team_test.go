package engine

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/genelet/molecule/godbi"
	"github.com/genelet/molecule/rdb"

	"github.com/genelet/team/micro"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabaseAdmin(t *testing.T) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "summer"
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	molecule, err := rdb.NewMolecule(db, godbi.MySQL, dbName)
	if err != nil {
		t.Fatal(err)
	}

	adv := NewTeamUserTable(molecule, "adv", "adv_id")

	generics, landings := adv.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv_balance", "adv_fail", "adv_ip", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pub", "pub_fail", "pub_ip", "pub_referer", "pub_site", "pub_slot"}, "") {
		t.Errorf("%#v", generics)
	}
	if strings.Join(landings, "") != strings.Join([]string{"adv_attrname", "adv_campaign", "pay_payment"}, "") {
		t.Errorf("%#v", landings)
	}

	u, err := url.Parse("file:///testdata/temp/adv.json")
	if err != nil {
		t.Fatal(err)
	}
	service, err := adv.getFileService(u)
	if err != nil {
		t.Fatal(err)
	}

	item := service.GetObject().(*Team).GetColorful("adv")
	if item.IsUser != true || item.Protected != "adv_id" || item.Atom.GetTable().TableName != "adv" {
		t.Errorf("%#v", item)
	}

	err = service.Write(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	u, err = url.Parse("file:///testdata/temp/admin.json")
	if err != nil {
		t.Fatal(err)
	}
	admin := NewTeamUserTable(molecule, "admin")
	serviceA, err := admin.getFileService(u)
	if err != nil {
		t.Fatal(err)
	}
	err = serviceA.Write(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTeamWrite(t *testing.T) {
	type Role struct {
		Team
		Redirect string `json:"redirect"`
	}

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (id int auto_increment not null primary key,
        x varchar(8), y varchar(8), z varchar(8))`)
	db.Exec(`drop table if exists m_b`)
	db.Exec(`CREATE TABLE m_b (tid int auto_increment not null primary key,
        child varchar(8), id int, foreign key (id) REFERENCES m_a(id) ON DELETE CASCADE)`)
	db.Exec(`drop table if exists m_c`)
	db.Exec(`CREATE TABLE m_c (itemid int auto_increment not null primary key,
        item varchar(8), childid int, foreign key (childid) REFERENCES m_b(tid) ON DELETE CASCADE)`)
	db.Exec(`drop table if exists m_d`)
	db.Exec(`CREATE TABLE m_d (createid int auto_increment not null primary key,
        creative varchar(8), itemid int, foreign key (itemid) REFERENCES m_c(itemid) ON DELETE CASCADE)`)
	db.Exec(`drop table if exists m_e`)
	db.Exec(`CREATE TABLE m_e (fid int auto_increment not null primary key,
        friend varchar(8), id int, foreign key (id) REFERENCES m_a(id) ON DELETE CASCADE)`)

	molecule, err := rdb.NewMolecule(db, godbi.MySQL, dbName)
	if err != nil {
		t.Fatal(err)
	}

	team := NewTeamUserTable(molecule, "xTeam", "id", "m_a")
	//role := &Role{Team:*team, Redirect:"/index.html"}

	u, err := url.Parse("file:///tmp/xTeam.json")
	service, err := team.getFileService(u)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Write(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	db.Exec(`drop table if exists m_e`)
	db.Exec(`drop table if exists m_d`)
	db.Exec(`drop table if exists m_c`)
	db.Exec(`drop table if exists m_b`)
	db.Exec(`drop table if exists m_a`)
}

func TestTeamEndpoint(t *testing.T) {
	ctx := context.Background()

	teamName := "adv"
	tables := []string{"adv", "adv_campaign", "adv_item"}

	myURL, err := url.Parse("file:///testdata/temp/" + teamName + ".json")
	if err != nil {
		t.Fatal(err)
	}

	service := newTeamFileService(myURL, tables)
	conf, endpoint, err := service.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	team := conf.(*Team)
	if team.UserIDName != "adv_id" ||
		team.Colorfuls[0].Atom.GetTable().Pks[0] != "adv_id" ||
		team.Colorfuls[1].Atom.GetTable().Pks[0] != "campaign_id" ||
		team.Colorfuls[2].Atom.GetTable().Pks[0] != "item_id" {
		t.Errorf("%#v", team)
		t.Errorf("%#v", endpoint)
	}
}

func TestTeamCreation(t *testing.T) {
	ctx := context.Background()

	teamName := "adv"
	tables := []string{"adv", "adv_campaign", "adv_item"}

	myURL, err := url.Parse("file:///testdata/temp/" + teamName + ".json")
	if err != nil {
		t.Fatal(err)
	}

	creation := func(colorful micro.Resolver, args ...interface{}) (micro.Microservice, error) {
		return micro.NewZeroFileSingle(colorful, args...), nil
	}

	service, err := NewTeamService(creation, myURL, tables)
	if err != nil {
		t.Fatal(err)
	}
	conf, endpoint, err := service.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	team := conf.(*Team)
	if team.UserIDName != "adv_id" ||
		team.Colorfuls[0].Atom.GetTable().Pks[0] != "adv_id" ||
		team.Colorfuls[1].Atom.GetTable().Pks[0] != "campaign_id" ||
		team.Colorfuls[2].Atom.GetTable().Pks[0] != "item_id" {
		t.Errorf("%#v", team)
		t.Errorf("%#v", endpoint)
	}
}

func TestTeamSlice(t *testing.T) {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := "summer"
    db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    molecule, err := rdb.NewMolecule(db, godbi.MySQL, dbName)
    if err != nil {
        t.Fatal(err)
    }
    adv := NewTeamUserTable(molecule, "adv", "adv_id")

    teamURL, err := url.Parse("file:///testdata/temp/adv.json")
    if err != nil {
        t.Fatal(err)
    }
    cs, err := adv.getFileService(teamURL)
    if err != nil {
        t.Fatal(err)
    }

    ctx := context.Background()
    err = cs.Write(ctx)
    if err != nil {
        t.Fatal(err)
    }

    iconf, endpoint, err := cs.Read(ctx)
    if err != nil {
        t.Fatal(err)
    }
    conf := *(iconf.(*Team))

    table1 := conf.Colorfuls[0].Atom.GetTable()
    table2 := conf.Colorfuls[1].Atom.GetTable()
    if table1.TableName != "ac" || table2.Pks[0] != "address_id" {
        t.Errorf("%#v", table1)
        t.Errorf("%#v", table2)
        t.Errorf("%#v", endpoint)
    }

    teamURL, err = url.Parse("file:///testdata/temp2/adv.json")
    if err != nil {
        t.Fatal(err)
    }
    cs.SetMyURL(teamURL)
    err = cs.Write(ctx)
    if err != nil {
        t.Fatal(err)
    }
}
