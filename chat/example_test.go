package chat_test

import (
	"fmt"

	"github.com/minero/minero/chat"
)

func ExampleTranslate() {
	t := chat.Translate("Roses are &cred&r. Violets are &9blue§r. Let's f***!", "&")
	fmt.Println(t)
	// Output: Roses are §cred§r. Violets are §9blue§r. Let's f***!
}
