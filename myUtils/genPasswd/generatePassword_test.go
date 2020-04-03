package genPasswd

import (
	"log"
	"testing"
)

func TestGeneratePasswd(t *testing.T) {
	RandomSeed()
	for i := 0; i < 10; i++ {
		log.Printf("%s\n", GeneratePasswd(8, FlagNumber))
	}

	for i := 0; i < 10; i++ {
		log.Printf("%s\n", GeneratePasswd(10, FlagNumber|FlagLowerChar))
	}

	for i := 0; i < 10; i++ {
		log.Printf("%s\n", GeneratePasswd(10, FlagNumber|FlagLowerChar|FlagUpperChar))
	}

	for i := 0; i < 10; i++ {
		log.Printf("%s\n", GeneratePasswd(10, FlagNumber|FlagLowerChar|FlagUpperChar|FlagSpecialChar))
	}
}
