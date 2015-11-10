#!/bin/bash

F=$1
while read e; do
	  curl -XPOST 'http://192.168.29.130:9200/tracerdemo/event/' -d "${e[@]}"
done < $F
