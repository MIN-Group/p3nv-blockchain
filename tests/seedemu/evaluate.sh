#!/bin/bash

script_path="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "$script_path"

dimension=6
count=4

rm -rf result/dimension_$dimension/
mkdir -p result/dimension_$dimension/

output=$(sudo docker ps -a)

for ((i = 0; i < 2 ** dimension; i++)); do
    binary=$(echo "obase=2; $i" | bc)
    binary=$(printf "%0${dimension}d\n" "$binary")

    regex="([0-9a-f]{12})\s+output_rnode_150_node$binary"
    if [[ $output =~ $regex ]]; then
        container_id="${BASH_REMATCH[1]}"
    fi

    echo "collecting results for node$binary corresponding to container $container_id"
    sudo docker cp $container_id:/app/consensus.csv result/dimension_$dimension/consensus_$i.csv
    sudo docker cp $container_id:/app/batch.csv result/dimension_$dimension/batch_$i.csv
    sudo docker cp $container_id:/app/output.log result/dimension_$dimension/output_$i.log

    ((count--))
    if [ $count -eq 0 ]; then
        break
    fi
done
