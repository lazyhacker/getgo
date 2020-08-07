#!/bin/bash                                                                     
                                                                                
#if [ $UID != 0 ]; then                                                         
#    echo "Please run this script with sudo:"                                   
#    echo "sudo $0 $*"                                                          
#    exit 1                                                                     
#fi                                                                             
                                                                                
dl=`getgo -dir ~/Downloads -show true`                                          
                                                                                
sudo rm -rf /usr/local/go                                                       
sudo tar -C /usr/local -xzf $dl
