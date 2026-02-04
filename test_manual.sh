#!/bin/bash

# Manual curl commands for testing
# Copy and paste these commands one by one

BASE_URL="http://localhost"
API_URL="${BASE_URL}/api"

echo "=========================================="
echo "Manual Test Commands"
echo "=========================================="
echo ""

echo "1. Health Check:"
echo "curl -X GET ${BASE_URL}/health"
echo ""

echo "2. Get Pack Sizes (initially empty):"
echo "curl -X GET ${API_URL}/pack-sizes"
echo ""

echo "3. Update Pack Sizes [250,500,1000]:"
echo "curl -X POST ${API_URL}/pack-sizes -H 'Content-Type: application/json' -d '{\"sizes\":[250,500,1000]}'"
echo ""

echo "4. Get Pack Sizes (after update):"
echo "curl -X GET ${API_URL}/pack-sizes"
echo ""

echo "5. Update Pack Sizes [23,31,53] (for edge case):"
echo "curl -X POST ${API_URL}/pack-sizes -H 'Content-Type: application/json' -d '{\"sizes\":[23,31,53]}'"
echo ""

echo "6. Calculate 251 items:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":251}'"
echo ""

echo "7. Calculate 250 items (exact match):"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":250}'"
echo ""

echo "8. Calculate 1 item:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":1}'"
echo ""

echo "9. Calculate 12001 items:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":12001}'"
echo ""

echo "10. Edge Case - Calculate 500,000 items:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":500000}'"
echo "Expected: {\"packs\":[{\"size\":23,\"quantity\":2},{\"size\":31,\"quantity\":7},{\"size\":53,\"quantity\":9429}]}"
echo ""

echo "11. Invalid - Empty pack sizes:"
echo "curl -X POST ${API_URL}/pack-sizes -H 'Content-Type: application/json' -d '{\"sizes\":[]}'"
echo ""

echo "12. Invalid - Zero items:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":0}'"
echo ""

echo "13. Invalid - Negative items:"
echo "curl -X POST ${API_URL}/calculate -H 'Content-Type: application/json' -d '{\"items\":-1}'"
echo ""

echo "14. CORS Preflight Check:"
echo "curl -X OPTIONS ${API_URL}/pack-sizes -H 'Origin: http://localhost:3000' -H 'Access-Control-Request-Method: POST' -v"
echo ""
