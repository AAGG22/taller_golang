# Branch Policy (Política de Ramas)
## Dinamica
Todo trabajo debe estar realizado en un a rama y cada rama debe establecer una convención de nombre para poder ser identificable y fácilmente determinar su funcionalidad.
## Convención de Nombres para Ramas
Ejemplo: nombre/tarea (Alfredo/modificar_metodo).

| Componente       | Ejemplo       | Descripción                                  |
|------------------|---------------|----------------------------------------------|
| **Nombre**       | `Alfredo`     | Nombre del desarrollador responsable         |
| **Separador**    | `/`           | Diagonal para separar componentes            |
| **Referencia**   | `issue-2`     | Número del Issue asociado (sin #)            |

## Estructura de Ramas
La rama contara con el nombre del autor, una diagonal y con el numero de “Issue”, el merge (que fusiona la rama con la rama principal "main") será consensuado por el grupo.
```mermaid
gitGraph
  commit
  branch Yamila/issue-1
  commit
  branch Gabriel/issue-3
  commit
  branch Alfredo/issue-2
  commit
  checkout main
  merge Yamila/issue-1
  merge Gabriel/issue-3
  merge Alfredo/issue-2
```

## Flujo de trabajo para Ramas
```mermaid 
flowchart TD
    A[["Crear Issue en GitHub"]] --> B[["git checkout -b nombre/issue-N"]]
    B --> C[["Desarrollar feature <br/> • Commits atómicos <br/> • Mensajes claros"]]
    C --> D[["git push origin nombre/issue-N"]]
    D --> E[["Crear Pull Request"]]
    E --> F[["Revisión de los 3 compañeros <br/> • Comentar cambios <br/> • Aprobar/Rechazar"]]
    F --> G{¿Cumple requisitos?}
    G -->|✅ Aprobado| H[["Merge a main <br/> • Squash commits <br/> • Borrar rama"]]
    G -->|❌ Rechazado| C
```
## Restricciones del Equipo
- [x] Nunca commitear directamente en main.
- [x] Nombrar las ramas con la convension de nombres establecida.
- [x] Probar los cambios localmente antes de subirlos.
- [x] Cada PR debe ser revisado por los 3 compañeros.

> [!NOTE]
> Un PR (Pull Request) es una funcionalidad clave en Git/GitHub que permite proponer cambios en un repositorio y solicitar que se fusionen con la rama principal (como main). Es el mecanismo central para la colaboración en proyectos de desarrollo.

| El equipo | | |
|-------------|------|--------|
| Yamila Silva | <img src="https://avatars.githubusercontent.com/u/113130495?s=64&v=4?size=80" width="50"> | [@yamila](https://github.com/YamilaAS) |
| Gabriel Villalobo | <img src="https://avatars.githubusercontent.com/u/40302142?s=64&v=4size=80" width="50"> | [@gabriel](https://github.com/kuhg) |
| Alfredo Galván | <img src="https://avatars.githubusercontent.com/u/95595323?v=4size=80" width="50"> | [@alfredo](https://github.com/AAGG22) |
