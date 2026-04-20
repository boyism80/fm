"use strict";

const { Repository } = require("./repository");

class ValueRepository extends Repository {
    async getAll() {
        throw new Error(`${this.constructor.name}.getAll is not supported for value repositories`);
    }

    async delAll() {
        throw new Error(`${this.constructor.name}.delAll is not supported for value repositories`);
    }
}

module.exports = { ValueRepository };
