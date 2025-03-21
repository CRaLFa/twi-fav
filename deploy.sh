#!/bin/bash

set -eu

deploy_webhook () {
    gcloud --project=${1} functions deploy twi-fav-webhook \
        --gen2 \
        --runtime=go123 \
        --region=asia-northeast1 \
        --source=./webhook \
        --entry-point=SaveLikedTweet \
        --trigger-http \
        --allow-unauthenticated \
        --service-account=cloud-run@${1}.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

deploy_api () {
    gcloud --project=${1} run deploy twi-fav-api \
        --region=asia-northeast1 \
        --source=./api \
        --allow-unauthenticated \
        --service-account=cloud-run@${1}.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

main () {
    cd "$(dirname "$0")"
    local project_id=$(yq -r '.PROJECT_ID' < .env.yaml)
    [[ $# -lt 1 || "$1" = 'webhook' ]] && {
        echo -e "Deploying webhook...\n"
        deploy_webhook "$project_id"
        echo
    }
    [[ $# -lt 1 || "$1" = 'api' ]] && {
        echo -e "Deploying api...\n"
        deploy_api "$project_id"
    }
}

main "$@"
