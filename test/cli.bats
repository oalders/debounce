#!/usr/bin/env bats

@test "Check version output" {
    result="$(./bin/debounce --version)"
    [ "$result" == "debounce 0.2.0" ]
}

@test "Check error message when not enough arguments are provided" {
    run ./bin/debounce
    [ "$status" -eq 1 ]
    [[ "$output" == "Usage: debounce <integer> <hours|minutes|seconds> <command>"* ]]
}

@test "Check success exit code" {
    run ./bin/debounce 1 s bash -c 'echo test'
    echo "Output: $output"
    echo "Status: $status"
    [ "$status" -eq 0 ]
}

@test "Check error exit code" {
    run ./bin/debounce 1 s "nonexistentcommand"
    [ "$status" -eq 1 ]
}
