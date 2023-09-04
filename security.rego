package dockerfile_validation



untrusted_base_image{
    input[i].cmd == "from"
    val := split(input[i].value, "/")
    val[0] == "cgr.dev"
}

# latest_base_image{
#     input[i].cmd == "from"
#     val1 := split(input[i].value, ":")
#     contains(val1[1], "latest")
# }


# Do not use root user
deny_root_user {
    input[i].cmd == "user"
    val2:= input[i].value
    val2 != "root"
    val2 != "0" 
}

# Do not sudo
deny_sudo{
    input[i].cmd == "run"
    val3:= input[i].value
    not contains(val3, "sudo")
}

# Avoid using cached layers CIS 4.7
deny_caching{
    input[i].cmd == "run"
    val4:= input[i].value
    matches := regex.match(".*?(apk|yum|dnf|apt|pip).+?(install|[dist-|check-|group]?up[grade|date]).*", val4)
    matches == true
    contains(val4, "--no-cache")
}

# Ensure that COPY is used instead of ADD CIS 4.9
deny_add{
    input[i].cmd != "add"
}

# Ensure ADD does not include unpack archives or download files 
deny_image_expansion{
	input[_].cmd == "add"
	val5 := input[_].value
	words := regex.match(".*?(curl|wget|.tar|.tar.).*", val5)
	words != true
}

# Ensure secrets are not stored CIS 4.10

# deny_root_user{
# forbidden_users := [
#     "root",
#     "toor",
#     "0"
# ]
#     # command := "user"
#     users := [name | input[i].cmd == "user"; name := input[i].value]
#     lastuser := users[count(users)-1]
#     contains(lower(lastuser[j]), forbidden_users[k])
# }
# https://play.openpolicyagent.org/p/x9pwuimWKd


   