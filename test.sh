#!/bin/bash

# Pack Calculator - Comprehensive Test Suite
# Run this after: make up (or docker-compose up -d)

BASE_URL="http://localhost"
API_URL="${BASE_URL}/api"

echo "=========================================="
echo "Pack Calculator - Comprehensive Test Suite"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local expected_status=$5
    local expected_content=$6
    
    echo -e "${YELLOW}Testing: ${name}${NC}"
    echo "  Method: $method"
    echo "  URL: $url"
    if [ -n "$data" ]; then
        echo "  Body: $data"
    fi
    echo "  Expected Status: $expected_status"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    echo "  Actual Status: $http_code"
    echo "  Response: $body"
    
    if [ "$http_code" = "$expected_status" ]; then
        if [ -n "$expected_content" ]; then
            if echo "$body" | grep -q "$expected_content"; then
                echo -e "  ${GREEN}✓ PASSED${NC}"
                ((TESTS_PASSED++))
            else
                echo -e "  ${RED}✗ FAILED - Content mismatch${NC}"
                echo "    Expected to contain: $expected_content"
                ((TESTS_FAILED++))
            fi
        else
            echo -e "  ${GREEN}✓ PASSED${NC}"
            ((TESTS_PASSED++))
        fi
    else
        echo -e "  ${RED}✗ FAILED - Status mismatch${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 5

# Test 1: Health Check
echo "=========================================="
echo "1. Health Check Endpoint"
echo "=========================================="
test_endpoint "Health Check" "GET" "${BASE_URL}/health" "" "200" "ok"

# Test 2: Get Pack Sizes (initially empty)
echo "=========================================="
echo "2. Pack Sizes - Get (Initial)"
echo "=========================================="
test_endpoint "Get Pack Sizes (empty)" "GET" "${API_URL}/pack-sizes" "" "200" "sizes"

# Test 3: Update Pack Sizes - Valid
echo "=========================================="
echo "3. Pack Sizes - Update (Valid Cases)"
echo "=========================================="
test_endpoint "Update Pack Sizes [250,500,1000]" "POST" "${API_URL}/pack-sizes" \
    '{"sizes":[250,500,1000]}' "204" ""

test_endpoint "Update Pack Sizes [23,31,53]" "POST" "${API_URL}/pack-sizes" \
    '{"sizes":[23,31,53]}' "204" ""

# Test 4: Get Pack Sizes (after update)
echo "=========================================="
echo "4. Pack Sizes - Get (After Update)"
echo "=========================================="
test_endpoint "Get Pack Sizes [23,31,53]" "GET" "${API_URL}/pack-sizes" "" "200" "23"

# Test 5: Update Pack Sizes - Invalid Cases
echo "=========================================="
echo "5. Pack Sizes - Update (Invalid Cases)"
echo "=========================================="
test_endpoint "Update Pack Sizes (empty array)" "POST" "${API_URL}/pack-sizes" \
    '{"sizes":[]}' "400" ""

test_endpoint "Update Pack Sizes (missing field)" "POST" "${API_URL}/pack-sizes" \
    '{}' "400" ""

test_endpoint "Update Pack Sizes (invalid JSON)" "POST" "${API_URL}/pack-sizes" \
    '{"invalid"}' "400" ""

# Test 6: Calculate Packs - Valid Cases
echo "=========================================="
echo "6. Calculate Packs - Valid Cases"
echo "=========================================="
test_endpoint "Calculate 251 items" "POST" "${API_URL}/calculate" \
    '{"items":251}' "200" "packs"

test_endpoint "Calculate 250 items (exact)" "POST" "${API_URL}/calculate" \
    '{"items":250}' "200" "packs"

test_endpoint "Calculate 1 item" "POST" "${API_URL}/calculate" \
    '{"items":1}' "200" "packs"

test_endpoint "Calculate 12001 items" "POST" "${API_URL}/calculate" \
    '{"items":12001}' "200" "packs"

# Test 7: Edge Case - 500,000 items
echo "=========================================="
echo "7. Edge Case - 500,000 items"
echo "=========================================="
response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"items":500000}' \
    "${API_URL}/calculate")

echo "Response: $response"
echo ""

# Verify edge case result
if echo "$response" | grep -q '"size":23.*"quantity":2'; then
    if echo "$response" | grep -q '"size":31.*"quantity":7'; then
        if echo "$response" | grep -q '"size":53.*"quantity":9429'; then
            echo -e "${GREEN}✓ Edge case PASSED - Correct combination found${NC}"
            ((TESTS_PASSED++))
        else
            echo -e "${RED}✗ Edge case FAILED - Quantity 9429 not found${NC}"
            ((TESTS_FAILED++))
        fi
    else
        echo -e "${RED}✗ Edge case FAILED - Quantity 7 not found${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "${RED}✗ Edge case FAILED - Quantity 2 not found${NC}"
    ((TESTS_FAILED++))
fi
echo ""

# Test 8: Calculate Packs - Invalid Cases
echo "=========================================="
echo "8. Calculate Packs - Invalid Cases"
echo "=========================================="
test_endpoint "Calculate 0 items" "POST" "${API_URL}/calculate" \
    '{"items":0}' "400" ""

test_endpoint "Calculate -1 items" "POST" "${API_URL}/calculate" \
    '{"items":-1}' "400" ""

test_endpoint "Calculate (missing field)" "POST" "${API_URL}/calculate" \
    '{}' "400" ""

test_endpoint "Calculate (invalid JSON)" "POST" "${API_URL}/calculate" \
    '{"invalid"}' "400" ""

# Test 9: Verify Results
echo "=========================================="
echo "9. Verification - Get Pack Sizes"
echo "=========================================="
test_endpoint "Verify pack sizes are saved" "GET" "${API_URL}/pack-sizes" "" "200" "23"

# Test 10: Complete Workflow
echo "=========================================="
echo "10. Complete Workflow Test"
echo "=========================================="
echo "Step 1: Update pack sizes to [250,500,1000]"
test_endpoint "Update to [250,500,1000]" "POST" "${API_URL}/pack-sizes" \
    '{"sizes":[250,500,1000]}' "204" ""

echo "Step 2: Verify pack sizes"
test_endpoint "Get updated sizes" "GET" "${API_URL}/pack-sizes" "" "200" "250"

echo "Step 3: Calculate with new sizes"
test_endpoint "Calculate 251 with new sizes" "POST" "${API_URL}/calculate" \
    '{"items":251}' "200" "packs"

# Test 11: CORS Headers
echo "=========================================="
echo "11. CORS Headers Verification"
echo "=========================================="
cors_response=$(curl -s -I -X OPTIONS \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST" \
    "${API_URL}/pack-sizes")

if echo "$cors_response" | grep -q "Access-Control-Allow-Origin"; then
    echo -e "${GREEN}✓ CORS headers present${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${RED}✗ CORS headers missing${NC}"
    ((TESTS_FAILED++))
fi
echo ""

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Tests Passed: ${TESTS_PASSED}${NC}"
echo -e "${RED}Tests Failed: ${TESTS_FAILED}${NC}"
echo "Total Tests: $((TESTS_PASSED + TESTS_FAILED))"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
