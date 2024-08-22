#!/bin/bash

# Usage function to display help for the script
usage() {
    echo "Usage: $0 <AWS_PROFILE> <CONFIG_FILE>"
    echo "Example: $0 openshift-dev-vmware ./config/conf.env"
    exit 1
}

# Check if the correct number of arguments was provided
if [ "$#" -ne 2 ]; then
    usage
fi

# Assign arguments to variables
AWS_PROFILE=$1
CONFIG_FILE=$2

# Load the configuration from config file
source $CONFIG_FILE

# Now you can use the environment variables set in conf.env
echo "The configured vSphereApiVIP is: ${vSphereApiVIP}"
echo "The configured vSphereIngressVIP is: ${vSphereIngressVIP}"

IP_ADDRESSES=(${vSphereApiVIP} ${vSphereIngressVIP})

export AWS_PROFILE=$AWS_PROFILE

# Get the Hosted Zone ID
HOSTED_ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name=='devqe.ibmc.devcluster.openshift.com.'].Id" --output text)

for IP_ADDRESS in "${IP_ADDRESSES[@]}"; do
    echo "Checking for IP address: $IP_ADDRESS"
    output=$(aws route53 list-resource-record-sets --hosted-zone-id $HOSTED_ZONE_ID --query "ResourceRecordSets[?ResourceRecords[?Value=='$IP_ADDRESS']]" --output json)

    # Check if the output is an empty JSON array
    if [ "$output" == "[]" ]; then
        echo "No DNS record found for IP $IP_ADDRESS"
    else
        echo "DNS record found for IP $IP_ADDRESS:"
        echo "$output" | jq -r '.[].Name'
    fi
done;
