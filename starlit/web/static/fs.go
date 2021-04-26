package static

import (
	"embed"
)

//go:generate rm -rf fonts images js legacy
//go:generate cp -r ../app/dist/static/* .

//go:embed *
var FS embed.FS
