# promoter

manifest project structure should be like this:

chart-repo/
-- project1/
---- service/
------ values-${env}.yaml
-- project2/
----  service/
------  values-${env}.yaml

the values files should be hooked to some GitOps provider (ArgoCD, Flux..)
