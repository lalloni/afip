Utilidades AFIP
===============

[![Build Status](https://travis-ci.org/lalloni/afip.svg?branch=master)](https://travis-ci.org/lalloni/afip)

Este proyecto contiene un pequeño conjunto de librerías de código de uso
frecuente en la interacción con servicios y herramientas de AFIP así como
en la implementación de funcionalidad de negocio.

Paquetes
--------

- `github.com/lalloni/afip/cuit` contiene funciones útiles para validar, parsear y formatear CUIT y CUIL. Ver su [documentación](https://godoc.org/github.com/lalloni/afip/cuit) para obtener más detalles.

Roadmap
-------

- `github.com/lalloni/afip/token` contiene funciones útiles para validar, parsear y generar tokens de autenticación.
- `github.com/lalloni/afip/signature` contiene funciones útiles para validar tokens de autenticación con la firma correspondiente.
- `github.com/lalloni/afip/clavefiscal` contiene middleware HTTP y funciones útiles para implementar autenticación con Clave Fiscal.
- `github.com/lalloni/afip/sua` contiene  middleware HTTP y funciones útiles para implementar autenticación con SUA.
- `github.com/lalloni/afip/wsaa` contiene  middleware HTTP y funciones útiles para implementar autenticación con WSAA.
