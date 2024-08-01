#!/usr/bin/env bats

@test "Check version output" {
  result="$(./debounce --version)"
  [ "$result" == "Version: 0.1.0" ]
}

@test "Check error message when not enough arguments are provided" {
  run ./debounce
  [ "$status" -eq 1 ]
  [ "$output" == "Please provide the quantity, unit, and command" ]
}

@test "Check success exit code" {
  run ./debounce 1 s "bash -c" echo test
  echo "Output: $output"
  echo "Status: $status"
  [ "$status" -eq 0 ]
}

@test "Check error exit code" {
  run ./debounce 1 s "nonexistentcommand"
  [ "$status" -eq 1 ]
}
