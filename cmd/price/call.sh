#!/user/bin/env sh

URL="localhost:8180"

RED='\033[0;31m'
NC='\033[0m'
YELLOW='\033[0;93m'

###############################
### sum service
###############################
echo -e "${YELLOW}==> Sum endpoint${NC}"
req=$(curl -s -i -X POST $URL/sum -d '{"price": 1, "fee":1}' ) 2>&1
status=$(echo "${req}" | head -1 ) 
data=$(echo "${req}" | tail -1 ) 
echo "Status: $status" 
echo -e "$data\n"

req=$(curl -s -i -X POST $URL/sum -d '{"price": 1, "fee":"1"}' ) 2>&1
status=$(echo "${req}" | head -1 ) 
data=$(echo "${req}" | tail -1 ) 
echo "Status: $status" 
echo -e "$data\n"
