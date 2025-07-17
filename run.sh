#!/usr/bin/env bash

# --- Argument Parsing ---
selected_examples=""
for arg in "$@"; do
  case $arg in
    --examples=*)
      selected_examples="${arg#*=}"
      shift # Remove --examples= from the list of arguments
      ;;
  esac
done

# --- Determine which examples to run ---
examples_to_run=()
if [ -z "$selected_examples" ]; then
  # If no examples are passed, find them dynamically.
  echo "🚀 No examples specified, searching for all in ./docs/examples..."
  all_examples=()
  for filepath in ./docs/examples/example*_input.mecha; do
    # This check handles the case where no files match the glob pattern.
    [ -e "$filepath" ] || continue

    # Extract the filename from the full path (e.g., "example1_input.mecha")
    filename=$(basename "$filepath")

    # Remove the prefix "example" and the suffix "_input.mecha" to get the number.
    num_part="${filename#example}"
    example_num="${num_part%_input.mecha}"

    all_examples+=("$example_num")
  done

  if [ ${#all_examples[@]} -eq 0 ]; then
    echo "⚠️ Error: No example files found in ./docs/examples/ matching the pattern 'example*_input.mecha'."
    exit 1
  fi

  echo "✅ Found examples to run: ${all_examples[*]}"
  examples_to_run=("${all_examples[@]}")

else
  # If examples are passed, split the comma-separated string into an array.
  echo "🚀 Running specified examples: $selected_examples"
  IFS=',' read -r -a examples_to_run <<< "$selected_examples"
fi

# --- Execution ---
# Create the output directory if it doesn't exist.
if [ ! -d "output" ]; then
  echo "Creating output directory..."
  mkdir output
fi

# Loop through the chosen examples and run the Go program.
for i in "${examples_to_run[@]}"; do
  input_file="docs/examples/example${i}_input.mecha"
  output_file="output/output${i}.txt"

  if [ -f "$input_file" ]; then
    echo "-> Compiling example $i: $input_file"
    go run cmd/mecha/main.go -i "$input_file" -o "$output_file"
  else
    echo "-> ⚠️ Warning: Skipping example $i. Input file not found: $input_file"
  fi
done

echo "✅ Done."