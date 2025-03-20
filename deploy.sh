#!/bin/bash

set -eu

deploy_webhook () {
    gcloud --project=sandbox-morita-1-441408 functions deploy twi-fav-webhook \
        --gen2 \
        --runtime=go123 \
        --region=asia-northeast1 \
        --source=./webhook \
        --entry-point=SaveLikedTweet \
        --trigger-http \
        --allow-unauthenticated \
        --service-account=cloud-run@sandbox-morita-1-441408.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

deploy_api () {
    gcloud --project=sandbox-morita-1-441408 run deploy twi-fav-api \
        --region=asia-northeast1 \
        --source=./api \
        --allow-unauthenticated \
        --service-account=cloud-run@sandbox-morita-1-441408.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

main () {
    cd "$(dirname "$0")"
    [[ $# -lt 1 || "$1" = 'webhook' ]] && {
        echo -e "Deploying webhook...\n"
        deploy_webhook
        echo
    }
    [[ $# -lt 1 || "$1" = 'api' ]] && {
        echo -e "Deploying api...\n"
        deploy_api
    }
}

main "$@"
