#!/bin/bash
# Split CSV files >100MB into ~45MB parts, preserving headers
set -euo pipefail

TARGET_MB=45

split_csv() {
    local file="$1"
    local dir
    dir=$(dirname "$file")
    local base
    base=$(basename "$file" .csv)
    local target_bytes=$((TARGET_MB * 1024 * 1024))

    local size
    size=$(stat -c%s "$file")
    local lines
    lines=$(wc -l < "$file")
    local data_lines=$((lines - 1))

    local num_parts=$(( (size + target_bytes - 1) / target_bytes ))
    local lines_per_part=$(( (data_lines + num_parts - 1) / num_parts ))
    [ "$lines_per_part" -lt 1 ] && lines_per_part=1

    echo "    Size: $(numfmt --to=iec $size) | Lines: $lines | Parts: $num_parts (~${lines_per_part} lines/part)"

    local header
    header=$(head -1 "$file")

    tail -n +2 "$file" | split -l "$lines_per_part" - "${dir}/${base}_part"

    for part in "${dir}/${base}_part"*; do
        mv "$part" "${part}.csv"
        { echo "$header"; cat "${part}.csv"; } > "${part}.tmp" && mv "${part}.tmp" "${part}.csv"
    done

    echo "    Created $(ls "${dir}/${base}_part"*.csv | wc -l) parts"
}

echo "=== Finding CSV files > 100MB ==="
files=$(find dataCSV/ -name "*.csv" -size +100M | sort)
count=$(echo "$files" | wc -l)
echo "Found $count files"
echo ""

total_parts=0
i=0
while IFS= read -r file; do
    [ -z "$file" ] && break
    i=$((i + 1))
    echo "[$i/$count] $file"
    split_csv "$file"
    total_parts=$((total_parts + $(ls "${file%.csv}_part"*.csv 2>/dev/null | wc -l)))
    echo ""
done <<< "$files"

echo "=== Done! Split $count files into $total_parts parts ==="
