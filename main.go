package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/tcnksm/go-latest"
	"os"
)

func main() {
	var sconstraint string
	var sver string
	var owner string
	var repository string
	var quiet bool
	var filterfunc latest.TagFilterFunc

	flag.BoolVar(&quiet, "quiet", false, "just prints the latest version to use")
	flag.StringVar(&sconstraint, "match", "", `version matching (eg. >=3.6.0,<3.7)`)
	flag.StringVar(&sver, "version", "", `current version (eg. 3.6.0)`)
	flag.StringVar(&owner, "owner", "python", "github owner (eg. python)")
	flag.StringVar(&repository, "repository", "cpython", "github repository (eg. cpython)")
	flag.Parse()
	if len(sver) == 0 || len(owner) == 0 || len(repository) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(sconstraint) != 0 {
		constraint, err := version.NewConstraint(sconstraint)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Constraint error:", err)
			os.Exit(1)
		}
		filterfunc = func(s string) bool {
			v, err := version.NewVersion(s)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return false
			}
			return constraint.Check(v)
		}
	} else {
		filterfunc = func(s string) bool { return true }
	}
	githubTag := &latest.GithubTag{
		Owner:         owner,
		Repository:    repository,
		TagFilterFunc: filterfunc,
	}
	res, err := latest.Check(githubTag, sver)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching latest versions:", err)
		os.Exit(3)
	}
	if quiet {
		fmt.Println(res.Current)
	} else {
		if res.Outdated {
			fmt.Println(sver, "is not latest, you should upgrade to", res.Current)
		}
	}
	os.Exit(0)
}
