# Kiln : Backend Home Exercise

<div style="text-align:center; width: 100%">
  <img src="https://financialit.net/sites/default/files/kiln1.png"  width="40%">
</div>

Write a small Golang service that proxifies the `eth_gasPrice` RPC from an Ethereum execution node.

Link of the exercise : https://www.notion.so/External-Backend-Home-Exercice-c8fade2ea7df4bf7afc14e0e7fec519f

### How to test this small Web Service ?

**Step 1 :** Setup your Infura Key in the file `constants.go`

    InfuraApiKey = "<REPLACE-BY-YOUR-INFURA-KEY>"

**Step 2 :** Build the go binaries

    go build

**Step 3 :** Start minikube

    minikube start

**Step 4 :** point your shell to minikube’s docker-daemon

    eval $(minikube -p minikube docker-env)

**Step 5 :** Build the docker image

    docker build -t yani/kiln-backend-home-exercise .

**Step 6 :** Deploy the app to the minikube cluster

    kubectl apply -f k8s-deployment.yml

**Step 7 :** List the pods

    kubectl get pods

**Step 8 :**  Map a local port to a port inside the pod

> **Info :** Don't forget to replace `{POD_ID}` by your local pod ID

    kubectl port-forward kiln-backend-home-exercise-{POD_ID} 8080:8080

**Step 9 :** With another console fetch the Gas Price endpoint

    curl --url http://localhost:8080/eth/gasprice