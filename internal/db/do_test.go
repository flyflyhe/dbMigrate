package db

import (
	"fmt"
	"github.com/dbMigrate/v2/config"
	"github.com/dbMigrate/v2/internal/db/connection"
	"testing"
)

func init() {
	config.InitConfig("../../config.yaml")
}

func TestWrapper_AllTables(t *testing.T) {
	db0, err := connection.InitDb(config.GetApp().DbConfig.Mysql0)
	if err != nil {
		t.Error(err)
	}

	w := &Wrapper{db0}
	if result, err := w.AllTables(); err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestWrapper_TableSchema(t *testing.T) {
	db0, err := connection.InitDb(config.GetApp().DbConfig.Mysql0)
	if err != nil {
		t.Error(err)
	}

	w := &Wrapper{db0}
	if result, err := w.TableSchema("person"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}

func TestWrapper_ChangeDDL(t *testing.T) {
	db0, _ := connection.InitDb(config.GetApp().DbConfig.Mysql0)

	w0 := &Wrapper{db0}
	if result, err := w0.TableSchema("person"); err != nil {
		t.Error(err)
	} else {
		nDDL := w0.ChangeDDL("person", "person1", result)
		t.Log(nDDL)
	}
}
