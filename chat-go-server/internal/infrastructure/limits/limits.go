package limits

// ------------------
// ------ USER ------
// ------------------

// -- Username
const USERNAME_MAX_LENGTH int = 20
const USERNAME_MIN_LENGTH int = 2

// -- Password
const PASSWORD_MAX_LENGTH int = 60
const PASSWORD_MIN_LENGTH int = 6

var PASSWORD_REQUIRE_SYMBOLS = []string{
	"!?.,:;@#$%^&*()[]{}\\/+-=~`_|",
	"0123456789",
	"ZXCVBNMASDFGHJKLQWERTYUIOP",
}

// ------------------
// ------ CHAT ------
// ------------------

// --

// ------------------
// ------ MESSAGE ------
// ------------------

// --
