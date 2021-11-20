package constant

import (
	"fmt"
	"math/rand"
)

func GenerateName() string {
	prefix := []string{"anjing", "monyet", "gajah", "kupu-kupu"}
	suffix := []string{"marah", "senang", "sedih", "galau", "bahagia"}
	randomIndexPrefix := rand.Intn(len(prefix))
	randomIndexSuffix := rand.Intn(len(suffix))
	pickPrefix := prefix[randomIndexPrefix]
	pickSuffix := suffix[randomIndexSuffix]
	return fmt.Sprintf(pickPrefix + " " + pickSuffix)
}
