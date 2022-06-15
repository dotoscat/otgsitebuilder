package builder

import (
    "testing"
    "os"
    // "path/filepath"
    "flag"
    "log"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

// update-site

const testdataContent = "testdata/content"
var updateSite bool

func updateGoldenSite() {
    content := manager.OpenContent(testdataContent)
    log.Println(content)
    website := NewWebsite(content)
    outputPath := website.Render("testdata")
    log.Println("render path:", outputPath)
}

func TestMain(m *testing.M) {
    flag.BoolVar(&updateSite, "update-site", false, "Re-generate or update golden site.")
    flag.Parse()
    if updateSite == true {
       // log.Println("Update golden site.")
       updateGoldenSite()
       os.Exit(0)
    }
	os.Exit(m.Run())
}
