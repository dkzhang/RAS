package genPasswd

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	numStr       = "0123456789"
	lowerCharStr = "abcdefghijklmnopqrstuvwxyz"
	upperCharStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specCharStr  = "+=-@#~,.[]()!%^*$"
)

const (
	FlagNumber      = 0b1
	FlagLowerChar   = 0b10
	FlagUpperChar   = 0b100
	FlagSpecialChar = 0b1000
)

func GeneratePasswd(pwLen int, pwType uint) string {
	//初始化密码切片
	var passwd []byte = make([]byte, pwLen, pwLen)

	//根据密码字符类型合成字符源
	sourceNum := ""
	if pwType&FlagNumber != 0 {
		sourceNum = numStr
	}
	sourceLowerChar := ""
	if pwType&FlagLowerChar != 0 {
		sourceLowerChar = lowerCharStr
	}
	sourceUpperChar := ""
	if pwType&FlagUpperChar != 0 {
		sourceUpperChar = upperCharStr
	}
	sourceSpecialChar := ""
	if pwType&FlagSpecialChar != 0 {
		sourceSpecialChar = specCharStr
	}
	sourceStr := fmt.Sprintf("%s%s%s%s", sourceNum, sourceLowerChar, sourceUpperChar, sourceSpecialChar)

	//遍历，生成一个随机index索引,
	for i := 0; i < pwLen; i++ {
		index := rand.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}
	return string(passwd)
}

func RandomSeed() {
	//随机种子
	rand.Seed(time.Now().UnixNano())
}
