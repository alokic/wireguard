#!/bin/bash

#setup trap 
finish () {
    wg-quick down wg0
    exit 0
}
trap finish SIGTERM SIGINT SIGQUIT

#setup wireguard
setup_wg() {
  whoami
  wg-quick up /etc/wireguard/wg0.conf
}

# Execute the CMD from the Dockerfile.
# If you want to run as server, run it as daemon.
run_cmd() {
  exec "$@"
}

# setup open-ssh server
setup_ssh_server() {
  rc-update add sshd boot
  rc-status && touch /run/openrc/softlevel && /usr/sbin/sshd -D -e
}

# setup application services
setup_application_services() {
  set -x
  wireguard &
  set +x
}

# Inifinite sleep
infinite_loop() {
  while true; do
    sleep 86400
    wait $!
  done
}

# setup client
setup_client() {
  if [ -z "${SSH_AUTHORIZED_KEYS}" ]; then
    echo "Client confg: Need your ssh public key as SSH_AUTHORIZED_KEYS env variable. Abnormal exit ..."
    exit 1
  fi

  echo "Client confg: Populating /root/.ssh/authorized_keys with the value from SSH_AUTHORIZED_KEYS env variable ..."
  echo "${SSH_AUTHORIZED_KEYS}" > $HOME/.ssh/authorized_keys

  chmod 0700 $HOME/.ssh
  chmod 0600 $HOME/.ssh/*
}

# setup server
setup_server() {
  if [ -z "${SSH_PUBLIC_KEY}" ]; then
    echo "Server config: Need your ssh public key as PUBLIC_KEY env variable. Abnormal exit ..."
    exit 1
  fi

  echo "Server config: Populating /root/.ssh/id_rsa.pub with the value from PUBLIC_KEY env variable ..."
  echo "${SSH_PUBLIC_KEY}" > $HOME/.ssh/id_rsa.pub

  if [ -z "${SSH_PRIVATE_KEY}" ]; then
    echo "Server config: Need your ssh public key as PRIVATE_KEY env variable. Abnormal exit ..."
    exit 1
  fi

  echo "Server config: Populating /root/.ssh/id_rsa with the value from PRIVATE_KEY env variable ..."
  echo "${SSH_PRIVATE_KEY}" > $HOME/.ssh/id_rsa

  chmod 0700 $HOME/.ssh
  chmod 0600 $HOME/.ssh/*
}

# setup container and bring wireguard interface up.
setup_container() {
  if [[ "${WG_TYPE,,}" != "client" ]] &&  [[ "${WG_TYPE,,}" != "server" ]]; then
    echo "WG_TYPE invalid. Possible values are 'client' or 'server'. Abnormal exit ..."
    exit 1
  fi

  # converting WG_TYPE to client
  if [ "${WG_TYPE,,}" = "client" ]; then 
    setup_client
    setup_server
  fi
  
  # converting WG_TYPE to server
  if [ "${WG_TYPE,,}" = "server" ]; then
    setup_server
  fi

  run_cmd

  setup_application_services
  setup_wg
  setup_ssh_server
  infinite_loop
}

setup_container