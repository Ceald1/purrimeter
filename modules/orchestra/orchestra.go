package orchestra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Ceald1/purrimeter/api/crypto"

	YAML "github.com/goccy/go-yaml"

	"github.com/gin-gonic/gin"
	surrealdb "github.com/surrealdb/surrealdb.go"
	EnrichmentTypes "github.com/Ceald1/purrimeter/modules/enrichment/types"
	AlertTypes "github.com/Ceald1/purrimeter/modules/alerts/types"
	types "github.com/Ceald1/purrimeter/modules/orchestra/types"
)
var (
  SURREAL_ADMIN string = os.Getenv("SURREAL_ADMIN")
  SURREAL_PASS string = os.Getenv("SURREAL_PASS")
  ctx = context.Background()
)
