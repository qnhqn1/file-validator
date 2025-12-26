package swagger

import _ "embed" // required to enable //go:embed for embedding the swagger JSON


var validatorSpec []byte


func Validator() []byte {
	return validatorSpec
}


