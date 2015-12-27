:: wget for windows can be downloaded at http://nebm.ist.utl.pt/~glopes/wget/
@echo off
wget http://10.0.0.7/steamtest/serverdump.json
mkdir dump
move serverdump.json dump/