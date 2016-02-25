#/bin/bash
#YOU NEED SOME FACEBOOK SECRETS EXPORTS TO RUN THIS
#FACEBOOK_APP_ID, FACEBOOK_SECRET
#CONSIDER BLACKBOX OR SOME FILECRYPTING APPS
export DB_PORT_28015_TCP_ADDR=192.168.99.100
export DB_PORT_28015_TCP_PORT=28015
export DNS_HOSTNAME="localhost:8000"
export SESSION_SECRET="session_secret"

if [ $# -gt 0 ] ; then
   p=""
   for var in "$@" ; do
     if [[  $p != "" ]] ; then
        p="$p,$var"
     else
	p=$var
     fi
   done
   godebug run -instrument=$p *.go

else
   go run *.go
fi
