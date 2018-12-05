# Demandes pour aider à l'utilisation de l'API

## Important: champ `bgName` manquant dans les templates

Quand on fait un GET sur les créations DCaaS/FAST :
`/catalog-service/api/consumer/entitledCatalogItems/<id>/requests/template`,
il manque le champ `bgName` alors qu'il est obligatoire.

## Important: certains champs ne parsent pas les strings

Il s'agit d'un embêtement car il faut beaucoup modifier le provider Terraform
proposé par VMWare pour corriger ça :

1. la plupart des champs (comme `cpu`, `ram`) acceptent un string qui sera
  parsé par l'API en integer ou laissée en string (en fonction du champs).
2. les champs `hasBasicat`", `isBGFast`, `leaseUnlimited`, `customRole`,
  `addDRSGroup`, `useCloudInit` demandent que la valeur dans le json soit
  true/false.
3. les champs `lease`, `targetDiskSizeOfVm` demandent des integers.

Serait-il possible de rendre les points (2) et (3) cohérents avec le point
(1) ? C'est à dire faire en sorte que (par exemple pour `hasBasicat`) on
puisse donner `"false"` (un string) au lieu de `false` (un boolean).

Pareil pour les integers `lease` et `targetDiskSizeOfVm`, faudrait qu'on puisse
donner `"2"` à la place de `2` par exemple.