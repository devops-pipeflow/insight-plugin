#!/bin/bash

# Set defaults
FIX_ME="true"
SILENT_MODE="false"
VERSION_INFO="1.2.0"

# Print pass message
print_pass() {
    if [ "$SILENT_MODE" = "false" ]; then
        echo -e "$1"
    fi
}

# Print fail message
print_fail() {
    echo -e "$1" 1>&2
}

# Print fix message
print_fix() {
    echo -e "$1" 1>&2
}

# Check /etc/hosts
check_hosts() {
    local pattern="localhost"
    local name="/etc/hosts"
    local cmd="grep -i $pattern $name"
    local ret

    eval "$cmd 1> /dev/null"
    ret=$?
    if [ $ret -ne 0 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($pattern found)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: remove $pattern from $name"
    fi

    return 1
}

# Check /etc/network/interfaces
check_interfaces() {
    local name="/etc/network/interfaces"
    local server1="8.8.8.8"
    local server2="8.8.4.4"
    local cmd="grep -E '$server1|$server2' $name"
    local ret

    eval "$cmd 1> /dev/null"
    ret=$?
    if [ $ret -eq 0 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($server1, $server2 required)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo bash -c \"cat >> $name\" << EOF"
        print_fix "dns-nameservers $server1"
        print_fix "dns-nameservers $server2"
        print_fix "EOF"
    fi

    return 2
}

# Check /etc/resolv.conf
check_resolv() {
    local name="/etc/resolv.conf"
    local server1="8.8.8.8"
    local server2="8.8.4.4"
    local cmd="grep -E '$server1|$server2' $name"
    local ret

    eval "$cmd 1> /dev/null"
    ret=$?
    if [ $ret -eq 0 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($server1, $server2 required)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo bash -c \"cat >> $name\" << EOF"
        print_fix "nameserver $server1"
        print_fix "nameserver $server2"
        print_fix "search localhost"
        print_fix "EOF"
    fi

    return 3
}

# Check /etc/sysctl.conf
check_sysctl() {
    local name="/etc/sysctl.conf"
    local config="net.ipv4.ip_forward"
    local cmd="grep -i $config $name"
    local arr
    local ret

    ret=$(eval "$cmd | sed '/#.*/d'")
    arr=(${ret//=/ })
    if [ "${arr[-2]}" = $config ] && [ "${arr[-1]}" -eq 1 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($config=1 required)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo bash -c \"cat >> /etc/sysctl.conf\" << EOF"
        print_fix "net.ipv4.ip_forward=1"
        print_fix "EOF"
        print_fix "sudo sysctl -p"
    fi

    return 4
}

# Check Docker
check_docker() {
    local name="docker"
    local cmd="$name version"
    local ret

    eval "$cmd 1> /dev/null"
    ret=$?
    if [ $ret -eq 0 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($name missing)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo apt autoremove docker docker.io"
        print_fix "sudo apt update && sudo apt install -y docker docker.io"
        print_fix "sudo service docker restart"
        print_fix "sudo chmod 666 /var/run/docker.sock"
    fi

    return 5
}

# Check /var/run/docker.sock
check_group() {
    local name="docker"
    local cmd="groups $USER | grep -q $name"
    local ret

    eval "$cmd 1> /dev/null"
    ret=$?
    if [ $ret -eq 0 ]; then
        print_pass "\e[1mINFO\e[0m: check user:$USER in group:$name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check user:$USER in group:$name \e[91mFAIL\e[0m ($USER missing)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: sudo usermod -a -G docker $USER"
    fi

    return 6
}

# Check /etc/default/docker
check_default() {
    local name="/etc/default/docker"
    local config="DOCKER_OPTS"
    local value="\"--insecure-registry 0.0.0.0/0\""
    local cmd="grep -i $config $name"
    local arr
    local ret

    ret=$(eval "$cmd | sed '/#.*/d'")
    arr=(${ret//=/ })
    if [ "${arr[0]}" = "$config" ] && [ "${arr[1]} ${arr[2]}" = "$value" ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ($config=$value required)"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "echo \"DOCKER_OPTS=\\\"--insecure-registry 0.0.0.0/0\\\"\" | sudo tee --append /etc/default/docker"
        print_fix "sudo service docker restart"
    fi

    return 7
}

# Check /etc/docker/daemon.json
check_daemon() {
    local name="/etc/docker/daemon.json"
    local config_registries="\"insecure-registries\""
    local value_registries="\"0.0.0.0/0\""
    local err_registries="false"
    local arr
    local ret

    ret=$(eval "grep -i $config_registries $name")
    if [ -z "${ret}" ]; then
        err_registries="true"
    fi

    ret=$(eval "grep -i $value_registries $name")
    if [ -z "${ret}" ]; then
        err_registries="true"
    fi

    if [ "$err_registries" = "false" ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    if [ "$err_registries" = "true" ]; then
        print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m ({$config_registries:[$value_registries]} required)"
    fi

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo touch /etc/docker/daemon.json"
        print_fix "echo \"{\\\"insecure-registries\\\": [ \\\"0.0.0.0/0\\\" ]}\" | sudo tee /etc/docker/daemon.json"
        print_fix "sudo service docker restart"
    fi

    return 8
}

# Check netstat
check_netstat() {
    local name="netstat"
    local cmd="$name --version 2> /dev/null"
    local ret

    _=$(eval "$cmd")
    ret=$?
    if [ "$ret" -ne 127 ]; then
        print_pass "\e[1mINFO\e[0m: check $name \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check $name \e[91mFAIL\e[0m"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "sudo apt install net-tools"
    fi

    return 9
}

# Check $HOME/.ssh/
check_ssh() {
    local addr

    addr=$(ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v inet6 | grep 10. | awk '{print $2}' | tr -d "addr:"​)
    if [[ "$addr" == "10."* ]]; then
        host=$addr
    else
        host="127.0.0.1"
    fi

    print_pass "\e[1mNOTE\e[0m: check ssh host \e[93m$USER@$host:22\e[0m"
    print_pass "\e[1mNOTE\e[0m: check ssh-keygen \e[93mPress <ENTER> when 'Enter passphrase' asked\e[0m"

    return 10
}

# Check disk
check_disk() {
    local free
    local limit=100

    free=$(df -m --output=avail "/boot" | tail -n1)
    if [[ $free -gt $limit ]]; then
        print_pass "\e[1mINFO\e[0m: check /boot \e[92mPASS\e[0m"
        return 0
    fi

    print_fail "\e[1mERROR\e[0m: check /boot \e[91mFAIL\e[0m"

    if [ "$FIX_ME" = "true" ]; then
        print_fix "\e[1mFIXME\e[0m: copy command below, then run it"
        print_fix "export RELEASE=$(uname -r)"
        echo "sudo find /boot -type f -name \"config-*\" -o -name \"initrd.img-*\" -o -name \"System.map-*\" -o -name \"vmlinuz-*\" | grep -vE \"\$RELEASE\$\" | xargs sudo rm -rf"
    fi

    return 11
}

run_check() {
    local ret

    print_pass "----------------------------------------------"

    check_hosts
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_interfaces
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_resolv
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_sysctl
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_docker
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_group
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_default
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_daemon
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_netstat
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_ssh
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    print_pass "----------------------------------------------"

    check_disk
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    return 0
}

show_usage() {
cat <<USAGE

Usage:
    bash $0 [OPTIONS]

Description:
    Perform Ubuntu server health check

OPTIONS:
    -h, --help
        Display this help message

    -s, --silent
        Show error message only

    -v, --version
        Display version information
USAGE
}

set_silent() {
    SILENT_MODE="true"
}

print_version() {
cat <<VERSION
$VERSION_INFO
VERSION
}

parse_opts() {
    local long_opts="help,silent,version,"
    local short_opts="hsv"
    local getopt_cmd

    getopt_cmd=$(getopt -o $short_opts --long "$long_opts" \
                -n "$(basename "$0")" -- "$@") || \
                { show_usage; return 1; }

    eval set -- "$getopt_cmd"

    while true; do
        case "$1" in
            -h|--help) show_usage; return 1;;
            -s|--silent) set_silent;;
            -v|--version) print_version; return 1;;
            --) shift; break;;
        esac
        shift
    done

    return 0
}

main() {
    local ret

    parse_opts "$@"
    ret=$?
    if [ $ret -ne 0 ]; then
        return 0
    fi

    run_check
    ret=$?
    if [ $ret -ne 0 ]; then
        return $ret
    fi

    return 0
}

main "$@"
ret=$?

exit $ret
