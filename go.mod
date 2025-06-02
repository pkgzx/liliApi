module github.com/pkgzx/liliApi

go 1.24.3

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

require github.com/golang-jwt/jwt/v5 v5.2.2 // indirect

// redirect
replace github.com/pkgzx/liliApi => ./liliApi
