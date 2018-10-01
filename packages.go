package main

import (
	pkgdb "github.com/lemmi/xbpspkgdb"
	"log"
)

var architectures = []string{
	"x86_64",
	"x86_64-musl", 
	"i686",
	"aarch64",
	"aarch64-musl",
	"armv6l",
	"armv6l-musl",
	"armv7l",
	"armv7l-musl",
}

type Pkg struct {
	Repo string
	Pkgname string
	pkgdb.Package
}

type PkgIndex struct {
	packages []Pkg
}

func (p Pkg) String() string {
	return p.Pkgver
}

var myrepos = map[string]pkgdb.Pkgdb {
	"x86_64": nil,
	"x86_64-musl": nil, 
}

func (pi PkgIndex) Architectures() []string {
	return architectures
}

func (pi PkgIndex) Filter() PkgIndex {
	var newpi PkgIndex
	return newpi
}

func (pi PkgIndex) Sort() PkgIndex {
	var newpi PkgIndex
	return newpi
}

func NewPkgIndex() PkgIndex {
	var pi PkgIndex
	for name, _ := range myrepos {
		log.Println("Loading", name + "-repodata")
		var err error
		myrepos[name], err = pkgdb.DecodeRepoDataFile(name + "-repodata")
		if err != nil {
			log.Fatal(err)
		}
		for pkgname, pkg := range myrepos[name] {
			ref, ok := pkgindex[pkgname]
			if !ok {
				pkgindex[pkgname] = make(map[string]pkgdb.Package)
				ref = pkgindex[pkgname]
			}
			ref[name] = pkg
		}
	}
	return pi
}
