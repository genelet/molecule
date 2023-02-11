package protein

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/genelet/molecule/godbi"
	"github.com/genelet/molecule/rdb"

	_ "github.com/go-sql-driver/mysql"
)

func newTeamAdmin(molecule *godbi.Molecule) *Team {
	team := &Team{IsAdmin: true}
    team.AutoUserTable(molecule)
    return team
}

func newTeamUserTable(molecule *godbi.Molecule, uerIDName, userTable string) *Team {
	team := &Team{UserIDName: uerIDName}
    team.AutoUserTable(molecule, userTable)
    return team
}

func TestUserTableTeam(t *testing.T) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "summer"
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil { t.Fatal(err) }
	defer db.Close()

	molecule, err := rdb.NewMolecule(db, godbi.MySQL, dbName)
	if err != nil { t.Fatal(err) }

	admin := newTeamAdmin(molecule)
	generics, landings := admin.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv", "adv_attrname", "adv_attrvalue", "adv_balance", "adv_campaign", "adv_creative", "adv_fail", "adv_ip", "adv_item", "adv_media", "adv_targetname", "adv_targetvalue", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pay_alipay", "pay_cc", "pay_cheque", "pay_payment", "pay_wechat", "pub", "pub_fail", "pub_ip", "pub_referer", "pub_site", "pub_slot", "pub_weight", "pub_white"}, "") || landings != nil {
		t.Errorf("%#v", generics)
		t.Errorf("%#v", landings)
	}

	adv := newTeamUserTable(molecule, "adv_id", "adv")
	generics, landings = adv.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv_balance", "adv_fail", "adv_ip", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pub", "pub_fail", "pub_ip", "pub_referer", "pub_site", "pub_slot"}, "") {
		t.Errorf("%#v", generics)
	}
	if strings.Join(landings, "") != strings.Join([]string{"adv_attrname", "adv_campaign", "pay_payment"}, "") {
		t.Errorf("%#v", landings)
	}
}

func getProtein() (*Protein, *sql.DB, error) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "summer"
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil { return nil, nil, err }

	molecule, err := rdb.NewMolecule(db, godbi.MySQL, dbName)
	if err != nil { return nil, nil, err }

	protein := NewDefaultProtein(godbi.MySQL, map[string]string{"adv":"adv_id", "pub":"pub_id"})
	err = protein.AutoTeams(molecule, map[string]string{"adv":"adv", "pub":"pub"})
	return protein, db, err
}

func TestUserTableProtein(t *testing.T) {
	protein, _, err := getProtein()
    if err != nil { t.Fatal(err) }

	admin := protein.Teams["admin"]
	generics, landings := admin.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv", "adv_attrname", "adv_attrvalue", "adv_balance", "adv_campaign", "adv_creative", "adv_fail", "adv_ip", "adv_item", "adv_media", "adv_targetname", "adv_targetvalue", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pay_alipay", "pay_cc", "pay_cheque", "pay_payment", "pay_wechat", "pub", "pub_fail", "pub_ip", "pub_referer", "pub_site", "pub_slot", "pub_weight", "pub_white"}, "") || landings != nil {
		t.Errorf("%#v", generics)
		t.Errorf("%#v", landings)
	}

	adv := protein.Teams["adv"]
	generics, landings = adv.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv_balance", "adv_fail", "adv_ip", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pub", "pub_fail", "pub_ip", "pub_referer", "pub_site", "pub_slot"}, "") {
		t.Errorf("%#v", generics)
	}
	if strings.Join(landings, "") != strings.Join([]string{"adv_attrname", "adv_campaign", "pay_payment"}, "") {
		t.Errorf("%#v", landings)
	}

	pub := protein.Teams["pub"]
	generics, landings = pub.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv", "adv_attrname", "adv_attrvalue", "adv_balance", "adv_campaign", "adv_creative", "adv_fail", "adv_ip", "adv_item", "adv_media", "adv_targetname", "adv_targetvalue", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pay_alipay", "pay_cc", "pay_cheque", "pay_payment", "pay_wechat", "pub_fail", "pub_ip"}, "") {
		t.Errorf("%#v", generics)
	}
	if strings.Join(landings, "") != strings.Join([]string{"pub_site"}, "") {
		t.Errorf("%#v", landings)
	}

	public := protein.Teams["public"]
	generics, landings = public.Dashboard()
	if strings.Join(generics, "") != strings.Join([]string{"ac", "add_address", "admin", "adv_balance", "adv_fail", "adv_ip", "agent", "ch_ac", "ch_belong", "cron_halfhour", "daily_adv", "daily_log", "daily_pub", "daily_pub_adv", "def_channel", "def_city", "def_continent", "def_country", "def_dma", "def_entitytype", "def_isp", "def_paytype", "def_size", "def_state", "his_balance", "his_payment", "ip", "ledger_adv", "ledger_log", "ledger_pub", "ledger_pub_adv", "pub_fail", "pub_ip"}, "") || landings != nil {
		t.Errorf("%#v", generics)
		t.Errorf("%#v", landings)
	}
}
