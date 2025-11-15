#!/bin/bash
# Compare existing proto files with extracted official definitions

echo "=========================================="
echo "PROTO COMPARISON REPORT"
echo "=========================================="
echo ""

for extracted in proto/xai/v1/*.proto.extracted; do
    base=$(basename "$extracted" .proto.extracted)
    existing="proto/xai/v1/${base}.proto"
    
    echo "### $base.proto ###"
    
    if [ -f "$existing" ]; then
        echo "Status: EXISTS"
        echo "Comparing..."
        
        # Extract just message/enum definitions for comparison
        diff -u \
            <(grep -E "^(message|enum|service|  )" "$existing" | grep -v "^  //" || true) \
            <(grep -E "^(message|enum|service|  )" "$extracted" || true) \
            > /tmp/proto_diff_${base}.txt 2>&1
        
        if [ $? -eq 0 ]; then
            echo "✓ MATCHES official Python SDK"
        else
            echo "✗ DIFFERS from official Python SDK"
            echo "  See: /tmp/proto_diff_${base}.txt"
            head -20 /tmp/proto_diff_${base}.txt
        fi
    else
        echo "Status: MISSING"
        echo "✗ Need to create this proto file"
    fi
    
    echo ""
done

echo "=========================================="
echo "SUMMARY"
echo "=========================================="
echo "Extracted files: $(ls proto/xai/v1/*.proto.extracted | wc -l)"
echo "Existing files: $(ls proto/xai/v1/*.proto 2>/dev/null | wc -l)"
echo ""
echo "Review diff files in /tmp/proto_diff_*.txt"
