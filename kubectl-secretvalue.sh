#!/bin/bash

if [ "$1" = "" ] || [ "$&" = "-h" ] || [ "$&" = "-help" ] 
then
    echo "help"
    echo ""
else
    secret_name=$1
    secret_field=$2

    if [ "$secret_field" = "" ] || [ "$secret_field" = "-h" ] || [ "$secret_field" = "-help" ] 
    then
        echo " - secret field: string"
        echo " - namespace: -n / default"
        echo ""
        exit 1
    fi

    if [ "$secret_field" = "-n" ]
    then
        echo "missing secret field"
        echo ""
        exit 1
    fi

    if [ "$3" = "" ]
    then
        # default namespace
        kubectl get secret $secret_name -o jsonpath="{.data."$secret_field"}" -n default | base64 -d
        echo ""
    elif [ "$3" = "-n" ]
    then
        # get specified namespace
        kubectl get secret $secret_name -o jsonpath="{.data."$secret_field"}" -n $4 | base64 -d
        echo ""
    else 
        echo "Namespace format incomplite"
        echo ""
    fi
fi
