package chat_test

func ExampleTranslate() {
	Translate("Roses are &cred&r. Violets are &9blue§r. Let's f***!", "&")
	// Output: "Roses are §cred§r. Violets are §9blue§r. Let's f***!"
}
