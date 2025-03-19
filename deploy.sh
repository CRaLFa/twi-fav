#!/bin/bash

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

main () {
    cd "$(dirname "$0")" || return
    deploy_webhook
}

main "$@"
