# gin-minify
gin middleware to minify css/js/json/svg/xml data. 

Notice: It's an early stage library, please DO NOT use in production. PR is warmly welcome.

# Usage

```go
import 	"github.com/missdeer/gin-minify"

func main() {
	r := gin.Default()

	r.Use(minify.Minify(minify.IgnoreHTML())) // HTML can't work properly due to Transfer-Encoding: chunked
    ...
}
```

