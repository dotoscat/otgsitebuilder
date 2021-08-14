package builder

import (
    "fmt"
    "embed"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed templates/base.tmpl
//go:embed templates/postspage.tmpl
var basicTemplates embed.FS

func Build(base string) {
    //to output
    content := manager.OpenContent(base)
    fmt.Println(content)

}
