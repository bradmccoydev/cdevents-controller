#1 /usr/bin/env sh

set -e

# wait for cdevents-controller
kubectl rollout status deployment/cdevents-controller --timeout=3m

# test cdevents-controller
helm test cdevents-controller
