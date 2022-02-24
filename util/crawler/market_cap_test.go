package crawler

import (
	"fmt"
	"testing"
)

func TestQueryEnterpriseLatestName(t *testing.T) {
	name, err := QueryEnterpriseLatestName("陌陌")
	fmt.Println(name, err)
}
