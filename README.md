#TestRingOver
Pour exécuter le projet run la commande :
`make run`

Vous devez ajouter un fichier `.env ` avec cette exemple de template 

```
db_host=localhost
db_name=ring-over-api
db_user=user
db_pwd=password
db_port=3306
```
Apres avoir saisit les parametres de votre BDD vous pouvez tester les endpoints

GET `/todolists`

Permet d’afficher la todo liste en JSON avec un limit/offset en paramètre.

POST `/todolists`

 Permet de créer un nouveau todo dans la liste

GET `/todolists/:id`

Permet de supprimer un todo dans la liste