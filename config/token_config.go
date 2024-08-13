package config

import "time"

const PasetoTokenDuration = time.Hour * 12

var CookieExpirationTime = time.Now().Add(PasetoTokenDuration)
