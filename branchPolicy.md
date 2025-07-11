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
    G -->|✅ Aprobado| H[["Merge a main <br/> • Squash commits "]]
    G -->|❌ Rechazado| C
```
## El Equipo

| Foto | Colaborador |
|------|-------------|
| <img src="https://avatars.githubusercontent.com/u/113130495?v=4" width="50" style="border-radius:50%"> | [Yamila Silva](https://github.com/YamilaAS) |
| <img src="https://avatars.githubusercontent.com/u/40302142?v=4" width="50" style="border-radius:50%"> | [Gabriel Villalobos](https://github.com/kuhg) |
| <img src="https://avatars.githubusercontent.com/u/95595323?v=4" width="50" style="border-radius:50%"> | [Alfredo Galván](https://github.com/AAGG22) |

## Restricciones del Equipo
1. **Nunca commitear directamente en main.**
2. **Nombrar las ramas con la convension de nombres establecida.**
3. **Probar los cambios localmente antes de subirlos.**
4. **Cada PR debe ser revisado por los 3 compañeros.**
   - [x] Require pull request approvals before merging
   - [x] Required number of approvals before merging
   - [x] Require review from Code Owners (opcional pero recomendado)

> [!NOTE]
> Un PR (Pull Request) es una funcionalidad clave en Git/GitHub que permite proponer cambios en un repositorio y solicitar que se fusionen con la rama principal (como main). Es el mecanismo central para la colaboración en proyectos de desarrollo.


