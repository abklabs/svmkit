package agave

import (
	"testing"
)

func TestValidatorEnv(t *testing.T) {
	{
		v := ValidatorEnv{}

		if v.ToString() != "" {
			t.Error("empty validator did not return an empty stirng")
		}
	}

	{
		m := Metrics{
			URL:      "noproto://nowhere",
			Database: "nodb",
			User:     "notauser",
			Password: "notapassword",
		}

		v := ValidatorEnv{
			Metrics: &m,
		}

		if v.ToString() != `SOLANA_METRICS_CONFIG="host=noproto://nowhere,db=nodb,u=notauser,p=notapassword"` {
			t.Error("validator env output didn't match expectations")
		}
	}
}
