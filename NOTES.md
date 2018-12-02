# Pourquoi le plugin terraform-provider-vra7 doit être réécrit pour BMRC

Le plugin ogirinal terraform-provider-vra7 est une sort de proof-of-concept
qui ne contient que `vra7_resource` qui ne permet que de créer, mettre à jour
ou détruire une VM. Pour faire en sorte qu'il soit suffisemment maléable
(car vmware est très maléable), ils utilisent des

```hcl
resource_configuration = {
    ComponentName.property = value
    Rehl74.cpu = 1
    ComponentName.MetaProperty.property = value
}
```

Je propose d'y aller "plain and simple" en mode Terraform. C'est à dire que
tous les champs seront dans le schema. Cela permet d'éviter le problème
des requestTemplate et de bgName qui y a été oublié mais qui est obligatoire.

Terraform State {
    resource_configuration {
        ensemble du POST /requests (tous les champs ?)
    }
    
}

Problèmes:

- `POST /catalog-service/api/consumer/entitledCatalogItems/{{id}}/requests`
  demande d'avoir un Json avec les bons types. Or, comme terraform-provider-vra7
  utilise JSON.Marshal() qui va tout mettre en "string".

  Exemple de json requis :

  ```json
  {
      ...
      "data": {
          "cpu": 1, // ne peut pas être "1", doit être de type int
          "cactusGroups": ["BRMC000001"], // doit être un array
      }
  }
  ```

  Mais dans le schema, Terraform prend ça :

  ```go
  "resource_configuration": {Type: schema.TypeMap, Optional: true, Computed: true,
      Elem: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString}
  }
  ```

  Donc terraform va parser le bloc

  ```hcl
  resource_configuration {
      ""
  }
  ```

  Du coup il faut re-transformer en bons types par la suite...