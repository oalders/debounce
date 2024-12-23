#!/usr/bin/env bats

temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT

usage_message="Usage: debounce <quantity> <unit> <command>"

@test "Check version output" {
    result="$(./bin/debounce --version)"
    [ "$result" == "0.5.0" ]
}

@test "Check error message when not enough arguments are provided" {
    run ./bin/debounce
    [ "$status" -eq 1 ]
    [[ "$output" == "$usage_message"* ]]
}

@test "Check success exit code" {
    run ./bin/debounce --cache-dir "$temp_dir" 1 s bash -c 'echo test'
    echo "Output: $output"
    echo "Status: $status"
    [ "$status" -eq 0 ]
}

@test "Check error exit code" {
    run ./bin/debounce 1 s "nonexistentcommand"
    [ "$status" -eq 1 ]
}

@test "Check --status flag with and without cache file" {
    random_number=$(perl -e 'print time')
    command="echo $random_number"

    run ./bin/debounce --debug --cache-dir "$temp_dir" --status 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Cache file does not exist. \"bash -c $command\" will run on next debounce"* ]]

    run ./bin/debounce --debug --cache-dir "$temp_dir" 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"$random_number"* ]]

    run ./bin/debounce --debug --cache-dir "$temp_dir" --status 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"📁 cache location:"* ]]
    [[ "$output" == *"🚧 cache last modified:"* ]]
    [[ "$output" == *"⏲️ debounce interval:"* ]]
    [[ "$output" == *"🕰️ cache age:"* ]]
}

@test "Check non-existent cache directory" {
    non_existent_dir="$temp_dir/non_existent"
    run ./bin/debounce --cache-dir "$non_existent_dir" 1 s bash -c 'echo test'
    [ "$status" -ne 0 ]
    [[ "$output" == *"provided cache directory does not exist"* ]]
}

@test "Check invalid unit error message" {
    run ./bin/debounce 1 invalid_unit echo "test"
    [ "$status" -eq 1 ]
    [[ "$output" == *"<unit> must be one of "* ]]
    [[ "$output" == "$usage_message"* ]]
}

@test "Check invalid quantity error message" {
    run ./bin/debounce invalid_quantity s echo "test"
    [ "$status" -eq 1 ]
    [[ "$output" == *"quantity invalid_quantity is not a valid integer"* ]]
    [[ "$output" == *"$usage_message"* ]]
}

@test "Check all available units" {
    units=("s" "second" "seconds" "m" "minute" "minutes" "h" "hour" "hours" "d" "day" "days")
    for unit in "${units[@]}"; do
        run ./bin/debounce --cache-dir "$temp_dir" 1 "$unit" bash -c 'echo test'
        [ "$status" -eq 0 ]
    done
}

@test "Check --local flag behavior" {
    command="echo local_test"

    # Run the command twice without the --local flag
    run ./bin/debounce --cache-dir "$temp_dir" 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"local_test"* ]]

    run ./bin/debounce --cache-dir "$temp_dir" 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"🚥 will not run"* ]]

    # Run the command twice with the --local flag
    run ./bin/debounce --cache-dir "$temp_dir" --local 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"local_test"* ]]

    run ./bin/debounce --cache-dir "$temp_dir" --local 10 s bash -c "$command"
    [ "$status" -eq 0 ]
    [[ "$output" == *"🚥 will not run"* ]]
}

@test "Check command that exits with error code" {
    run ./bin/debounce 1 s bash -c "exit 2"
    [ "$status" -eq 1 ]
    [[ "$output" == *"running command: bash -c exit 2"* ]]
    [[ "$output" != *"$usage_message"* ]]
}
