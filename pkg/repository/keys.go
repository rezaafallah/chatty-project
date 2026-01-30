package repository

import "fmt"

//(Key Builders)

// HistoryKey
func HistoryKey(userID string) string {
	return fmt.Sprintf("chat:history:%s", userID)
}

// OnlineKey
func OnlineKey(userID string) string {
	return fmt.Sprintf("user:online:%s", userID)
}