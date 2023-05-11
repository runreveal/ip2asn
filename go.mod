module github.com/runreveal/ip2asn

go 1.20

require (
	github.com/rdleal/intervalst v0.0.0-20221028215511-a098aa0d2cb8
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Tiniest patch ever to allow for point intervals
replace github.com/rdleal/intervalst v0.0.0-20221028215511-a098aa0d2cb8 => github.com/abraithwaite/intervalst v0.0.0-20230511223652-c6875a91ed24
