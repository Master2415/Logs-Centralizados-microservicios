const { Given, When, Then, Before } = require("@cucumber/cucumber");
const assert = require("assert");
const axios = require("axios");
const messageSchema = require("../../schemas/messaje-schema");
const errorSchema = require("../../schemas/error-schema");
const listSchema = require("../../schemas/logList-schema");
const faker = require("@faker-js/faker");
const Ajv = require("ajv");
const ajv = new Ajv();
const addFormats = require("ajv-formats");

addFormats(ajv);

let response;
let URL_LOG = "http://localhost:8081/logs/";

Given('el usuario desea obtener registros de logs de la API', function () {
    return Promise.resolve();
});

Given('los registros están disponibles en la base de datos', function () {
    return Promise.resolve();
});

Given('los siguientes parámetros de consulta están disponibles para el usuario:',  function (dataTable) {
    return Promise.resolve();
});

When('el usuario envía una solicitud GET a \\/api\\/logs\\/', async function () {
    try {
        response = await axios.get(URL_LOG);
    } catch (error) {
        if (error.response) {
            response = error.response;
        } else {
            throw new Error('No se recibió una respuesta válida o el mensaje de error está vacío.');
        }
    }
});

Then('el código de respuesta de \\/api\\/logs debe ser {int}', function (expectedStatusCode) {
    assert.strictEqual(response.status, expectedStatusCode);
});

Then('el cuerpo de la respuesta debe contener un arreglo de LOGS', function () {
    if (response && response.data) {
        const body = response.data;
        valid = ajv.validate(listSchema, body);
        assert.strictEqual(valid, true);
    } else {
        throw new Error('No se recibió una respuesta válida o los datos están vacíos.');
    }
});

When('el usuario envía una solicitud GET a \\/api\\/logs\\/ con parametros que no retornan info', async function () {
    try {
        URL_LOG = "http://localhost:8081/logs/?page=1&pageSize=30&startDate=2024-04-06T10:25:00-05:00&endDate=2024-04-14T11:05:00-05:00&logType=NOINFO"
        response = await axios.get(URL_LOG);
    } catch (error) {
        if (error.response) {
            response = error.response;
        } else {
            throw new Error('No se recibió una respuesta válida o el mensaje de error está vacío.');
        }
    }
});


Then('el cuerpo de la respuesta debe contener un mensaje de error', function () {
    if (response && response.data) {
        const error = response.data;
        valid = ajv.validate(errorSchema, error);
        assert.strictEqual(valid, true);
    } else {
        throw new Error('No se recibió una respuesta válida o el mensaje de error está vacío.');
    }
});
