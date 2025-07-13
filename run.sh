#!/usr/bin/env bash

mkdir output

# run the project
go run cmd/mecha/main.go -i docs/examples/example1_input.mecha -o output/outpu1.txt
go run cmd/mecha/main.go -i docs/examples/example2_input.mecha -o output/outpu2.txt
go run cmd/mecha/main.go -i docs/examples/example3_input.mecha -o output/outpu3.txt
go run cmd/mecha/main.go -i docs/examples/example4_input.mecha -o output/outpu4.txt