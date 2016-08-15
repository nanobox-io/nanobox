package models

import (
	"testing"
)

func TestComponentSave(t *testing.T) {
	// clear the components table when we're finished
	defer truncate("123")

	component := Component{
		AppID: "123",
		Name:  "web.main",
	}

	err := component.Save()
	if err != nil {
		t.Error(err)
	}

	// fetch the component
	component2 := Component{}

	if err = get("123", "web.main", &component2); err != nil {
		t.Errorf("failed to fetch component: %s", err.Error())
	}

	if component2.AppID != "123" {
		t.Errorf("component doesn't match")
	}
}

func TestComponentDelete(t *testing.T) {
	// clear the components table when we're finished
	defer truncate("123")

	component := Component{
		AppID: "123",
		Name:  "web.main",
	}

	if err := component.Save(); err != nil {
		t.Error(err)
	}

	if err := component.Delete(); err != nil {
		t.Error(err)
	}

	// make sure the component is gone
	keys, err := keys("123")
	if err != nil {
		t.Error(err)
	}

	if len(keys) > 0 {
		t.Errorf("component was not deleted")
	}
}

func TestFindComponentBySlug(t *testing.T) {
	// clear the components table when we're finished
	defer truncate("123")

	component := Component{
		AppID: "123",
		Name:  "web.main",
	}

	if err := component.Save(); err != nil {
		t.Error(err)
	}

	component2, err := FindComponentBySlug("123", "web.main")
	if err != nil {
		t.Error(err)
	}

	if component2.AppID != "123" {
		t.Errorf("did not load the correct component")
	}
}

func TestAllComponentsByApp(t *testing.T) {
	// clear the components table when we're finished
	defer truncate("1")
	defer truncate("2")

	component1 := Component{AppID: "1", Name: "web.main"}
	component2 := Component{AppID: "1", Name: "data.db"}
	component3 := Component{AppID: "2", Name: "web.main"}

	if err := component1.Save(); err != nil {
		t.Error(err)
	}
	if err := component2.Save(); err != nil {
		t.Error(err)
	}

	if err := component3.Save(); err != nil {
		t.Error(err)
	}

	components, err := AllComponentsByApp("1")
	if err != nil {
		t.Error(err)
	}

	if len(components) != 2 {
		t.Errorf("did not load all components, got %d", len(components))
	}
}
