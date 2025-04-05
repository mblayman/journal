package entries

import (
	"github.com/mblayman/journal/model"
)

func MakeEmailContentProcessor(requiredToAddress string) model.EmailContentProcessor {
	return func(emailContent model.EmailContent) {
	}
}
