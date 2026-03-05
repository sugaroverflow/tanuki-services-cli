package formatter

import (
	"testing"

	"github.com/example/tanuki/internal/catalog"
)

func TestListEmpty(t *testing.T) {
	List(nil)
	List([]catalog.Service{})
}

func TestListWithServices(t *testing.T) {
	svcs := []catalog.Service{
		{Name: "payments-api", Version: "1.0", Team: "payments", Owner: "jane"},
		{Name: "auth-service", Version: "2.0", Team: "platform", Owner: "alice"},
	}
	List(svcs)
}

func TestStatusNil(t *testing.T) {
	Status(nil)
}

func TestOwnersNil(t *testing.T) {
	Owners(nil)
}
