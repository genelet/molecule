package engine

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/genelet/team/micro"
)

func TestColorful(t *testing.T) {
	bs, err := ioutil.ReadFile("tmp/adv/adv-adv.json")
	if err != nil {
		t.Fatal(err)
	}
	colorful := Colorful{}
	err = json.Unmarshal(bs, &colorful)
	if err != nil {
		t.Fatal(err)
	}
	if !colorful.Color.IsUser || colorful.Atom.Table.Pks[0] != "adv_id" {
		t.Errorf("%#v", colorful)
		t.Errorf("%#v", colorful.Atom)
		t.Errorf("%#v", colorful.Atom.Table.Pks[0])
	}
}

func TestColorfulEndpoint(t *testing.T) {
	ctx := context.Background()

	myURL, err := url.Parse("file:///tmp/adv/adv-adv.json")
	if err != nil {
		t.Fatal(err)
	}

	service := newColorfulFileService(myURL)
	conf, endpoint, err := service.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	colorful := conf.(*Colorful)
	if colorful.Protected != "adv_id" || colorful.Atom.GetTable().Pks[0] != "adv_id" {
		t.Errorf("%#v", colorful)
		t.Errorf("%#v", endpoint)
	}
}

func TestColorfulCreation(t *testing.T) {
	ctx := context.Background()

	myURL, err := url.Parse("file:///tmp/adv/adv-adv.json")
	if err != nil {
		t.Fatal(err)
	}

	/*
		endpoint := &micro.Endpoint{
			MetaType: micro.METASingle,
			SingleDescriptor:&micro.Descriptor{
				ObjectName: "Colorful",
				NextService: "adv-adv.json",
			},
		}
	*/

	creation := func(colorful micro.Resolver, args ...interface{}) (micro.Microservice, error) {
		return micro.NewZeroFileSingle(colorful, args...), nil
	}

	service, err := NewColorfulService(creation, myURL)
	if err != nil {
		t.Fatal(err)
	}
	conf, endpoint, err := service.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	colorful := conf.(*Colorful)
	if colorful.Protected != "adv_id" || colorful.Atom.GetTable().Pks[0] != "adv_id" {
		t.Errorf("%#v", colorful)
		t.Errorf("%#v", endpoint)
	}
}

func TestColorfulItem(t *testing.T) {
    u, err := url.Parse("file:///testdata/temp/adv_campaign/adv-adv_campaign.json")
    if err != nil {
        t.Fatal(err)
    }
    colorful := new(Colorful)
    cs, err := colorful.getFileSingle(u)
    if err != nil {
        t.Fatal(err)
    }

    ctx := context.Background()
    iconf, endpoint, err := cs.Read(ctx)
    if err != nil {
        t.Fatal(err)
    }
    conf := iconf.(*Colorful)

    table1 := conf.Atom.GetTable()
    if table1.TableName != "adv_campaign" {
        t.Errorf("%#v", table1)
        t.Errorf("%#v", conf)
        t.Errorf("%#v", endpoint)
    }

    u, err = url.Parse("file:///testdata/adv-adv_campaign.json")
    if err != nil {
        t.Fatal(err)
    }
    cs.SetMyURL(u)
    err = cs.Write(ctx)
    if err != nil {
        t.Errorf("%v", err)
    }
}
