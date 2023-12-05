#!/bin/bash

########## Change to the Path where the Script Is Located ##########
script_path="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "$script_path"
rm -rf ./workdir/benchmarks/

########## Some Functions for Modifying Global Variables ##########
#tests/main.go#NodeCount
function set_node_count() {
  echo "set_node_count $1"
  sed -E -i 's/NodeCount\s*=\s*[0-9]+/NodeCount='$1'/' ./main.go
}

########## Experiment 1: Basic Performance ##########
#ExecuteTxFlag = false
#PreserveTxFlag = true
#BlockTxLimit = 20000
#LoadJobPerTick = 500
#LoadSubmitNodes = []int{0}
#CheckRotation = false
#BroadcastTx = false
#BenchLoads = []int{5000}
function run_experiment_basic() {
  mkdir -p ./workdir/experiment-pov/
  >./workdir/experiment-pov.log

  echo "> starting experiment 1: PoV"
  for i in {4..28..2}; do
    set_node_count "$i"
    go run . >>./workdir/experiment-pov.log 2>&1
  done
  mv ./workdir/benchmarks/* ./workdir/experiment-pov/
  echo -e "> finished experiment 1: PoV\n"
}
run_experiment_basic