package main

import (
	"github.com/dustin/go-humanize"
	pkgdb "github.com/lemmi/xbpspkgdb"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"fmt"
)

var pkgindex = map[string]map[string]pkgdb.Package{}

type pkgData struct {
	Pkgname string
	Versions map[string]string
	Pkg pkgdb.Package
}

type page struct {
	Text string
	URL string
}

type indexData struct {
	Pkgs []pkgdb.Package
	Architectures []string
	Pages []template.HTML
}

var archs = []string{
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

type repository struct {
	arch string
	name string
}

// var repos = []repository{
// 	{ "x86_64", "" },
// 	{ "x86_64-musl", "" },
// }

var repos = map[string]pkgdb.Pkgdb {
	"x86_64": nil,
	"x86_64-musl": nil, 
	// "i686": nil,
	// "aarch64": nil,
	// "aarch64-musl": nil,
	// "armv6l": nil,
	// "armv6l-musl": nil,
	// "armv7l": nil,
	// "armv7l-musl": nil,
}



var pathPkg = regexp.MustCompile("^/pkg/([a-zA-Z0-9_-]+)(/([a-zA-Z0-9_-]+))?$")

var templates = template.Must(
	template.New("").Funcs(template.FuncMap{
		"pkgname": tmplPkgname,
		"ghtemplate": tmplGithubTemplate,
		"strjoin": strings.Join,
		"humanbytes": func (n int) string {
			return humanize.Bytes(uint64(n))
		},
	}).ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/pkg.html",
		"templates/search.html",
	))

func tmplPkgname(pkgver string) string {
	idx := strings.LastIndex(pkgver, "-")
	return pkgver[:idx]
}

func tmplGithubTemplate(sourcerev string) string {
	a := strings.Split(sourcerev, ":")
	return "https://github.com/void-linux/void-packages/blob/"+a[1]+"/srcpkgs/"+a[0]+"/template"
}

func renderTemplate(w http.ResponseWriter, f string, d interface{}) {
	err := templates.ExecuteTemplate(w, f, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var per_page = 20

func index(w http.ResponseWriter, r *http.Request) {
	var pkgs []pkgdb.Package
	var i int

	_ = len(pkgindex) / per_page

	for _, ref := range pkgindex {
		for _, pkg := range ref {
			pkgs = append(pkgs, pkg)
			i++
			if i > per_page {
				break
			}
		}
		if i > per_page {
			break
		}
	}

	pages := make([]template.HTML, 5)
	for i = 0; i < 5; i++ {
		pages[i] = template.HTML(fmt.Sprintf(`<a>%d</a>`, i))
	}

	renderTemplate(w, "index.html", indexData {
		Architectures: archs,
		Pkgs: pkgs,
		Pages: pages,
	})
}

func pkg(w http.ResponseWriter, r *http.Request) {
	m := pathPkg.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	log.Println(m)
	pkgname := m[3]
	repo := repos[m[1]]

	if repo == nil {
		http.Redirect(w, r, "/pkg/x86_64/"+pkgname, 302)
		return
	}
	pkg, ok := repo[pkgname]
	if !ok {
		http.NotFound(w, r)
		return
	}

	versions := make(map[string]string)

	for arch, repo := range repos {
		versions[arch] = repo[pkgname].Pkgver
	}

	renderTemplate(w, "pkg.html", pkgData {
		Pkgname: pkgname,
		Versions: versions,
		Pkg: pkg,
	})
}

func main() {
	for name, _ := range repos {
		log.Println("Loading", name + "-repodata")
		var err error
		repos[name], err = pkgdb.DecodeRepoDataFile(name + "-repodata")
		if err != nil {
			log.Fatal(err)
		}
		for pkgname, pkg := range repos[name] {
			ref, ok := pkgindex[pkgname]
			if !ok {
				pkgindex[pkgname] = make(map[string]pkgdb.Package)
				ref = pkgindex[pkgname]
			}
			ref[name] = pkg
		}
	}

	http.HandleFunc("/pkg/", pkg)
	http.HandleFunc("/", index)

	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
