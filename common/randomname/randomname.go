package randomname

import (
	"math/rand"
	"time"
)

// GenerateName 生成随机昵称
func GenerateName() string {
	var name string
	rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedType := RandomType(rand.Intn(2))
	switch selectedType {
	case AdjectiveAndPerson:
		name = AdjectiveSlice[rand.Intn(AdjectiveSliceCount)] + PersonSlice[rand.Intn(PersonSliceCount)]
	case PersonActSomething:
		name = PersonSlice[rand.Intn(PersonSliceCount)] + ActSomethingSlice[rand.Intn(ActSomethingSliceCount)]
	}
	return name
}
