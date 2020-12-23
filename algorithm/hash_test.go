package algorithm

import (
	"fmt"
	"testing"
)

func TestHashSha(t *testing.T) {
	testcases := []struct {
		GivenHashType string
		Want          string
	}{
		{
			GivenHashType: "sha224",
			Want:          "ea09ae9cc6768c50fcee903ed054556e5bfc8347907f12598aa24193",
		},
		{
			GivenHashType: "sha256",
			Want:          "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			GivenHashType: "sha384",
			Want:          "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f",
		},
		{
			GivenHashType: "sha512",
			Want:          "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
	}

	for _, scenario := range testcases {
		t.Run(fmt.Sprintf("given hash(%s) returns(%s)", scenario.GivenHashType, scenario.Want), func(t *testing.T) {
			givenHashType := scenario.GivenHashType
			givenMessage := "hello"

			want := scenario.Want

			get, _ := hashSha(givenHashType, givenMessage)
			if want != fmt.Sprintf("%x", get) {
				t.Errorf("given hash %q message %q want %q but got %q\n", givenHashType, givenMessage, want, bToB64(get))
			}
		})
	}
}
