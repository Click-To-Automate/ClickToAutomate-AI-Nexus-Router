package apiport

// Port is the HTTP listen port for the API.
// Dev (wails dev / go run) keeps the default 20128.
// Production builds override this via -ldflags:
//
//	-X ainexusrouter-core/apiport.Port=30128
var Port = "20128"

// ProductionPort is the port baked into production release binaries.
const ProductionPort = "30128"
