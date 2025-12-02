#!/bin/bash

# Tests against public endpoints to demonstrate real-world performance

set -e

REQUESTS=200
CONNECTIONS=(10 50 100)

echo ""
echo "Select test endpoint:"
echo "  1) google.com"
echo "  2) example.com"
echo "  3) Custom URL"
echo ""
read -p "Choice [1]: " choice
choice=${choice:-1}

case $choice in
    1)
        TEST_URL="https://www.google.com"
        TEST_NAME="google.com"
        ;;
    2)
        TEST_URL="http://example.com"
        TEST_NAME="example.com"
        ;;
    3)
        read -p "Enter URL: " TEST_URL
        TEST_NAME=$(echo $TEST_URL | sed -e 's|^[^/]*//||' -e 's|/.*$||')
        ;;
esac

echo ""
echo "Testing against: ${TEST_URL}"
echo ""

# Check endpoint availability
echo "Checking endpoint availability..."
if ! curl -s -f --max-time 5 ${TEST_URL} > /dev/null 2>&1; then
    echo "❌ Endpoint not responding or too slow"
    exit 1
fi

# Measure baseline RTT
echo "Measuring baseline latency..."
HOST=$(echo $TEST_URL | sed -e 's|^[^/]*//||' -e 's|/.*$||' -e 's|:.*$||')
if ping -c 5 -W 2 ${HOST} > /dev/null 2>&1; then
    RTT=$(ping -c 5 ${HOST} | tail -1 | awk '{print $4}' | cut -d '/' -f 2)
    echo "   Average RTT: ${RTT}ms"
else
    echo "   (ICMP blocked, measuring HTTP latency...)"
    HTTP_TIME=$(curl -o /dev/null -s -w '%{time_total}\n' ${TEST_URL})
    echo "   HTTP response time: ${HTTP_TIME}s"
fi
echo ""

# Create results directory
RESULTS_DIR="results/remote/${TEST_NAME}_$(date +%Y%m%d_%H%M%S)"
mkdir -p ${RESULTS_DIR}

cat > ${RESULTS_DIR}/info.txt << EOF
Endpoint: ${TEST_URL}
Target: ${TEST_NAME}
RTT: ${RTT:-N/A}ms
Date: $(date)
Client: $(uname -a)
Duration: ${REQUESTS}
EOF

echo "Running benchmarks..."
echo ""

run_benchmark() {
    local tool=$1
    local connections=$2
    local url=${TEST_URL}

    echo "  Testing ${tool} with ${connections} connections..."

    case $tool in
        "volt")
            # Changed -json to redirect to volt.json to keep output cleaner
            .././volt bench -url ${url} -c ${connections} -n ${REQUESTS} -json > ${RESULTS_DIR}/${tool}_c${connections}.json 2>&1 || echo "Failed"
            ;;
        "hey")
            hey -c ${connections} -n ${REQUESTS} ${url} > ${RESULTS_DIR}/${tool}_c${connections}.txt 2>&1 || echo "Failed"
            ;;
    esac

    sleep 1  # Cool down between tests
}

TOOLS=()
[ -f ".././volt" ] && TOOLS+=("volt")
command -v hey &> /dev/null && TOOLS+=("hey")


if [ ${#TOOLS[@]} -eq 0 ]; then
    echo "❌ No benchmarking tools found!"
    exit 1
fi

echo "Available tools: ${TOOLS[@]}"
echo ""

# Run benchmarks
for conn in "${CONNECTIONS[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing with ${conn} connections"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    for tool in "${TOOLS[@]}"; do
        run_benchmark ${tool} ${conn}
    done
    echo ""
done

echo "Results saved to: ${RESULTS_DIR}"
echo ""

# Generate summary (combined with advantage calculation)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "QUICK SUMMARY - ${TEST_NAME}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

for conn in "${CONNECTIONS[@]}"; do
    echo "Connections: ${conn}"
    echo "─────────────────────────────────────────────────────"

    # Correctly Parse Volt results using .summary.throughput
    VOLT_RPS="N/A"
    if [ -f "${RESULTS_DIR}/volt_c${conn}.json" ]; then
        VOLT_RPS=$(jq -r '.summary.throughput // "N/A"' ${RESULTS_DIR}/volt_c${conn}.json 2>/dev/null)
        printf "  %-10s %s req/s\n" "Volt:" "${VOLT_RPS}"
    fi

    # Parse hey results
    HEY_RPS="N/A"
    if [ -f "${RESULTS_DIR}/hey_c${conn}.txt" ]; then
        HEY_RPS=$(grep "Requests/sec:" ${RESULTS_DIR}/hey_c${conn}.txt | awk '{print $2}')
        printf "  %-10s %s req/s\n" "hey:" "${HEY_RPS:-N/A}"
    fi

    # Calculate advantage for the current connection count
    if [ "$VOLT_RPS" != "N/A" ] && [ "$HEY_RPS" != "N/A" ] && [ "$HEY_RPS" != "0" ]; then
        # Round to one decimal place: add 0.05 and truncate using scale=1
        ADVANTAGE=$(echo "scale=1; ($VOLT_RPS / $HEY_RPS) + 0.05" | bc)
        echo "  Volt advantage: ${ADVANTAGE}x faster"
    elif [ "$VOLT_RPS" != "N/A" ] && [ "$HEY_RPS" = "0" ]; then
        echo "  Volt advantage: Infinite (hey RPS is zero)"
    fi
    echo ""
done