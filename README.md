#TestCDIDataImpact

Pour exécuter le projet run la commande :
`run main.go`

Vous devez ajouter un fichier `.env ` avec cette exemple de template 

```
db_host=localhost
 db_name=user
 db_user=user
 db_pwd=
 mongo_uri=mongodb+srv://user:<password>@test-cdi.agxkj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority
```
Apres avoir saisit les parametres de votre BDD,

Vous devez executer l'endpoint :
``
POST /api/users
`` 
qui permet d'ajouter à la BDD la liste des users depuis le fichier ``DataSet.json`` tout en hachant le password avec bcrypt

Ensuite vous devez vous connecter tout en executant l'endpoint : 
``
POST /api/auth
``

``
BODY { "email": "nikkifarley@anivet.com",
          "password":"CGUsfQkS06lo7LM2guSV" }
``

Ceci permet de générer un BearerToken qu'on doit utiliser pour avoir l'accès aux autres Endpoints