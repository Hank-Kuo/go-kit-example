#!/user/bin/env sh

URL="localhost:8180"

NC='\033[0m'
RED='\033[0;31m'
YELLOW='\033[0;93m'

###############################
### sum endpoint(http)
###############################
echo -e "${YELLOW}==> Sum endpoint http${NC}"
req=$(curl -s -i -X POST $URL/sum -d '{"price": 1, "fee":1}' ) 2>&1
status=$(echo "${req}" | head -1 ) 
version=$(echo "${req}" | sed -n -e 4p) 
data=$(echo "${req}" | tail -1 ) 
echo "Status: $status" 
if [[ $version =~ "X-Api-Version" ]]; then
   echo -e "Version: $version"
fi
echo -e "$data\n"

req=$(curl -s -i -X POST $URL/sum -d '{"price": 1, "fee":"1"}' ) 2>&1
status=$(echo "${req}" | head -1 ) 
version=$(echo "${req}" | sed -n -e 4p) 
data=$(echo "${req}" | tail -1 ) 
echo "Status: $status" 
if [[ $version =~ "X-Api-Version" ]]; then
   echo -e "Version: $version"
fi
echo -e "$data\n"

req=$(curl -s -i -X POST $URL/sum -d '{"price": 0, "fee":0}' ) 2>&1
status=$(echo "${req}" | head -1 ) 
version=$(echo "${req}" | sed -n -e 4p) 
data=$(echo "${req}" | tail -1 ) 
echo "Status: $status" 
if [[ $version =~ "X-Api-Version" ]]; then
   echo -e "Version: $version"
fi
echo -e "$data\n" 

###############################
### sum endpoint(grpc)
###############################
echo -e "${YELLOW}==> Sum endpoint grpc${NC}"
echo "health check api"
echo grpcurl -plaintext localhost:8181 grpc.health.v1.Health/Check
echo "Price sum api"
grpcurl -plaintext -d '{"price": 1, "fee":1}' localhost:8181 pb.Price/Sum
grpcurl -plaintext -d '{"price": 1, "fee":"1f"}' localhost:8181 pb.Price/Sum

###############################
### exchange endpoint
###############################








