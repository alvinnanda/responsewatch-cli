#!/bin/bash

# Test script for ResponseWatch CLI
# Run this to verify all commands work correctly

set -e

CLI="./rwcli"
FAILED=0
PASSED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

run_test() {
    local name="$1"
    local cmd="$2"
    
    echo -n "Testing: $name ... "
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        ((FAILED++))
    fi
}

echo "====================================="
echo "ResponseWatch CLI Test Suite"
echo "====================================="
echo ""

# Build first
echo "Building CLI..."
go build -o rwcli . > /dev/null 2>&1
echo ""

# Test 1: Basic commands
echo "--- Basic Commands ---"
run_test "version" "$CLI version"
run_test "help" "$CLI --help"
run_test "version with json" "$CLI version -o json"
run_test "version no-color" "$CLI version --no-color"

echo ""
echo "--- Authentication Commands (Help Only) ---"
run_test "login --help" "$CLI login --help"
run_test "logout --help" "$CLI logout --help"
run_test "me --help" "$CLI me --help"
run_test "profile --help" "$CLI profile --help"
run_test "password --help" "$CLI password --help"

echo ""
echo "--- Request Commands (Help Only) ---"
run_test "request --help" "$CLI request --help"
run_test "request list --help" "$CLI request list --help"
run_test "request create --help" "$CLI request create --help"
run_test "request get --help" "$CLI request get --help"
run_test "request update --help" "$CLI request update --help"
run_test "request delete --help" "$CLI request delete --help"
run_test "request reopen --help" "$CLI request reopen --help"
run_test "request assign --help" "$CLI request assign --help"
run_test "request stats --help" "$CLI request stats --help"
run_test "request export --help" "$CLI request export --help"
run_test "request start --help" "$CLI request start --help"
run_test "request finish --help" "$CLI request finish --help"

echo ""
echo "--- Group Commands (Help Only) ---"
run_test "group --help" "$CLI group --help"
run_test "group list --help" "$CLI group list --help"
run_test "group create --help" "$CLI group create --help"
run_test "group get --help" "$CLI group get --help"
run_test "group update --help" "$CLI group update --help"
run_test "group delete --help" "$CLI group delete --help"

echo ""
echo "--- Monitor Commands (Help Only) ---"
run_test "monitor --help" "$CLI monitor --help"
run_test "monitor public --help" "$CLI monitor public --help"

echo ""
echo "--- Note Commands (Help Only) ---"
run_test "note --help" "$CLI note --help"
run_test "note list --help" "$CLI note list --help"
run_test "note create --help" "$CLI note create --help"
run_test "note update --help" "$CLI note update --help"
run_test "note delete --help" "$CLI note delete --help"
run_test "note reminders --help" "$CLI note reminders --help"

echo ""
echo "--- Notification Commands (Help Only) ---"
run_test "notif --help" "$CLI notif --help"
run_test "notif list --help" "$CLI notif list --help"
run_test "notif unread --help" "$CLI notif unread --help"
run_test "notif read --help" "$CLI notif read --help"
run_test "notif read-all --help" "$CLI notif read-all --help"

echo ""
echo "--- Admin Commands (Help Only) ---"
run_test "admin --help" "$CLI admin --help"
run_test "admin users --help" "$CLI admin users --help"
run_test "admin upgrade --help" "$CLI admin upgrade --help"

echo ""
echo "--- Auth Required Commands (Should Fail Without Auth) ---"
# These should fail because we're not logged in
if $CLI me 2>&1 | grep -q "not authenticated"; then
    echo -e "Test: me without auth ... ${GREEN}✓ PASSED${NC} (correctly rejected)"
    ((PASSED++))
else
    echo -e "Test: me without auth ... ${RED}✗ FAILED${NC}"
    ((FAILED++))
fi

if $CLI request list 2>&1 | grep -q "not authenticated"; then
    echo -e "Test: request list without auth ... ${GREEN}✓ PASSED${NC} (correctly rejected)"
    ((PASSED++))
else
    echo -e "Test: request list without auth ... ${RED}✗ FAILED${NC}"
    ((FAILED++))
fi

echo ""
echo "====================================="
echo "Test Results:"
echo "  Passed: $PASSED"
echo "  Failed: $FAILED"
echo "====================================="

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
