module github.com/ingotmc/ingot

go 1.15

replace (
	github.com/ingotmc/mc => ../mc
	github.com/ingotmc/protocol => ../protocol
	github.com/ingotmc/worldgen => ../worldgen
)

require (
	github.com/ingotmc/mc v0.0.1
	github.com/ingotmc/protocol v0.0.1
	github.com/ingotmc/worldgen v0.0.0-00010101000000-000000000000
	github.com/ojrac/opensimplex-go v1.0.1
)
