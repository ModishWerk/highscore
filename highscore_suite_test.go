package main_test

import (
	// . "github.com/emicklei/forest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

// var host = NewClient("http://localhost:"+strconv.Itoa(Port), new(http.Client))

func TestHighscore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Highscore Suite")
}
