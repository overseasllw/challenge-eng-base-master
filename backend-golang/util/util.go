package chat

import "math/rand"

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	uniqueCheck := map[int]bool{}
	b := make([]rune, n)
	for i := range b {
		randId := rand.Intn(len(letter))
		if _, ok := uniqueCheck[randId]; ok {
			randId = rand.Intn(len(letter))
		}

		b[i] = letter[randId]
	}
	return string(b)
}
