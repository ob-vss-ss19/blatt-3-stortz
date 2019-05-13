## Ausführen mit Docker

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    sudo docker run --rm --net actors --name treeservice terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stortz:develop-treeservice --bind="treeservice.actors:8090"
    ```

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    sudo docker run --rm --net actors --name treecli terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stortz:develop-treecli --bind="treecli.actors:8091" --remote="treeservice.actors:8090" trees
    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass eine Liste aller Trees anzeigt.
    Für weitere Verwendung siehe CLI-Befehlsübersicht.


-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

    
## Ausführen ohne Docker
    
-  Clone Repository
    
    ```
    git clone https://github.com/ob-vss-ss19/blatt-3-stortz
    ```
    
    Starten des Services.
    ```
    cd blatt-3-stortz/treeservice
    go build
    ./treeservice
    ```
    
    Starten der CLI in einem zweiten Terminal.
    ```
    cd blatt-3-lallinger/treecli
    go build
    ./treecli trees
    ```
    Für weitere Verwendung siehe CLI-Befehlsübersicht.

## CLI-Befehlsübersicht

### Befehle
  ```
  create
    Creates a new tree; returns the id and a token
  add key value
    Inserts a key-value-pair into the specified tree
  remove key
    Removes the key-value-pair with the given key
  find key
    Returns the key-value pair for the given key
  traverse
    Returns the traversed tree
  delete
    Deletes the specified tree
    -Requires authorized flag set to true
  ```
  
 ### Flags
  ```
  --token=string
    token for tree access 
  --id=int
    specifies tree-id
  --remote=string
    remote adress (default: "localhost:8092")
  --bind=string
    bind adress (default: "127.0.0.1:8091")
  --authorized=bool
    authorizes for dangerous commands like delete (default: false)
  ```
