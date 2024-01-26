# Fleet Infrastructure

```bash
k3d cluster create my-cluster
```


```bash
export GITHUB_TOKEN=<your-token>
export GITHUB_USER=<your-username>

flux check --pre

flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=CodersInControl \
  --branch=main \
  --path=./infra/clusters/my-cluster \
  --personal \
  --components-extra image-reflector-controller,image-automation-controller
```

Щоб отримати доступ до [Weave GitOps](http://localhost:9001) та [PodInfo](http://localhost:9001) здійсніть перенаправлення наступних портів.

```bash
# WeaveGitOps (user: admin, password: 12345)
k port-forward svc/dashboard-weave-gitops 9001:9001 -n flux-system
# PodInfo
k port-forward svc/dev-podinfo 9898:9898 -n dev
```

Для зміни версії додатку PodInfo, потрібно змінювати маніфест `./clusters/my-cluster/podinfo/helmrelease.yaml`, а саме `spec.chart.spec.version: x.x.x`.

Щоб переглянути перелік доступних версій, пройдіть за посиланням <https://stefanprodan.github.io/podinfo/index.yaml>
