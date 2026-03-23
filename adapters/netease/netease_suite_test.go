package netease

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNetease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Netease Suite")
}
