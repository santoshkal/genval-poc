package dockerfile_validation

# # Define policy rules as a list of policies to evaluate
# policy_rules = [
#     untrusted_base_image,
#     latest_base_image,
#     # warning_system_packages_upgrade,
#     any_user,
#     deny_root_user,
#     deny_sudo,
#     # ... other policies ...
# ]



untrusted_base_image{
    input[i].cmd == "from"
    val := split(input[i].value, "/")
    val[0] == "cgr.dev"
}

latest_base_image{
    input[i].cmd == "from"
    val1 := split(input[i].value, ":")
    contains(val1[1], "latest")
}

# # Do not upgrade your system packages
# warning_system_packages_upgrade{
#     input[i].cmd == "run"
#     val2 := concat(" ", input[i].value)
#     matches := regex.match(".*?(apk|yum|dnf|apt|pip).+?(install|[dist-|check-|group]?up[grade|date]).*", val2)
#     matches == true
# }

# Do not use any user
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


   