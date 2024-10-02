Feature: Recuperación de Logs de la API

    Background:
        Given el usuario desea obtener registros de logs de la API
        And los registros están disponibles en la base de datos
        And los siguientes parámetros de consulta están disponibles para el usuario:
            | Parameter | Description                             |
            | page      | Número de página                        |
            | pageSize  | Tamaño de página                        |
            | startDate | Fecha de inicio para filtrar logs       |
            | endDate   | Fecha de finalización para filtrar logs |
            | logType   | Tipo de log para filtrar                |

    Scenario: Recuperación exitosa de logs con parámetros válidos
        When el usuario envía una solicitud GET a /api/logs/
        Then el código de respuesta de /api/logs debe ser 200
        And el cuerpo de la respuesta debe contener un arreglo de LOGS

    Scenario: No se encuentran logs dentro de los parámetros de consulta
        When el usuario envía una solicitud GET a /api/logs/ con parametros que no retornan info
        Then el código de respuesta de /api/logs debe ser 404
        And el cuerpo de la respuesta debe contener un mensaje de error
