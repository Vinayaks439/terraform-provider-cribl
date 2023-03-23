resource "cribl_local_user" "provideruser" {
  items {
    id = "vinayak"
    username = "providerUser12345"
    first = "customprovider"
    last = "provider"
    email = "provider@gmail.com"
    roles = ["user"]
    password = "password1234asad"
    disabled = false
  }
}