package authentication

import (
    "kademlia"
)

type User struct {
    username string
    userID kademlia.ID
}

type UserGroup struct {
    name string
    access map[User]int
}