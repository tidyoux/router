#!/bin/bash

main_folder=${1}
target_name=${2}

go install ../cmd/${main_folder}
go build -o ./bin/${target_name} ../cmd/${main_folder}

echo