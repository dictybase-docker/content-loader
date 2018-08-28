package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func validateImport(c *cli.Context) error {
	for _, p := range []string{"akey", "skey", "rp", "namespace", "email"} {
		if !c.IsSet(p) {
			return cli.NewExitError(
				fmt.Sprintf("missing required argument %s", p),
				2,
			)
		}
	}
	return nil
}
