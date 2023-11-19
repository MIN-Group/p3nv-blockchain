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
#BatchTxLimit = 5000
#BlockBatchLimit: -1
#VoteBatchLimit:  -1
#LoadJobPerTick = 500
#LoadSubmitNodes = []int{0}
#CheckRotation = false
#BroadcastTx = false
#BenchDuration = 1 * time.Minute
#BenchLoads = []int{20000}
function run_experiment_basic() {
  mkdir -p ./workdir/experiment-ppov/
  >./workdir/experiment-ppov.log

  echo "> starting experiment 1: PPoV"
  for i in {4..28..4}; do
    set_node_count "$i"
    go run . >>./workdir/experiment-ppov.log 2>&1
  done
  mv ./workdir/benchmarks/* ./workdir/experiment-ppov/
  echo -e "> finished experiment 1: PPoV\n"
}
run_experiment_basic