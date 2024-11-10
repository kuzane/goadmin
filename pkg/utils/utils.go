package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	once       sync.Once
	randSource rand.Source
)

func init() {
	once.Do(func() { // 使用 sync.Once 确保初始化种子只进行一次
		randSource = rand.NewSource(time.Now().UnixNano())
	})
}

// GenerateCaptcha 生成6为随机验证码
func GenerateCaptcha() string {
	r := rand.New(randSource)
	code := r.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// GeneratePassword 生成6位的随机密码，包含大小写字母、数字和特殊字符
func GeneratePassword() string {
	randGen := rand.New(randSource)
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	specialChars := "@#$%&"
	// 将字符集合合并
	charset := lowerLetters + upperLetters + numbers + specialChars
	password := make([]byte, 6)
	// 保证至少包含一个小写字母、大写字母、数字和特殊字符
	password[0] = lowerLetters[randGen.Intn(len(lowerLetters))]
	password[1] = upperLetters[randGen.Intn(len(upperLetters))]
	password[2] = numbers[randGen.Intn(len(numbers))]
	password[3] = specialChars[randGen.Intn(len(specialChars))]
	// 剩余两位从所有字符中随机选择
	for i := 4; i < 6; i++ {
		password[i] = charset[randGen.Intn(len(charset))]
	}
	// 打乱密码字符顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}
