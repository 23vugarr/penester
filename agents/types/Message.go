package types

import (
	"masProject/penester/pkg"
)

type Message struct {
	Type    string
	Message string
	Content pkg.Config
}
