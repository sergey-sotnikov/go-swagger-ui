go-swagger-ui
============

Alternative to Redoc middleware in https://github.com/go-openapi/runtime package 
for [Swagger UI](https://swagger.io/tools/swagger-ui). This middleware takes latest version of
Swagger UI from CDN and renders page with your endpoints.

## Install
```shell
go get github.com/sergey-sotnikov/go-swagger-ui
```

## Example of usage
```shell
package main 

import "github.com/sergey-sotnikov/go-swagger-ui/swagger"

func main() {
    handler := swagger.Middleware(&swagger.Opts{SpecURL: "/swagger-spec/swagger.json"}, nil)
	if err := http.ListenAndServe(":5300", handler); err != nil {
		panic(err)
	}
}
```