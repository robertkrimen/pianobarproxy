#/bin/bash
#-----------------------------------------------------------
# Licensed separately from pianobarproxy
# Licensed as BSD 2-clause License
# Copyright (c) 2016, Brendan Horan
# Description :
# A wrapper script that starts/stops
# an amazon EC2 instance in the USA
# Then runs pianobarproxy so that
# you can use pianobar outside of USA.
# Depends on :
# * pianobarproxy
# * pianobar (https://6xq.net/pianobar)
# * euca2ools
# * Any Linux VM on Amazon EC2 (In USA)
#   https://github.com/eucalyptus/euca2ools
# Usage :
# ./piano-proxy.sh {start|stop}
#-----------------------------------------------------------


# Fill in the details for the variables

# Environment variables for euca commands
# Currently only has the following line :
# export EC2_URL=http://ec2.amazonaws.com
euca2ools_env_vars="euca2ools.conf"

# Path to the pianobarproxy binary
pianobarproxy_path=""

# Amazon key
# user attached to this key must have rights 
# to start/stop VM's
amazon_key=""

# Amazon secret key
amazon_secret=""

# Instance ID of VM to start/stop
amazon_instance_id=""

# Pem file from amazon
instance_ssh_key_file=""

# EC2 image username
user_name=""

# End variables section


cmd_args=$1

start(){
  source $euca2ools_env_vars

  euca-start-instances -I $amazon_key -S $amazon_secret $amazon_instance_id &
  wait
  
  sleep 30

  instance_dns_name=`euca-describe-instances -I $amazon_key -S $amazon_secret \
    $amazon_instance_id | grep publicDnsName | awk -F, {'print $3'} | awk -F\' {'print $4'}`

  nohup ssh -o StrictHostKeyChecking=no -TMNf -D 9050 $user_name@$instance_dns_name \
    -i $instance_ssh_key_file >/dev/null 2>&1 </dev/null &

  sleep 15

   $pianobarproxy_path -socks5 :9050 &
}


stop() {

  source $euca2ools_env_vars

  euca-stop-instances -I $amazon_key -S $amazon_secret $amazon_instance_id

  pid_of_pianobarproxy=`ps -ef | grep '[p]ianobarproxy' | awk -F" " {'print $2'}`

  kill -9 $pid_of_pianobarproxy
}

case "$cmd_args" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  *)
    echo $"Usage: $0 {start|stop}"
    exit 1
esac
